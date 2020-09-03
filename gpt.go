package gpt

import (
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
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

func run(pass *analysis.Pass) (interface{}, error) {

	// 競技プログラミングライブラリを全て捜査して，必要な情報を取ってくる
	dir, err := parser.ParseDir(token.NewFileSet(), "testdata/src/a/lib", nil, 0)
	if err != nil {
		return nil, err
	}

	fmt.Println(len(pass.Files))

	// import 文を集計する
	var imports []*ast.ImportSpec
	for _, v := range dir {
		for _, f := range v.Files {
			for _, spec := range f.Imports {
				imports = append(imports, spec)
			}
		}
	}

	// -----

	const code = `
package p

type I interface{
	Hoge() string
}

type S struct {
}

func (s *S) Hoge() string{
	return ""
}

type Y struct {
}

func (y Y) Error() string{
	return ""
}

`

	fmt.Println("sample finish")

	// ----
	var decls []ast.Decl
	for _, v := range dir {
		for _, f := range v.Files {
			fset := token.NewFileSet()

			conf := types.Config{
				Importer: importer.Default(),
				Error: func(err error) {
					fmt.Printf("!!! %#v\n", err)
				},
			}

			info := &types.Info{
				// Types: map[ast.Expr]types.TypeAndValue{},
				Defs: map[*ast.Ident]types.Object{},
				// Uses:  map[*ast.Ident]types.Object{},
			}

			pkg, err := conf.Check("p", fset, []*ast.File{f}, info) // FIXME: ここで panic が発生する
			if err != nil {
				log.Fatal(err)
			}

			/* TODO: pkgやinfoを使う処理 */
			fmt.Println("success")
			fmt.Println("pkg = ", pkg)

			// 名前が衝突しないように，全ての識別子を置き換える
			// ただし，基本型は除く
			ast.Inspect(f, func(n ast.Node) bool {
				ident, _ := n.(*ast.Ident)
				if ident == nil {
					return true
				}
				// fmt.Println("ident = ", ident)

				expr, _ := n.(ast.Expr)
				if expr == nil {
					return true
				}
				// fmt.Println("expr = ", expr)

				objType := pass.TypesInfo.TypeOf(expr)
				// fmt.Println("objType = ", objType)
				// fmt.Println("types.Typ[types.Int] = ", types.Typ[types.Int])
				if types.Identical(objType, types.Typ[types.Int]) {
					// fmt.Println("ident = ", ident)
					return true
				}
				// if ident.(type) == ast.Expr {
				// fmt.Println("expr!!!!")
				// }

				// fmt.Println("rename ident = ", ident)
				return true
			})

			// decl を集計する
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
