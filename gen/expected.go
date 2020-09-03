package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
)

func _lib_unUsedFunction() string {
	return "I am unused unexported function"
}
func _lib_UnUsedFunction() string {
	return "I am unused exported function"
}
func _lib_swap(a int, b int) (int, int) {
	return b, a
}
func (u _lib_UnionFind) _lib_Size(x int) int {
	return -u.par[u._lib_Find(x)]
}
func (u _lib_UnionFind) _lib_Same(x, y int) bool {
	return u._lib_Find(x) == u._lib_Find(y)
}
func (u _lib_UnionFind) _lib_Union(x, y int) {
	xr := u._lib_Find(x)
	yr := u._lib_Find(y)
	if xr == yr {
		return
	}
	if u._lib_Size(yr) < u._lib_Size(xr) {
		yr, xr = _lib_swap(yr, xr)
	}
	u.par[yr] += u.par[xr]
	u.par[xr] = yr
}
func (u _lib_UnionFind) _lib_Find(x int) int {
	if u.par[x] < 0 {
		return x
	}
	u.par[x] = u._lib_Find(u.par[x])
	return u.par[x]
}
func _lib_NewUnionFind(N int) *_lib_UnionFind {
	u := new(_lib_UnionFind)
	u.par = make([]int, N)
	for i := range u.par {
		u.par[i] = -1
	}
	return u
}

type _lib_UnionFind struct{ par []int }

func _lib_ModPow() int64 {
	return int64(math.Max(1, 3))
}
func main() {
	r := bufio.NewReader(os.Stdin)
	w := bufio.NewWriter(os.Stdout)
	defer w.Flush()
	var n, m int
	fmt.Fscan(r, &n, &m)
	uf := _lib_NewUnionFind(n)
	for i := 0; i < m; i++ {
		var a, b int
		fmt.Fscan(r, &a, &b)
		a--
		b--
		uf._lib_Union(a, b)
	}
	ans := 0
	for i := 0; i < n; i++ {
		if ans < uf._lib_Size(i) {
			ans = uf._lib_Size(i)
		}
	}
	fmt.Fprintln(w, ans)
}

func swap(a int, b int) (int, int) {
	return 10 * b, 10 * a
}
