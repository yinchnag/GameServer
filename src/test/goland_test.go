package test

import (
	"fmt"
	"reflect"
	"testing"
)

func Add(a, b int) int {
	return a + b
}

func Add1(a, b int) int {
	return a + b
}

func equip(ab, cb reflect.Value) bool {
	return ab == cb
}

func TestReflect(t *testing.T) {
	if !equip(reflect.ValueOf(Add), reflect.ValueOf(Add)) {
		t.Fail()
	}
}

func TestGetFuncProgm(t *testing.T) {
	ref := reflect.TypeOf(Add)
	for i := 0; i < ref.NumIn(); i++ {
		fmt.Println("参数:", ref.In(i))
	}
}
