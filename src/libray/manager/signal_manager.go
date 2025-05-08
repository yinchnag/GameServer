package manager

import (
	"sync"

	"wgame_server/libray/core"
)

const (
	SIGNAL_START_SERVER   = 1 // 启动服务器
	SIGNAL_CONNECT_SERVER = 2 // 连接服务器

	SIGNAL_PLAYER_NEW = 3 // 新建角色
)

var (
	signalManager     *SignalManager
	signalManagerOnce sync.Once
)

func GetSignalManager() *SignalManager {
	signalManagerOnce.Do(func() {
		signalManager = &SignalManager{}
		signalManager.Init()
	})
	return signalManager
}

type SignalManager struct {
	core.Delegate
}

func (that *SignalManager) Init() {
	that.Delegate.Init()
}

func (that *SignalManager) AddListener(event int, cb interface{}) {
	that.Delegate.AddListener(event, cb)
}

// 移除事件
func (that *SignalManager) RemoveListener(callback interface{}) {
	that.Delegate.RemoveListener(1, callback)
}

// 事件通知
func (that *SignalManager) Notify(event int, args ...interface{}) {
	that.Delegate.Notify(event, args...)
}
