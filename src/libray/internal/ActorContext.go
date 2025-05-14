package internal

import (
	"wgame_server/libray/interfaces"
)

type ActorContext struct {
	ActorSystem *ActorSystem
	Receiver    interfaces.IRceiver
	runner      *ActorRunner       // actor对象
	suspendChan chan *ActorMessage // 挂起消息通道
	ref         int32              // 引用计数
}

// 获取别名
func (that *ActorContext) Alias() string {
	return that.Receiver.GetName()
}

func (that ActorContext) ActorID() uint32 {
	return that.runner.actorID
}

func (that *ActorContext) Send(source uint32, data any, wait bool) {
}
