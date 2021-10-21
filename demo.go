package main

func sum(x int,y int) int {
	return x + y
}

func main() {
	s := make([]int, 5)
	s[0] = 1
	s[1] = 9

	sum(s[0], s[1])
}