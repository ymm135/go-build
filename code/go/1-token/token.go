package main

import (
	"fmt"
	"go/scanner"
	"go/token"
	"io/ioutil"
)

func main() {

	// 读取demo.go 文件内容
	filepath := "/Users/zero/work/go/workspace/go-build/demo.go"
	src, err := ioutil.ReadFile(filepath)

	if err != nil {
		panic(err)
	}

	var s scanner.Scanner
	fset := token.NewFileSet()
	file := fset.AddFile("", fset.Base(), len(src))
	s.Init(file, src, nil, 0)

	for {
		pos, tok, lit := s.Scan()
		fmt.Printf("%-6s%-8s%q\n", fset.Position(pos), tok, lit)

		if tok == token.EOF {
			break
		}
	}
}
