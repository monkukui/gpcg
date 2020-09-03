package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
)

func unUsedFunction() string {
	return "I am unused unexported function"
}
func UnUsedFunction() string {
	return "I am unused exported function"
}
func swap(a int, b int) (int, int) {
	return b, a
}
func (u UnionFind) Size(x int) int {
	return -u.par[u.Find(x)]
}
func (u UnionFind) Same(x, y int) bool {
	return u.Find(x) == u.Find(y)
}
func (u UnionFind) Union(x, y int) {
	xr := u.Find(x)
	yr := u.Find(y)
	if xr == yr {
		return
	}
	if u.Size(yr) < u.Size(xr) {
		yr, xr = swap(yr, xr)
	}
	u.par[yr] += u.par[xr]
	u.par[xr] = yr
}
func (u UnionFind) Find(x int) int {
	if u.par[x] < 0 {
		return x
	}
	u.par[x] = u.Find(u.par[x])
	return u.par[x]
}
func NewUnionFind(N int) *UnionFind {
	u := new(UnionFind)
	u.par = make([]int, N)
	for i := range u.par {
		u.par[i] = -1
	}
	return u
}

type UnionFind struct{ par []int }

func ModPow() int64 {
	return int64(math.Max(1, 3))
}
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
