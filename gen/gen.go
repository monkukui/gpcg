package main

import (
	"fmt"
)

func main() {
	fmt.Println("Hello")
	fmt.Println("World")
	gopher := 3
	fmt.Println(gopher)
	fmt.Println("lib.ModPow() = ", ModPow())
	uf := UnionFind{4}
	fmt.Println(uf.N)
}

func ModModPow() int64 {
	return 1
}

func ModPow() int64 {
	return 111
}

type UnionFind struct{ N int }
