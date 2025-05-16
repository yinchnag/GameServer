package actor

import (
	"sync/atomic"
	"time"

	"wgame_server/libray/core"
)

type ActorContext struct {
	ActorSystem *ActorSystem
	Receiver    IRceiver
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
	switch msg := data.(type) {
	case *ActorMessage:
		if source == that.runner.actorID { // 本地消息
			msg.Result, msg.Err = that.Receiver.Invoker(msg.UID, msg.FunName, msg.Args...)
		} else if wait || that.Receiver.GetNumOut(msg.FunName) > 0 {
			msg.SourceCtx = that.ActorSystem.FindActorByID(source)
			msg.Add(1)
			if msg.SourceCtx != nil {
				atomic.AddInt32(&msg.SourceCtx.suspendNum, 1)
				worker := that.ActorSystem.workers[msg.SourceCtx.workerIdx]
				cursor := atomic.AddUint32(&worker.cursor, 1)
				go that.ActorSystem.runRecvLoop(worker, cursor)

				// 冻结状态允许重新加入调度队列
				msg.SourceCtx.runner.lock.Lock()
				msg.SourceCtx.runner.inGlobal = 0
				msg.SourceCtx.runner.lock.Unlock()
			}
			that.runner.send(msg)
			msg.timer = time.AfterFunc(time.Millisecond*msg.timeout, func() {
				if msg.Suspend(false) {
					during := time.Now().UnixMilli() - msg.startTime
					srcAlias := core.TernaryF(msg.SourceCtx != nil, func() string { return msg.SourceCtx.Alias() }, "nil")
					core.Logger.Errorf("[%s] ActorContext.Send UID=%d funcName=%s超时,总耗时%dms src=%s", msg.Alias, msg.UID, msg.FunName, during, srcAlias)
				}
			})
			msg.Wait()
		} else {
			that.runner.send(msg)
		}
	default:
		that.runner.send(msg)
	}
}
