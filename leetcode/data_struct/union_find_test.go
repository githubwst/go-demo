package data_struct

import "testing"

func TestUnionFind(t *testing.T) {
	uf := NewUnionFind(10)
	uf.merge(0, 1)
	uf.merge(2,3)
	uf.merge(5,6)
	t.Logf("%+v\n", uf)
	t.Log(uf.SetCount())
}
