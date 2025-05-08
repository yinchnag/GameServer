package test

import (
	"testing"
	"wgame_server/libray/manager"
)

func TestSignal(t *testing.T) {
	signal := manager.GetSignalManager()
	signal.AddListener(manager.SIGNAL_START_SERVER, func(a, b int) {
		t.Log("start server1")
	})
	signal.AddListener(manager.SIGNAL_START_SERVER, func(a, b int) {
		t.Log("start server2")
	})
	signal.Notify(manager.SIGNAL_START_SERVER, 12, 1)
}
