package reflect

import (
	"fmt"
	"reflect"
	"testing"
)

type TestStructWithFunc struct {
	*TestStruct
	f func(string) error
}

func (t TestStructWithFunc) SayHello(name string) string {
	return fmt.Sprint("hello " + name)
}

func (t *TestStructWithFunc) SayFuck(name string) string {
	return fmt.Sprint("fuck " + name)
}

func (t *TestStructWithFunc) sayFuck(name string) {
	fmt.Println("fuck" + name)
}

func TestReflectMethod(t *testing.T) {
	var r TestStructWithFunc
	rt := reflect.TypeOf(r)
	for i := 0; i < rt.NumField(); i++ {
		field := rt.Field(i)
		t.Logf("%s: Type: %v, Kind: %v", field.Name, field.Type, field.Type.Kind())
	}

	fmt.Println()
	rt = reflect.TypeOf(&r)
	// rt.NumMethod() 只统计导出的方法
	for i := 0; i < rt.NumMethod(); i++ {
		f := rt.Method(i)
		t.Logf("%s: Type: %v, Kind: %v", f.Name, f.Type, f.Type.Kind())
	}
	sayFuckMethod, ok := rt.MethodByName("sayFuck")
	if !ok {
		t.Log("can't access un export method event if by method name.")
	} else {
		t.Logf("%s: Type: %v, Kind: %v", sayFuckMethod.Name, sayFuckMethod.Type, sayFuckMethod.Type.Kind())
	}

}

func TestCallMethod(t *testing.T) {
	var r TestStructWithFunc
	// 指针接收器有两个方法：
	// 0: SayFuck
	// 1: SayHello
	// 非指针接收器仅有一个方法： 0：SayHello

	// 可以通过指针调用struct方法；不可以通过struct 调用指针方法
	rt := reflect.TypeOf(&r)
	// 传入的接收器类型必须和Type的类型一致
	callReturn := rt.Method(0).Func.Call([]reflect.Value{reflect.ValueOf(&r), reflect.ValueOf("wst")})
	t.Log(callReturn)

	// 指针value和非指针value都可以正常调用方法
	rv := reflect.ValueOf(&r)
	callReturn = rv.Method(1).Call([]reflect.Value{reflect.ValueOf("wst")})
	t.Log(callReturn)
}
