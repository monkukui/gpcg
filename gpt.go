package gpt

import (
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"os"
	// "fmt"
	"strconv"
	"strings"
	// "reflect"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/astutil"
	// "golang.org/x/tools/go/ast/inspector"
)

const doc = "gpt is ..."
const lib = "a"

// Analyzer is ...
var Analyzer = &analysis.Analyzer{
	Name: "gpt",
	Doc:  doc,
	Run:  run,
	Requires: []*analysis.Analyzer{
		inspect.Analyzer,
	},
}

func generateCode(node interface{}) {
	format.Node(os.Stdout, token.NewFileSet(), node)
}

func generateEoln() {
	node, err := parser.ParseExpr("1 + 2")
	if err != nil {
		panic(err)
	}
	format.Node(os.Stdout, token.NewFileSet(), node)
}

func run(pass *analysis.Pass) (interface{}, error) {

	// main 関数に相当する，a file をよむ
	f, err := parser.ParseFile(token.NewFileSet(), "testdata/src/a/a.go", nil, 0)
	if err != nil {
		return nil, err
	}

	// main 関数に対するコードの編集
	n := astutil.Apply(f, func(cr *astutil.Cursor) bool {
		switch node := cr.Node().(type) {
		case *ast.ImportSpec:
			// import a/lib など，ローカルからインポートしている文を削除する

			path, err := strconv.Unquote(node.Path.Value)
			if err != nil {
				return true
			}
			if strings.HasSuffix(path, "lib") {
				cr.Delete()
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

	generateCode(n)

	// parser.ParseDir を読んで，ディレクトリ単位で ast を得る
	// fset := token.NewFileSet()
	d, err := parser.ParseDir(token.NewFileSet(), "testdata/src/a/lib", nil, 0)
	if err != nil {
		return nil, err
	}

	for _, v := range d {
		for _, file := range v.Files {

			// package lib を除いたコードを出力（ライブラリ内の import 文は一旦むし）
			generateCode(file.Decls)
		}
	}

	/*
	  // main 内の import 文を捜査
	  // fmt.Println("len = ", len(pass.Files))
	  paths := []*string{}
	  for _, f := range pass.Files {
	    // fmt.Println("file = ", f)
	    imports, err := findLocalImports(f)
	    if err != nil {
	      return nil, err
	    }
	    for _, i := range imports {
	      paths = append(paths, i)
	    }
	  }

	  // 対象ファイルの抽象構文木を取得
	  for i, path := range paths {

	    fmt.Println("依存ライブラリ ", i)
	    fmt.Println(*path)
	    fset := token.NewFileSet()
	    // f, err := parser.ParseFile(fset, "./testcase/src/" + *path + "/graph.go", nil, 0)
	    f, err := parser.ParseFile(fset, "testdata/src/a/lib/graph/union_find.go", nil, 0)
	    if err != nil {
	      return nil, err
	    }
	    fmt.Println(f)
	  }

		inspect := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)

		nodeFilter := []ast.Node{
			(*ast.Ident)(nil),
		}

		inspect.Preorder(nodeFilter, func(n ast.Node) {
			switch n := n.(type) {
			case *ast.Ident:
	      // fmt.Println(n)
				if n.Name == "gopher" {
	        fmt.Println(n.Name)
					pass.Reportf(n.Pos(), "identifier is gopher")
				}
			}
		})

	*/

	return nil, nil
}
