package main

type Man struct {
	Name string
	Age  int
}

func (man *Man) walk() string {
	return man.Name
}

func main() {
	man := Man{Name: "xiaoming", Age: 18}
	man.walk()
}
