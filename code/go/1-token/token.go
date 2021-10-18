package main

import (
	"fmt"
	"go/scanner"
	"go/token"
)

func main() {
	src := []byte(`
package main

import "fmt"

func main() {
	var s = "HelloWorld!"
	fmt.Println(s)
}
`)

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
