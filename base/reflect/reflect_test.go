package reflect

import (
	"encoding/json"
	"fmt"
	"reflect"
	"testing"
	"unsafe"
)

type TestStruct struct {
	A int    `json:"a"`
	B string `json:"b"`
	c string `json:"c"`
	D *int   `json:"d"`

	Func func(int) error
}

// 通过反射修改非导出字段
func TestChangeNotExportFiled(t *testing.T) {
	var r TestStruct
	// 获取字段对象
	v := reflect.ValueOf(&r).Elem().FieldByName("c")
	// 构建指向该字段的可寻址（addressable）反射对象
	rv := reflect.NewAt(v.Type(), unsafe.Pointer(v.UnsafeAddr())).Elem()
	// 设置值
	fv := reflect.ValueOf("pibigstar")
	rv.Set(fv)

	t.Logf("%+v", r)
}

// 通过反射设置字段值
func TestChangeNotExportField2(t1 *testing.T) {
	var t TestStruct
	v := reflect.ValueOf(&t)

	fv := reflect.ValueOf(&t.A)
	fv.Elem().Set(reflect.ValueOf(123))
	t1.Logf("%+v", t)

	exportedFieldValue := v.Elem().Field(1)
	if !exportedFieldValue.CanSet() {
		t1.Log("can't set un export field")
	} else {
		exportedFieldValue.Set(reflect.ValueOf("456"))
		t1.Logf("%+v", t)
	}

	notExportFieldValue := v.Elem().Field(2)
	if !notExportFieldValue.CanSet() {
		t1.Log("can't set un export field")
	} else {
		notExportFieldValue.Set(reflect.ValueOf("123"))
		t1.Logf("%+v", t)
	}

	exportedPointField := v.Elem().Field(3)
	if !exportedPointField.CanSet() {
		t1.Log("can't set un export field")
	} else {
		num := 123
		exportedPointField.Set(reflect.ValueOf(&num))
		t1.Logf("%+v", t)
	}

}

// 根据反射判断字段类型
func TestInterface(t *testing.T) {
	var value interface{}
	value = "pibigstar"
	switch value.(type) {
	case string:
		v, ok := value.(string)
		if ok {
			t.Logf("String ==> %s \n", v)
		}
	case map[string]string:
		v, ok := value.(map[string]string)
		if ok {
			t.Logf("Map ==> %v \n", v)
		}
	default:
		bs, _ := json.Marshal(value)
		t.Logf("Others ==> %s \n", string(bs))
	}
}

// 反射基本操作
func TestReflect(t *testing.T) {

	var str = "hello world"

	v := reflect.ValueOf(str)
	// 获取值
	t.Log("value:", v)
	t.Log("value:", v.String())

	// 获取类型
	t.Log("type:", v.Type())
	t.Log("kind:", v.Kind())

	// 修改值
	// 判断是否可以修改
	canSet := v.CanSet() // 不可以直接修改
	t.Log("can set:", canSet)

	// 如果想修改其值，必须传递的是指针
	v = reflect.ValueOf(&str)
	v = v.Elem()
	v.SetString("new world")
	t.Log("value:", v)

	// 通过反射修改结构体
	num := 123
	test := TestStruct{A: 23, B: "hello world", c: "world2", D: &num}
	s := reflect.ValueOf(&test).Elem()

	s.Field(0).SetInt(77)
	s.Field(1).SetString("new world")
	t.Logf("%+v", test)

	typeOfT := reflect.TypeOf(test)
	for i := 0; i < typeOfT.NumField(); i++ {
		f := s.Field(i)
		// 无法获取未导出的字段的值
		t.Logf("%s: Type ==>%s Value==> %v \n", typeOfT.Field(i).Name, typeOfT.Field(i).Type, f.Interface())
	}

}

// 获取tag
func TestGetTag(t *testing.T) {
	s := TestStruct{}
	rt := reflect.TypeOf(s)
	for i := 0; i < rt.NumField(); i++ {
		f := rt.Field(i)
		t.Logf("%s: Tag: %v, Type: %v, Kind: %v\n", f.Name, f.Tag.Get("json"), f.Type, f.Type.Kind())
	}
}

// 处理不定数量的chan
func TestChan(t *testing.T) {
	ch1 := make(chan string, 10)
	ch2 := make(chan string, 10)
	ch3 := make(chan string, 10)
	cases := createCases(ch1, ch2, ch3)

	// 进行10次select
	for i := 0; i < 10; i++ {
		// 从cases里随机选择一个可用case
		chosen, recv, ok := reflect.Select(cases)
		// 是否是接收
		if recv.IsValid() && ok {
			t.Logf("chosen: %v, recv:%v\n", chosen, recv)
		} else {
			t.Log("send:", cases[chosen].Send)
		}
	}
}

func createCases(chs ...chan string) []reflect.SelectCase {
	var cases []reflect.SelectCase
	// create receiver case
	for _, ch := range chs {
		cases = append(cases, reflect.SelectCase{
			Dir:  reflect.SelectRecv,
			Chan: reflect.ValueOf(ch),
		})
	}
	// create send case
	for i, ch := range chs {
		cases = append(cases, reflect.SelectCase{
			Dir:  reflect.SelectSend,
			Chan: reflect.ValueOf(ch),
			Send: reflect.ValueOf(fmt.Sprintf("Hello: %d", i)), // 发送的send值
		})
	}
	return cases
}
