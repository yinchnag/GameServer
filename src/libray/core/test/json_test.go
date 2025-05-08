package test

import (
	"testing"
	"wgame_server/libray/core"
)

type User struct {
	Name string `json:"Name"`
	Age  int    `json:"Age"`
}

func TestJsonMarshal(t *testing.T) {
	u := &User{Name: "yinchang", Age: 30}
	bt, err := core.Marshal(u)
	if len(bt) == 0 {
		t.Log(err)
		t.Fail()
	}
}

func TestJsonUnmarshal(t *testing.T) {
	u := &User{Name: "yinchang", Age: 30}
	bt, _ := core.Marshal(u)
	u1 := &User{}
	err := core.Unmarshal(bt, u1)
	if err != nil || u1.Age != 30 {
		t.Fail()
	}
}
