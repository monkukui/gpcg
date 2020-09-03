package gpt

import (
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	// "reflect"
	// "bufio"

	"go/ast"
	"go/format"
	"go/parser"
	"go/token"

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
	dir, err := parser.ParseDir(token.NewFileSet(), "testdata/src/a/lib", nil, 0)
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

	// decl を集計する
	var decls []ast.Decl
	for _, v := range dir {
		for _, f := range v.Files {

			packageName := f.Name.Name
			ast.Inspect(f, func(n ast.Node) bool {
				switch n := n.(type) {
				case *ast.GenDecl:
					// インポート文以外を全て集計
					switch n.Tok {
					case token.IMPORT:
						return true
					case token.TYPE:
						for _, spec := range n.Specs {
							spec, _ := spec.(*ast.TypeSpec)
							if spec == nil {
								return true
							}
							rename(&spec.Name.Name, packageName)
						}
					case token.CONST:
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
						// レシーバ名を変更
						for _, field := range n.Recv.List {
							ident := field.Type.(*ast.Ident)
							if ident == nil {
								return true
							}
							rename(&ident.Name, packageName)
						}
					} else {
						// レシーバを持たない関数（構造体のメンバ変数以外）は，関数名を変更
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
				cr.Replace(node.Sel)
			}
		}

		return true
	}, nil)

	// 標準出力ではなく，file 出力にしたい

	// コードを生成して終了
	file, err := os.Create("./gen/gen.go")
	defer file.Close()
	if err != nil {
		return nil, err
	}
	generateCode(file, n)
	fmt.Println("gpt: generate code successfully✨")

	return nil, nil
}
