package main

import (
	"bufio"
	"fmt"
	"os"

	alib "a/lib"
)

func main() {
	r := bufio.NewReader(os.Stdin)
	w := bufio.NewWriter(os.Stdout)
	defer w.Flush()

	var n, m int
	fmt.Fscan(r, &n, &m)
	uf := alib.NewUnionFind(n)

	for i := 0; i < m; i++ {
		var a, b int
		fmt.Fscan(r, &a, &b)
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

	fmt.Fprintln(w, ans)
}
