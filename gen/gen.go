package main

import (
	"fmt"
	"math"
	"strconv"
)

type UnionFind struct{ N int }

func ModInv() float64 {
	fmt.Println("hello mod inv")
	fmt.Println(strconv.Atoi("122"))
	return math.Min(1.0, 2.0)
}
func ModPow() float64 {
	fmt.Println("hello mod pow")
	return math.Max(1.0, 2.0)
}
func main() {
	fmt.Println("Hello")
	fmt.Println("World")
	gopher := 3
	fmt.Println(gopher)
	fmt.Println("lib.ModPow() = ", ModPow())
	fmt.Println("lib.ModInv() = ", ModInv())
	uf := UnionFind{4}
	fmt.Println(uf.N)
}
func ModModPow() int64 {
	return 1
}
