package core

import "reflect"

// 判空
func IsNil(c interface{}) bool {
	if btype := reflect.TypeOf(c); btype == nil {
		return true
	}
	vi := reflect.ValueOf(c)
	switch vi.Kind() {
	case reflect.Func,
		reflect.Chan,
		reflect.Map,
		reflect.Pointer,
		reflect.UnsafePointer,
		reflect.Interface,
		reflect.Slice:
		return vi.IsNil()
	default:
		return false
	}
}

// 反射拷贝
// goos: windows
// goarch: amd64
// pkg: wgame_server/libray/core/test
// cpu: 11th Gen Intel(R) Core(TM) i7-11700 @ 2.50GHz
// 163.3 ns/op
func HF_ReflectCopy(src interface{}) interface{} {
	if reflect.TypeOf(src).Kind() == reflect.Ptr {
		dst := reflect.New(reflect.ValueOf(src).Elem().Type())
		dst.Elem().Set(reflect.ValueOf(src).Elem())
		return dst.Interface()
	} else {
		dst := reflect.New(reflect.TypeOf(src))
		dst.Elem().Set(reflect.ValueOf(src))
		return dst.Elem().Interface()
	}
}

// 反射创建
func HF_ReflectNew(src interface{}) interface{} {
	if reflect.TypeOf(src).Kind() == reflect.Ptr {
		return reflect.New(reflect.ValueOf(src).Elem().Type()).Interface()
	} else {
		return reflect.New(reflect.TypeOf(src)).Elem().Interface()
	}
}
