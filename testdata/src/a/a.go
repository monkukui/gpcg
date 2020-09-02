package main

import (
	"fmt"

	"a/lib"
)

func main() {
	var n, m int
	fmt.Scan(&n, &m)
	uf := lib.NewUnionFind(n)

	for i := 0; i < m; i++ {
		var a, b int
		fmt.Scan(&a, &b)
		a--
		b--

		uf.Union(a, b)
	}

	ans := 0
	for i := 0; i < n; i++ {
		if ans < uf.Size(i) {
			ans = uf.Size(i)
		}
	}

	fmt.Println(ans)
}
