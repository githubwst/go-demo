package cmp

import "testing"

type User struct {
	Name     string
	Password string
}
type User2 struct {
	Name     string
	Password string
}

func TestDiff(t *testing.T) {
	u1 := &User{
		Name:     "派大星",
		Password: "123456",
	}

	u2 := &User{
		Name:     "海绵宝宝",
		Password: "123456",
	}

	u3 := &User2{
		Name:     "海绵宝宝",
		Password: "123456",
	}

	diff := Diff(u1, u2)
	t.Log(diff)

	// 只能diff相同类型struct的字段value的不同
	diff2 := Diff(u1, u3)
	t.Log(diff2)
}
