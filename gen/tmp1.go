package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
)

const (
	generated_alib_C1 = 1
	generated_alib_C2 = 2
	generated_alib_C3 = 3
	generated_alib_C4 = 4
)

var (
	generated_alib_v1 = 1
	generated_alib_V2 = 2
	generated_alib_v3 = 3
	generated_alib_v4 = 4
)

const generated_alib_Huga = 100

var generated_alib_Hoge = 10

func generated_alib_swap(a int, b int) (int, int) {
	return b, a
}
func (u generated_alib_UnionFind) Size(x int) int {
	return -u.par[u.Find(x)]
}
func (u generated_alib_UnionFind) Union(x, y int) {
	xr := u.Find(x)
	yr := u.Find(y)
	if xr == yr {
		return
	}
	if u.Size(yr) < u.Size(xr) {
		yr, xr = generated_alib_swap(yr, xr)
	}
	u.par[yr] += u.par[xr]
	u.par[xr] = yr
}
func (u generated_alib_UnionFind) Find(x int) int {
	if u.par[x] < 0 {
		return x
	}
	u.par[x] = u.Find(u.par[x])
	return u.par[x]
}
func generated_alib_NewUnionFind(N int) *generated_alib_UnionFind {
	u := new(generated_alib_UnionFind)
	u.par = make([]int, N)
	for i := range u.par {
		u.par[i] = -1
	}
	return u
}

type generated_alib_UnionFind struct{ par []int }

func main() {
	r := bufio.NewReader(os.Stdin)
	w := bufio.NewWriter(os.Stdout)
	defer w.Flush()
	fmt.Fprintln(w, "alib.V2 = ", generated_alib_V2)
	fmt.Fprintln(w, "alib.C4 = ", generated_alib_C4)
	fmt.Fprintln(w, "alib.Hoge = ", generated_alib_Hoge)
	fmt.Fprintln(w, "alib.Huga = ", generated_alib_Huga)
	var n, m int
	fmt.Fscan(r, &n, &m)
	uf := generated_alib_NewUnionFind(n)
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
	fmt.Fprintln(w, "ans = ", ans)
}
