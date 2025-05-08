package test

import (
	"testing"
	"wgame_server/libray/core"
)

func TestGetPath(t *testing.T) {
	path := core.GetExecutableAbsPath()
	if path == "" {
		t.Fail()
	}
}
