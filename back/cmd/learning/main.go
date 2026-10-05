package main

import (
	"fmt"
	"pouic/cmd/learning/a"
	"pouic/cmd/learning/b"
)

func main() {
	a := a.NewA(0)
	b := b.NewB(a)

	fmt.Printf("a = %d\n", a.GetI())
	a.IncrI()
	fmt.Printf("a = %d\n", a.GetI())
	fmt.Printf("a from b = %d\n", b.GetA().GetI())

	test := []byte(`hello`)

	fmt.Printf("%s\n", test)
}