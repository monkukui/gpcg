package main

import (
	"fmt"
	"math"
	"strconv"
)

func ModPow() float64 {
	fmt.Println("hello mod pow")
	return math.Max(1.0, 2.0)
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

func ModInv() float64 {
	fmt.Println("hello mod inv")
	fmt.Println(strconv.Atoi("122"))
	return math.Min(1.0, 2.0)
}
func main() {
	var n, m int
	fmt.Scan(&n, &m)
	uf := NewUnionFind(n)
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
