package test

import (
	"testing"
	"wgame_server/libray/core"
)

func TestTimestamp(t *testing.T) {
	tk := core.ServerTime()
	tk1 := core.TimestampToTime(tk.Unix() + 60*60)
	if tk1.Unix()-tk.Unix() != 60*60 {
		t.Fail()
	}
}
