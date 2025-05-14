package internal

import "sync/atomic"

type ActorSystem struct {
	queue    chan *ActorContext // 任务列表
	capacity int                // 任务队列容量
}

func (that *ActorSystem) Length() int {
	return len(that.queue)
}

func (that *ActorSystem) Push(ctx *ActorContext) {
	if that.Length() >= that.capacity-1 {
		go that.doPush(ctx)
	} else {
		that.doPush(ctx)
	}
}

// 加入Actor单元
func (that *ActorSystem) doPush(ctx *ActorContext) {
	if atomic.CompareAndSwapInt32(&ctx.ref, 0, 1) {
		that.queue <- ctx
	}
}
