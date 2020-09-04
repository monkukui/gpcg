package lib

type UnionFind struct {
	par []int
}

func NewUnionFind(N int) *UnionFind {
	u := new(UnionFind)
	u.par = make([]int, N)
	for i := range u.par {
		u.par[i] = -1
	}
	return u
}

func (u UnionFind) Find(x int) int {
	if u.par[x] < 0 {
		return x
	}
	u.par[x] = u.Find(u.par[x])
	return u.par[x]
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

func (u UnionFind) Same(x, y int) bool {
	return u.Find(x) == u.Find(y)
}

func (u UnionFind) Size(x int) int {
	return -u.par[u.Find(x)]
}

// not exported but used
func swap(a int, b int) (int, int) {
	return b, a
}

// exported but not used
func UnUsedFunction() string {
	return "I am unused exported function"
}

// not exported and not used
func unUsedFunction() string {
	return "I am unused unexported function"
}

// token.VAR
var hoge = 10
var Hoge = 10

// token.CONST
const huga = 100
const Huga = 100

var (
	v1 = 1
	V2 = 2 // exported
	v3 = 3
	v4 = 4
)

var (
	UnUsedVar1 = 1
	UnUsedVar2 = 2
	UnUsedVar3 = 3
)

const (
	C1 = 1
	C2 = 2
	C3 = 3
	C4 = 4
)

const (
	UnUsedConst1 = 1
	UnUsedConst2 = 2
	UnUsedConst3 = 3
)
