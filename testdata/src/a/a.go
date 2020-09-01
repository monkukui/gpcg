package main

import (
  "fmt"
  "a/lib"
)

func main() {
	fmt.Println("Hello")
	fmt.Println("World")
  gopher := 3
  fmt.Println(gopher)
  fmt.Println(lib.ModPow())
  uf := lib.UnionFind{4}
  fmt.Println(uf.N)
}
