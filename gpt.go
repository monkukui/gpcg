package gpt

import (
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
	// "sync"
	// "reflect"
	// "bufio"
	// "bytes"

	"go/ast"
	"go/format"
	"go/importer"
	"go/parser"
	"go/token"
	"go/types"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/ast/astutil"
)

const doc = "gpt is ..."
const lib = "a"

// Analyzer is ...
var Analyzer = &analysis.Analyzer{
	Name: "gpt",
	Doc:  doc,
	Run:  run,
}

func generateCode(w io.Writer, node interface{}) {
	format.Node(w, token.NewFileSet(), node)
}

func rename(name *string, packageName string) {
	*name = "generated_" + packageName + "_" + *name
}

func run(pass *analysis.Pass) (interface{}, error) {

	// 競技プログラミングライブラリを全て捜査して，必要な情報を取ってくる
	fset := token.NewFileSet()
	dir, err := parser.ParseDir(fset, "testdata/src/a/lib", nil, 0)
	if err != nil {
		return nil, err
	}

	// import 文を集計する
	var imports []*ast.ImportSpec
	for _, v := range dir {
		for _, f := range v.Files {
			for _, spec := range f.Imports {
				imports = append(imports, spec)
			}
		}
	}

	// ----
	var decls []ast.Decl
	for _, v := range dir {
		for _, f := range v.Files {
			packageName := f.Name.Name

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

			pkg, err := conf.Check("lib", fset, []*ast.File{f}, info)
			if err != nil {
				log.Fatal(err)
			}

			// 名前が衝突しないように，全ての識別子を置き換える
			// ただし，基本型は除く

			packageScope := pkg.Scope()
			// fmt.Println("packageScope = ", packageScope)
			// decl を集計する
			ast.Inspect(f, func(n ast.Node) bool {
				switch n := n.(type) {
				case *ast.Ident:
					if obj := packageScope.Lookup(n.Name); obj != nil {
						rename(&n.Name, packageName)
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
							rename(&spec.Name.Name, packageName)
						}
					case token.CONST:
						// 変数定義を rename
						for _, spec := range n.Specs {
							spec, _ := spec.(*ast.ValueSpec)
							if spec == nil {
								return true
							}
							for _, ident := range spec.Names {
								rename(&ident.Name, packageName)
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
								rename(&ident.Name, packageName)
							}
						}
					}

					decls = append(decls, n)
				case *ast.FuncDecl:
					// 関数定義は全て集計
					if n.Recv != nil {
						// レシーバ名を変更 （TODO 若干怪しい，レシーバ名は，このライブラリで定義された構造体であることを仮定においている）
						for _, field := range n.Recv.List {
							ident := field.Type.(*ast.Ident)
							if ident == nil {
								return true
							}
							rename(&ident.Name, packageName)
						}
					} else {
						// レシーバを持たない関数定義（構造体のメンバ変数以外）は，関数名を変更
						rename(&n.Name.Name, packageName)
					}
					decls = append(decls, n)
				}

				return true
			})
		}
	}

	// main 関数に対するコードの編集
	mainFile := pass.Files[0]

	insertImportsFlag := false
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

		case *ast.ImportSpec:
			// import a/lib など，ローカルからインポートしている文を削除する
			path, err := strconv.Unquote(node.Path.Value)
			if err != nil {
				return true
			}
			if strings.HasSuffix(path, "lib") {
				cr.Delete()
			}

			// ImportSpec の一番最初の時，蓄えた import 文を後ろに挿入していく
			if insertImportsFlag {
				return true
			}
			if cr.Index() == 0 {
				for _, spec := range imports {
					cr.InsertAfter(spec)
				}
			}
			insertImportsFlag = true

		case *ast.SelectorExpr:
			// lib.HogeHuga() -> HogeHuga() に置換する
			// lib.HogeHuga -> HogeHuga に置換する

			// node.X が 識別子ではない場合は無視
			ident, ok := node.X.(*ast.Ident)
			if !ok {
				return true
			}

			// 識別子の名前が lib の時は，現在のノードをごっそり node.Sel に置換
			if ident.Name == "lib" {
				rename(&node.Sel.Name, "lib")
				cr.Replace(node.Sel)
			}
		}

		return true
	}, nil)

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

	ff, _ := n.(*ast.File)
	if ff == nil {
		panic("file じゃありません")
	}
	_, err = conf.Check("lib", pass.Fset, []*ast.File{ff}, info)
	if err != nil {
		log.Fatal(err)
	}

	// ここから，定義だけされているが使われていない，関数，構造体，変数を消していく
	m := astutil.Apply(n, func(cr *astutil.Cursor) bool {
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
				// 構造体の定義を rename
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
				// 変数定義を rename
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
				// 変数定義を rename
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
		/*ident, _ := cr.Node().(*ast.Ident)
		  if ident == nil {
		    return true
		  }

		  if !isUsed(info, ident) {
		    fmt.Println("ident = ", ident)
		  }
		*/

		return true
	}, nil)

	// コードを生成して終了
	file, err := os.Create("./gen/gen.go")
	defer file.Close()
	if err != nil {
		return nil, err
	}
	generateCode(file, m)
	fmt.Println("gpt: generate code successfully✨")

	return nil, nil
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
