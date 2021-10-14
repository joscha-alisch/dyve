package main

import (
	"fmt"
)

type Func func(something Bla) error

func main() {
	a := []Bla{}
	doSomething(&a)
	fmt.Printf("%+v", a)
}

type Bla struct {
	Test string
}

func doSomething(a *[]Bla) {
	*a = append(*a, Bla{})
}
