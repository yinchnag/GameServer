package test

import (
	"server/src/library/internal/base"
	"testing"
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
