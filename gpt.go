package gpt

import (
	"go/ast"
  "go/parser"
  "go/token"
  "os"
  "go/format"
  // "fmt"
  "strconv"
  "strings"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
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

// Imports を走査し，目的のライブラリに対する import を探す
func findLocalImports(f *ast.File) ([]*string, error) {
  ret := []*string{}
  for _, spec := range f.Imports {
    path, err := strconv.Unquote(spec.Path.Value)
    if err != nil {
      return nil, err
    }

    if strings.HasPrefix(path, lib) {
      ret = append(ret, &path)
    }
  }
  return ret, nil
}

// ast.File を受け取って，ソースコードを出力する関数
func generateCodeByFile(file *ast.File) {
  format.Node(os.Stdout, token.NewFileSet(), file)
}

func run(pass *analysis.Pass) (interface{}, error) {

  // main 関数に相当する，a file をよむ
  f, err := parser.ParseFile(token.NewFileSet(), "testdata/src/a/a.go", nil, 0)
  if err != nil {
    return nil, err
  }
  generateCodeByFile(f)

  // parser.ParseDir を読んで，ディレクトリ単位で ast を得る
  // fset := token.NewFileSet()
  d, err := parser.ParseDir(token.NewFileSet(), "testdata/src/a/lib", nil, 0)
  if err != nil {
    return nil, err
  }

  for _, v := range d {
    for _, file := range v.Files {
      generateCodeByFile(file)
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
