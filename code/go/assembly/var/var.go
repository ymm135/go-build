package main

import "fmt"

type IMan interface {
	walk()
}

type Man struct {
	name string
}

func (man *Man) walk() {
	man.name = "xiaoming"
}

func main() {
	var a int
	a = 2

	b := new(int)
	b = &a
	*b = 5

	s := make([]int, 6)
	s[0] = 2

	c := make(chan int, 5)
	c <- 3

	man := Man{}
	var manImpl IMan
	manImpl = &man

	fmt.Println(a, b, s, <-c, manImpl)
}
