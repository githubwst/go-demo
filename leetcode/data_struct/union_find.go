package data_struct

type UnionFind struct {
	parent, rank []int
	cnt          int
}

func NewUnionFind(n int) *UnionFind {
	parent := make([]int, n)
	for i := range parent {
		parent[i] = i
	}
	return &UnionFind{parent, make([]int, n), n}
}

func (uf *UnionFind) find(x int) int {
	if uf.parent[x] != x {
		uf.parent[x] = uf.find(uf.parent[x])
	}
	return uf.parent[x]
}

func (uf *UnionFind) merge(x, y int) {
	x, y = uf.find(x), uf.find(y)
	if x == y {
		return
	}
	// rank 小的挂在 rank 大的下面
	rankX := uf.rank[x]
	rankY := uf.rank[y]
	if rankX < rankY {
		y += x
		x = y - x
		y = y - x
	}
	uf.parent[y] = x
	if rankX == rankY {
		uf.rank[x]++
	}
	uf.cnt--

	// rank 小的挂在 rank 大的下面
	//if uf.rank[x] > uf.rank[y] {
	//	uf.parent[y] = x
	//} else if uf.rank[x] < uf.rank[y] {
	//	uf.parent[x] = y
	//} else {
	//	uf.parent[y] = x
	//	uf.rank[x]++
	//}
}

func (uf *UnionFind) SetCount() int {
	return uf.cnt
}
