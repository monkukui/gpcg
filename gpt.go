package gpt

import (
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"os"
	"strconv"
	"strings"
	// "reflect"

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

func generateCode(node interface{}) {
	format.Node(os.Stdout, token.NewFileSet(), node)
}

func run(pass *analysis.Pass) (interface{}, error) {

	// 競技プログラミングライブラリを全て捜査して，必要な情報を取ってくる
	dir, err := parser.ParseDir(token.NewFileSet(), "testdata/src/a/lib", nil, 0)
	if err != nil {
		return nil, err
	}

	// import 文を集計する
	imports := []*ast.ImportSpec{}
	for _, v := range dir {
		for _, f := range v.Files {
			for _, spec := range f.Imports {
				imports = append(imports, spec)
			}
		}
	}

	// decl を集計する
	decls := []ast.Decl{}
	for _, v := range dir {
		for _, f := range v.Files {
			ast.Inspect(f, func(n ast.Node) bool {
				switch n := n.(type) {
				case *ast.GenDecl:
					// インポート文以外を全て集計
					if n.Tok != token.IMPORT {
						decls = append(decls, n)
					}
				case *ast.FuncDecl:
					// 関数定義は全て集計
					decls = append(decls, n)
				}

				return true
			})
		}
	}

	// main 関数に対するコードの編集
	mainFile := pass.Files[0]

	n := astutil.Apply(mainFile, func(cr *astutil.Cursor) bool {
		switch node := cr.Node().(type) {
		case *ast.GenDecl:
			// Decl の一番最初の時，蓄えた decl 文を後ろに挿入していく
			for _, spec := range decls {
				cr.InsertAfter(spec)
			}

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
			if cr.Index() == 0 {
				for _, spec := range imports {
					cr.InsertAfter(spec)
				}
			}

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

	// コードを生成して終了
	fmt.Println("code gen begin")
	fmt.Println("-------------------------")
	generateCode(n)
	fmt.Println("-------------------------")
	fmt.Println("code gen end")
	return nil, nil
}
