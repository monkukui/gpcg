package main

import (
  "fmt"
  "a/lib/graph"
  "a/lib/math"
)

func main() {
	fmt.Println("Hello")
	fmt.Println("World")
  gopher := 3
  fmt.Println(gopher)
  fmt.Println(math.ModPow())
  uf := graph.UnionFind{4}
  fmt.Println(uf.N)
}
