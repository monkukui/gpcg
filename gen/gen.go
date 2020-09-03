package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
)

func (y Y) Error() string {
	return ""
}

type Y struct{}

func (s *S) Hoge() string {
	return ""
}

type S struct{}
type I interface{ Hoge() string }

func (y Y) Error() string {
	return ""
}

type Y struct{}

func (s *S) Hoge() string {
	return ""
}

type S struct{}
type I interface{ Hoge() string }

func main() {
	r := bufio.NewReader(os.Stdin)
	w := bufio.NewWriter(os.Stdout)
	defer w.Flush()
	var n, m int
	fmt.Fscan(r, &n, &m)
	uf := NewUnionFind(n)
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
