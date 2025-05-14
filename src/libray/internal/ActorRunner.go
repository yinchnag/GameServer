package internal

import (
	"sync"

	"wgame_server/libray/core"

	"github.com/Workiva/go-datastructures/queue"
)

type ActorRunner struct {
	actorID       uint32 // ActorID
	queue         queue.RingBuffer
	dropping      bool          // 是否丢去溢出消息
	runCount      uint64        // 执行次数
	overload      uint64        // 超载数量
	overloadLimit uint64        // 超载上线
	inGlobal      int32         // 是否加入全局队列
	ctx           *ActorContext // 上下文
	actorSystem   *ActorSystem  // ActorSystem
	lock          sync.Mutex    // 异步锁
}

func (that *ActorRunner) length() int {
	return int(that.queue.Len())
}

// 发送消息
func (that *ActorRunner) send(data any) {
	switch msg := data.(type) {
	case *ActorMessage:
		if msg.Alias != that.ctx.Alias() {
			core.Logger.Errorf("actor message alias error, expect %s, but %s", that.ctx.Alias(), msg.Alias)
			return
		}
	}
	that.lock.Lock()
	defer that.lock.Unlock()
	if that.dropping {
		if that.queue.Len() > 0 && that.queue.Cap()-1 == that.queue.Len() {
			_, _ = that.queue.Get()
		}
	}
	err := that.queue.Put(data)
	if err != nil {
		core.Logger.Errorf("actor queue put error,%v", err)
	}
	if that.inGlobal == 0 {
		that.inGlobal = 1
		that.actorSystem.Push(that.ctx)
	}
}
