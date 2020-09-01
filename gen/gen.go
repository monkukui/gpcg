package main

import (
	// "a/lib"
	"fmt"
)

func main() {
	fmt.Println("Hello")
	fmt.Println("World")
	gopher := 3
	fmt.Println(gopher)
	fmt.Println(ModPow())
	uf := UnionFind{4}
	fmt.Println(uf.N)
}
// package lib

func ModPow() int64 {
	return 1000000007
}
// package lib

type UnionFind struct{ N int }

