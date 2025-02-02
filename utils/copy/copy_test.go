package copy

import (
	"testing"
)

func BenchmarkDeepCopy(b *testing.B) {
	type User struct {
		Name string
	}
	user1 := &User{Name: "pibigstar"}
	for i := 0; i < b.N; i++ {
		var user2 User
		_ = DeepCopy(user1, &user2)
	}
}

func BenchmarkDeepCopy2(b *testing.B) {
	type User struct {
		Name string
	}
	user1 := &User{Name: "pibigstar"}
	for i := 0; i < b.N; i++ {
		var user2 User
		_ = Copy(user1, &user2)
	}
}

func TestDeepCopy(t *testing.T) {
	type User struct {
		Name string
	}
	user1 := &User{Name: "pibigstar"}

	var user2 User
	err := DeepCopy(user1, &user2)
	if err != nil {
		t.Error(err)
	}
	t.Log(user2)

	var user3 User
	err = Copy(user1, &user3)
	if err != nil {
		t.Error(err)
	}
	t.Log(user3)
}

type TypeStruct1 struct {
	Field1 string
	Field2 string
	Field3 TypeStruct2
	Field4 *TypeStruct2
	Field5 []*TypeStruct2
	Field6 []TypeStruct2
	Field7 []*TypeStruct2
	Field8 []TypeStruct2
}

type TypeStruct2 struct {
	Field1 int
	Field2 string
}

type TypeStruct3 struct {
	Field1 interface{}
	Field2 string
	Field3 TypeStruct4
	Field4 *TypeStruct4
	Field5 []*TypeStruct4
	Field6 []*TypeStruct4
	Field7 []TypeStruct4
	Field8 []TypeStruct4
}

type TypeStruct4 struct {
	field1 int
	Field2 string
}

func (t *TypeStruct4) Field1(i int) {
	t.field1 = i
}

func TestCopyDifferentFieldType(t *testing.T) {
	ts := &TypeStruct1{
		Field1: "str1",
		Field2: "str2",
	}
	ts2 := &TypeStruct2{}

	Copy(ts2, ts)

	if ts2.Field2 != ts.Field2 || ts2.Field1 != 0 {
		t.Errorf("Should be able to copy from ts to ts2")
	}
}

func TestCopyDifferentTypeMethod(t *testing.T) {
	ts := &TypeStruct1{
		Field1: "str1",
		Field2: "str2",
	}
	ts4 := &TypeStruct4{}

	Copy(&ts4, &ts)
	t.Logf("%+v", ts4)
	if ts4.Field2 != ts.Field2 || ts4.field1 != 0 {
		t.Errorf("Should be able to copy from ts to ts4")
	}
}

func TestAssignableType(t *testing.T) {
	ts := &TypeStruct1{
		Field1: "str1",
		Field2: "str2",
		Field3: TypeStruct2{
			Field1: 666,
			Field2: "str2",
		},
		Field4: &TypeStruct2{
			Field1: 666,
			Field2: "str2",
		},
		Field5: []*TypeStruct2{
			{
				Field1: 666,
				Field2: "str2",
			},
		},
		Field6: []TypeStruct2{
			{
				Field1: 666,
				Field2: "str2",
			},
		},
		Field7: []*TypeStruct2{
			{
				Field1: 666,
				Field2: "str2",
			},
		},
	}

	ts3 := &TypeStruct3{}

	Copy(&ts3, &ts)

	if v, ok := ts3.Field1.(string); !ok {
		t.Error("Assign to interface{} type did not succeed")
	} else if v != "str1" {
		t.Error("String haven't been copied correctly")
	}

	if ts3.Field2 != ts.Field2 {
		t.Errorf("Field2 should be copied")
	}

	checkType2WithType4(ts.Field3, ts3.Field3, t, "Field3")
	checkType2WithType4(*ts.Field4, *ts3.Field4, t, "Field4")

	for idx, f := range ts.Field5 {
		checkType2WithType4(*f, *(ts3.Field5[idx]), t, "Field5")
	}

	for idx, f := range ts.Field6 {
		checkType2WithType4(f, *(ts3.Field6[idx]), t, "Field6")
	}

	for idx, f := range ts.Field7 {
		checkType2WithType4(*f, ts3.Field7[idx], t, "Field7")
	}

	for idx, f := range ts.Field8 {
		checkType2WithType4(f, ts3.Field8[idx], t, "Field8")
	}
}

func checkType2WithType4(t2 TypeStruct2, t4 TypeStruct4, t *testing.T, testCase string) {
	if t2.Field1 != t4.field1 || t2.Field2 != t4.Field2 {
		t.Errorf("%v: type data_struct 4 and type data_struct 2 is not equal", testCase)
	}
}
