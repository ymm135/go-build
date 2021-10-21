package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"io/ioutil"
	"os"
)

func main() {
	// 读取demo.go 文件内容
	filepath := "/Users/zero/work/go/workspace/go-build/demo.go"
	src, err := ioutil.ReadFile(filepath)

	if err != nil {
		panic(err)
	}

	fset := token.NewFileSet()

	file, err := parser.ParseFile(fset, "", src, 0)
	if err != nil {
		fmt.Println(err)
	}

	ast.Inspect(file, func(n ast.Node) bool {
		call, ok := n.(*ast.CallExpr)
		if !ok {
			return true
		}

		printer.Fprint(os.Stdout, fset, call.Fun)

		return false
	})
}