package main

import (
	"bufio"
	"fmt"
	"os"
)

func generated_monkukui_swap(a int, b int) (int, int) {
	return b, a
}
func (u generated_monkukui_UnionFind) Size(x int) int {
	return -u.par[u.Find(x)]
}
func (u generated_monkukui_UnionFind) Union(x, y int) {
	xr := u.Find(x)
	yr := u.Find(y)
	if xr == yr {
		return
	}
	if u.Size(yr) < u.Size(xr) {
		yr, xr = generated_monkukui_swap(yr, xr)
	}
	u.par[yr] += u.par[xr]
	u.par[xr] = yr
}
func (u generated_monkukui_UnionFind) Find(x int) int {
	if u.par[x] < 0 {
		return x
	}
	u.par[x] = u.Find(u.par[x])
	return u.par[x]
}
func generated_monkukui_NewUnionFind(N int) *generated_monkukui_UnionFind {
	u := new(generated_monkukui_UnionFind)
	u.par = make([]int, N)
	for i := range u.par {
		u.par[i] = -1
	}
	return u
}

type generated_monkukui_UnionFind struct{ par []int }

func main() {
	r := bufio.NewReader(os.Stdin)
	w := bufio.NewWriter(os.Stdout)
	defer w.Flush()
	var n, m int
	fmt.Fscan(r, &n, &m)
	uf := generated_monkukui_NewUnionFind(n)
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
