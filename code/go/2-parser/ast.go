package main

import (
	"go/ast"
	"go/parser"
	"go/token"
	"io/ioutil"
	"log"
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
		log.Fatal(err)
	}

	ast.Print(fset, file)
}