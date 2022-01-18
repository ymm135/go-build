package main

import "time"

func main() {
	a := 5
	go func() {
		a = 6
	}()

	time.Sleep(time.Second)
}
