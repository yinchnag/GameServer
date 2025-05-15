package internal

import (
	"wgame_server/libray/interfaces"
)

type ActorContext struct {
	ActorSystem *ActorSystem
	Receiver    interfaces.IRceiver
	runner      *ActorRunner       // actor对象
	workerIdx   int                // 所在工作器索引
	suspendChan chan *ActorMessage // 挂起消息通道
	suspendNum  int32              // 挂起消息数量
	Profile     bool               // 是否开启性能分析
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
