package internal

import (
	"errors"
	"reflect"
	"sync/atomic"
	"time"

	"wgame_server/libray/core"
	"wgame_server/libray/interfaces"
)

type (
	ActorProducer func() interfaces.IRceiver

	ActorSystem struct {
		queue      chan *ActorContext // 任务列表
		capacity   int                // 任务队列容量
		throughput int                // 负载
		stopper    chan struct{}      // 停止信号
		workers    []*actorWorker     // 工作上下文
		profile    bool               // 是否调试
		startTime  int64              // 启动时间
		storage    *ActorStorage      // 存储器
	}
)

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

func (that *ActorSystem) Start() {
	that.startTime = time.Now().Unix()
	for i := 0; i < that.throughput; i++ {
		go that.runRecvLoop(that.workers[i], 1)
	}
}

func (that *ActorSystem) runRecvLoop(worker *actorWorker, cursor uint32) {
	defer func() {
		err := recover()
		if err != nil {
			core.Logger.Errorln(err)
		}
	}()

	for ctx := range that.queue {
		if that.IsStoppend() {
			break
		}
		if ref := atomic.AddInt32(&ctx.ref, -1); ref != 0 {
			core.Logger.Warnf("actor ref=%d error", ref) // 防止多次加入引起协程一场
			continue
		}
		that.processActor(worker, ctx)
		if cursor != worker.cursor {
			break // 解冻的线程直接退出
		}
	}
}

func (that *ActorSystem) processActor(worker *actorWorker, ctx *ActorContext) {
	startTime := time.Now().UnixMilli()
	worker.context = ctx
	ctx.workerIdx = worker.index // 标记工作器方便后续重启工作器协程
	ctx.runner.run(ctx, RUNNER_WEIGHT[worker.index%32])
	if that.profile {
		core.Logger.Infof("[%s]好事%dms,执行任务次数%d", ctx.Receiver.GetName(), (time.Now().UnixMilli() - startTime), ctx.runner.runCount)
	}
	worker.context = nil
}

func (that *ActorSystem) Shutdown() {
	close(that.stopper)
}

func (that *ActorSystem) IsStoppend() bool {
	select {
	case <-that.stopper:
		return true
	default:
		return false
	}
}

// 分配actor
func (that *ActorSystem) AllocActor(alias string, producer ActorProducer, options ...ActorConfigOption) *ActorContext {
	recevier := producer()
	ctx := &ActorContext{
		ActorSystem: that,
		Receiver:    recevier,
		Profile:     that.profile,
		suspendChan: make(chan *ActorMessage, ACTORID_SUSPEND),
	}
	ctx.runner = newActorRunner(ctx, options...)
	that.storage.register(ctx)
	recevier.Init(ctx)
	recevier.(interfaces.IModule).Start()
	return ctx
}

func (that *ActorSystem) FreeActor(ctx *ActorContext) {
	ctx.Receiver.(interfaces.IModule).Destory()
}

// 根据ID获得上下文
func (that *ActorSystem) FindActorByID(actorID uint32) *ActorContext {
	return that.storage.findActorByID(actorID)
}

// 根据别名获得上下文
func (that *ActorSystem) FindActorByAlias(alias string) *ActorContext {
	return that.storage.findActorByAlias(alias)
}

// 别名转actorID
func (that *ActorSystem) ToActorID(alias string) uint32 {
	return that.storage.toActorID(alias)
}

// 遍历actor
func (that *ActorSystem) ForEach(cb func(*ActorContext)) {
	that.storage.ForEach(cb)
}

// 调用导出接口
// 带反参的调用务必加上来源上下文
// uid=-1 actor 0 mgr > 0 mod
func (that *ActorSystem) ModInvokeSafe(sourceCtx *ActorContext, uid int64, actorAlias string, funName string, args ...any) ([]reflect.Value, error) {
	actor := that.FindActorByAlias(actorAlias)
	if actor == nil {
		return nil, errors.New("actor not found")
	}
	result, err := actor.Receiver.Invoker(uid, funName, args...)
	return result, err
}

func (that *ActorSystem) ModInvoke(sourceCtx *ActorContext, uid int64, actorAlias string, funName string, args ...any) ([]reflect.Value, error) {
	actor := that.FindActorByAlias(actorAlias)
	if actor == nil {
		return nil, errors.New("actor not found")
	}
	msg := NewActorMessage(uid, actorAlias, funName, 1000, args...)
	source := core.TernaryF(sourceCtx != nil, func() uint32 { return sourceCtx.ActorID() }, 0)
	actor.Send(source, msg, true)
	if atomic.LoadInt32(&msg.suspend) == 2 {
		defer msg.Free()
	}
	return msg.Result, msg.Err
}
