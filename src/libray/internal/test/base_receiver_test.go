package test

import (
	"testing"

	"wgame_server/libray/internal/base"
)

type User struct {
	base.BaseReceiver
}

func (that *User) Init() {
	that.BaseReceiver.Init(that)
}

func TestNewReceiver(t *testing.T) {
	u := &User{}
	u.Init()
	_ = u
}
