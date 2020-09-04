package gpt

import (
	"fmt"
	"io"
	"log"
	"os"
	"strconv"

	"go/ast"
	"go/format"
	"go/importer"
	"go/parser"
	"go/token"
	"go/types"
	"io/ioutil"

	"golang.org/x/tools/go/ast/astutil"
	"golang.org/x/tools/imports"
)

const doc = "gpt is ..."

var libPackageName string = "lib"

const targetLibPath = "a/lib"

func Generate(mainPath, libPath, genPath string) error {

	// main 関数に対するコードの編集
	mainFileSet := token.NewFileSet()
	mainFile, err := parser.ParseFile(mainFileSet, mainPath, nil, 0)
	if err != nil {
		return err
	}

	for _, spec := range mainFile.Imports {
		path := spec.Path.Value
		path, err := strconv.Unquote(spec.Path.Value)
		if err != nil {
			return err
		}
		if path == targetLibPath {
			if spec.Name != nil {
				libPackageName = spec.Name.Name
				spec.Name = nil
			}
		}
	}

	// 競技プログラミングライブラリを全て捜査して，必要な情報を取ってくる
	fset := token.NewFileSet()
	dir, err := parser.ParseDir(fset, libPath, nil, 0)
	if err != nil {
		return err
	}

	var decls []ast.Decl
	for _, v := range dir {
		for _, f := range v.Files {
			if f.Name == nil {
				continue
			}

			conf := types.Config{
				Importer: importer.Default(),
				Error: func(err error) {
					fmt.Printf("!!! %#v\n", err)
				},
			}

			pkg, err := conf.Check("lib", fset, []*ast.File{f}, &types.Info{})
			if err != nil {
				log.Fatal(err)
			}

			packageScope := pkg.Scope()

			// 1. 名前が衝突しないように，識別子の名前を rename する（e.g. UnionFind -> generated_lib_UnionFind）
			// 2. コード生成のために Decl 文を集約する
			// a. 関数定義（ただし，レシーバを持つ関数は rename しない）
			// b. 構造体
			// c. 変数定義
			ast.Inspect(f, func(n ast.Node) bool {
				switch n := n.(type) {
				case *ast.Ident:
					if obj := packageScope.Lookup(n.Name); obj != nil {
						rename(&n.Name, libPackageName)
					}
				case *ast.GenDecl:
					switch n.Tok {
					case token.IMPORT:
						return true
					case token.TYPE:
						// 構造体の定義を rename
						for _, spec := range n.Specs {
							spec, _ := spec.(*ast.TypeSpec)
							if spec == nil {
								return true
							}
							rename(&spec.Name.Name, libPackageName)
						}
					case token.CONST:
						// 変数定義を rename
						for _, spec := range n.Specs {
							spec, _ := spec.(*ast.ValueSpec)
							if spec == nil {
								return true
							}
							for _, ident := range spec.Names {
								rename(&ident.Name, libPackageName)
							}
						}
					case token.VAR:
						// 変数定義を rename
						for _, spec := range n.Specs {
							spec, _ := spec.(*ast.ValueSpec)
							if spec == nil {
								return true
							}
							for _, ident := range spec.Names {
								rename(&ident.Name, libPackageName)
							}
						}
					}

					decls = append(decls, n)
				case *ast.FuncDecl:
					// 関数定義は全て集計
					if n.Recv != nil {
						// レシーバ名を変更
						for _, field := range n.Recv.List {
							ident := field.Type.(*ast.Ident)
							if ident == nil {
								return true
							}
							rename(&ident.Name, libPackageName)
						}
					} else {
						// レシーバを持たない関数定義（構造体のメンバ変数以外）は，関数名を変更
						rename(&n.Name.Name, libPackageName)
					}
					decls = append(decls, n)
				}

				return true
			})
		}
	}

	insertDeclsFlag := false
	n := astutil.Apply(mainFile, func(cr *astutil.Cursor) bool {
		switch node := cr.Node().(type) {
		case *ast.GenDecl:
			// Decl の一番最初の時，蓄えた decl 文を後ろに挿入していく
			if insertDeclsFlag {
				return true
			}
			for _, spec := range decls {
				cr.InsertAfter(spec)
			}
			insertDeclsFlag = true

		case *ast.SelectorExpr:
			// lib.HogeHuga() -> HogeHuga() に置換する
			// lib.HogeHuga -> HogeHuga に置換する

			// node.X が 識別子ではない場合は無視
			ident, ok := node.X.(*ast.Ident)
			if !ok {
				return true
			}

			// 識別子の名前が lib の時は，現在のノードをごっそり node.Sel に置換
			if ident.Name == libPackageName {
				rename(&node.Sel.Name, libPackageName)
				cr.Replace(node.Sel)
			}
		}

		return true
	}, nil)

	// 一旦ファイルに書き込む -> goimports をかける -> 再度 ast.File として読み込む
	f, _ := n.(*ast.File)
	if f == nil {
		log.Fatal("can not open the file")
	}

	f, mainFileSet, err = goimportsToFile(f)

	conf := types.Config{
		Importer: importer.Default(),
		Error: func(err error) {
			fmt.Printf("!!! %#v\n", err)
		},
	}

	info := &types.Info{
		Types:  map[ast.Expr]types.TypeAndValue{},
		Defs:   map[*ast.Ident]types.Object{},
		Uses:   map[*ast.Ident]types.Object{},
		Scopes: map[ast.Node]*types.Scope{},
	}

	_, err = conf.Check("lib", mainFileSet, []*ast.File{f}, info)
	if err != nil {
		log.Fatal(err)
	}

	// 定義だけされているが使われていない，関数，構造体，変数を消していく
	// 使用されていない識別子を cr.Delete() で削除する
	// 1. 関数定義（e.g. func ModPow() int など）
	// 2. 構造体定義 (e.g. type UnionFind struct など)
	// 3. 変数定義（e.g. var n int など）
	m := astutil.Apply(f, func(cr *astutil.Cursor) bool {
		switch node := cr.Node().(type) {
		case *ast.FuncDecl:
			// 関数定義
			if !isUsed(info, node.Name) {
				cr.Delete()
			}

		case *ast.GenDecl:
			switch node.Tok {
			case token.IMPORT:
				return true
			case token.TYPE:
				// 構造体の定義を削除
				for _, spec := range node.Specs {
					spec, _ := spec.(*ast.TypeSpec)
					if spec == nil {
						return true
					}
					if !isUsed(info, spec.Name) {
						cr.Delete()
						return true
					}
				}
			case token.CONST:
				// 変数定義を削除
				for _, spec := range node.Specs {
					spec, _ := spec.(*ast.ValueSpec)
					if spec == nil {
						return true
					}
					for _, ident := range spec.Names {
						if !isUsed(info, ident) {
							cr.Delete()
							return true
						}
					}
				}
			case token.VAR:
				// 変数定義を削除
				for _, spec := range node.Specs {
					spec, _ := spec.(*ast.ValueSpec)
					if spec == nil {
						return true
					}
					for _, ident := range spec.Names {
						if !isUsed(info, ident) {
							cr.Delete()
							return true
						}
					}
				}
			}
		}

		return true
	}, nil)

	// コードを生成して終了
	file, err := os.Create(genPath)
	defer file.Close()
	if err != nil {
		return err
	}

	f, _ = m.(*ast.File)
	if f == nil {
		log.Fatal("can not open the file")
	}

	f, _, err = goimportsToFile(f)
	if err != nil {
		return err
	}
	generateCode(file, f)
	fmt.Println("gpt: generate code successfully✨")
	return nil
}

