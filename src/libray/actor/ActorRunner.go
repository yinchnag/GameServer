package actor

import (
	"sync"
	"sync/atomic"

	"wgame_server/libray/core"

	"github.com/Workiva/go-datastructures/queue"
)

func newActorRunner(ctx *ActorContext, options ...ActorConfigOption) *ActorRunner {
	config := ActorConfigure(options...)
	return newActorRunnerWithConfig(ctx, config)
}

func newActorRunnerWithConfig(ctx *ActorContext, config *actorConfig) *ActorRunner {
	runner := &ActorRunner{
		actorID:       0,
		queue:         queue.NewRingBuffer(config.Capacity),
		dropping:      config.Dropping,
		runCount:      0,
		overload:      0,
		overloadLimit: ACTORID_OVERLOAD,
		inGlobal:      1,
		ctx:           ctx,
		actorSystem:   ctx.ActorSystem,
	}
	return runner
}

type ActorRunner struct {
	actorID       uint32 // ActorID
	queue         *queue.RingBuffer
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

func (that *ActorRunner) pop() any {
	that.lock.Lock()
	defer that.lock.Unlock()
	count := that.queue.Len()
	if count > 0 {
		msg, _ := that.queue.Get()
		for count > that.overloadLimit {
			that.overload = count
			that.overloadLimit *= 2
			core.Logger.Errorf("[%s]actor overload len=%d limit=%d", that.ctx.Alias(), count, that.overloadLimit)
		}
		return msg
	}
	that.overloadLimit = ACTORID_OVERLOAD
	that.inGlobal = 0
	return nil
}

func (that *ActorRunner) recover(ctx *ActorContext) bool {
	if atomic.LoadInt32(&ctx.suspendNum) < 1 {
		return false
	}
	select {
	case msg := <-ctx.suspendChan:
		atomic.AddInt32(&ctx.suspendNum, -1)
		msg.Done()
		return true
	default:
		return false
	}
}

func (that *ActorRunner) run(ctx *ActorContext, weight int) {
	for i, count := 0, 1; i < count; i++ {
		if atomic.LoadInt32(&ctx.ref) != 0 {
			core.Logger.Infof("[%s] repeat in global", that.ctx.Alias())
			return
		}
		data := that.pop()
		if data == nil {
			return
		}
		if i == 0 && weight > 0 {
			count = that.length()
			count >>= weight
		}
		if that.overload != 0 {
			core.Logger.Errorf("actor overload,message queue length=%d", that.overload)
			that.overload = 0
			switch msg := data.(type) {
			case *ActorMessage:
				if atomic.LoadInt32(&msg.suspend) != 2 {
					msg.Result, msg.Err = ctx.Receiver.Invoker(msg.UID, msg.FunName, msg.Args...)
					msg.Suspend(true)
				} else {
					srcAlias := core.TernaryF(msg.SourceCtx != nil, func() string { return msg.SourceCtx.Alias() }, "nil")
					core.Logger.Debugf("[%s] suspend timeout UID=%d funName=%s等待返回 src%s", msg.Alias, msg.UID, msg.FunName, srcAlias)
					msg.Free()
				}
			case *actorSuspend:
				that.recover(msg.ctx)
				return // 唤醒挂起协程并结束当前协程
			default:
				ctx.Receiver.Receive(msg)
			}
		}
	}
	that.actorSystem.Push(ctx)
}
