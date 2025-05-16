package test

import (
	"fmt"
	"strconv"
	"testing"

	"wgame_server/libray/internal"
)

type User struct {
	internal.ActorReceiver
}

func (that *User) Init() {
	that.ActorReceiver.Init(that)
}

func (that *User) Action(num int) string {
	return strconv.Itoa(num)
}

func TestNewReceiver(t *testing.T) {
	u := &User{}
	u.Init()
	ret, err := u.Invoker(1, "Action", 1)
	if err != nil {
		t.Error(err)
		return
	}
	fmt.Printf("%v", ret)

	_ = u
}