func generateCode(w io.Writer, node interface{}) {
	format.Node(w, token.NewFileSet(), node)
}

func rename(name *string, packageName string) {
	*name = "generated_" + packageName + "_" + *name
}

func isUsed(info *types.Info, ident *ast.Ident) bool {
	obj := info.ObjectOf(ident)
	switch obj := obj.(type) {
	case *types.Func:
		// main
		if obj.Name() == "main" {
			return true
		}

		// init
		if obj.Name() == "init" {
			return true
		}
	}

	// どこかで使われているか
	for _, o := range info.Uses {
		if o == obj {
			return true
		}
	}

	return false
}

func goimportsToFile(f *ast.File) (*ast.File, *token.FileSet, error) {

	// 1. 適当にファイル出力をする
	outputFile, err := os.Create("./gen/tmp1.go")
	defer outputFile.Close()
	if err != nil {
		return nil, nil, err
	}
	generateCode(outputFile, f)

	// 2. 出力されたファイルに goimports をかける
	generatedFile, err := os.Open("./gen/tmp1.go")
	defer generatedFile.Close()
	if err != nil {
		return nil, nil, err
	}
	generatedCode, err := ioutil.ReadAll(generatedFile)
	if err != nil {
		return nil, nil, err
	}
	formatedCode, err := imports.Process("tmp1.go", generatedCode, nil)
	if err != nil {
		return nil, nil, err
	}

	// 3. フォーマットされたコードを書き込む
	formatedFile, err := os.Create("./gen/tmp2.go")
	defer formatedFile.Close()
	if err != nil {
		return nil, nil, err
	}
	formatedFile.Write(formatedCode)

	// 4. 再度ファイルを読み込む
	fset := token.NewFileSet()
	retFile, err := parser.ParseFile(fset, "./gen/tmp2.go", nil, 0)
	if err != nil {
		log.Print(err)
		return nil, nil, err
	}

	return retFile, fset, nil
}
