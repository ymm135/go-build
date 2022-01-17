package main

import "fmt"

func main() {
	b := 8
	c := callMe(b)
	fmt.Println(c)
}

func callMe(a int) int {
	return a + 5
}
