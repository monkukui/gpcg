package gpt

import (
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"os"
	"fmt"
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

func run(pass *analysis.Pass) (interface{}, error) {

  // main と lib の imports 文を集計する
  imports := []*ast.ImportSpec{}

  // a.go 内の import 文を走査
  for _, f := range pass.Files {
    for _, spec := range f.Imports {
      imports = append(imports, spec)
      path, err := strconv.Unquote(spec.Path.Value)

      // ローカルのライブラリをインポートしていたら読み飛ばす
      if strings.HasSuffix(path, "lib") {
        continue
      }
      if err != nil {
        return nil, err
      }
    }
  }

  // lib 内の import 文を走査
	dir, err := parser.ParseDir(token.NewFileSet(), "testdata/src/a/lib", nil, 0)
	if err != nil {
		return nil, err
	}

	for _, v := range dir {
		for _, f := range v.Files {
      for _, spec := range f.Imports {
        imports = append(imports, spec)
      }
		}
	}

	// main 関数に相当する，a file をよむ
	f, err := parser.ParseFile(token.NewFileSet(), "testdata/src/a/a.go", nil, 0)
	if err != nil {
		return nil, err
	}

  decls := []ast.Decl{}

	for _, v := range dir {
		for _, f := range v.Files {

			// package lib を除いたコードを出力（ライブラリ内の import 文は一旦むし）
      // TODO Aooly ではなく，トラバース
      m := astutil.Apply(f, func(cr *astutil.Cursor) bool {
        switch node := cr.Node().(type) {
        case *ast.GenDecl:
          if node.Tok != token.IMPORT {
            decls = append(decls, node)
          }
        case *ast.FuncDecl:
          decls = append(decls, node)
        }

        return true
      }, nil)

      fmt.Println("m = ", m)
		}
	}

  for _, decl := range decls {
    fmt.Println(decl)
    f.Decls = append(f.Decls, decl)
  }

	// main 関数に対するコードの編集
	n := astutil.Apply(f, func(cr *astutil.Cursor) bool {
		switch node := cr.Node().(type) {
    case *ast.GenDecl:
      // lib 以下の import を追加していく
      if node.Tok == token.IMPORT {
        for _, spec := range imports {
          node.Specs = append(node.Specs, spec)
        }
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



	return nil, nil
}
