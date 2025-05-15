package internal

import (
	"reflect"
	"sync"
	"sync/atomic"
	"time"

	"wgame_server/libray/core"
)

var actorMessagePool = &sync.Pool{
	New: func() any {
		return &ActorMessage{}
	},
}

func NewActorMessage(uid int64, alias string, funName string, timeout time.Duration, args ...any) *ActorMessage {
	task := actorMessagePool.Get().(*ActorMessage)
	task.UID = uid
	task.SourceCtx = nil
	task.Alias = alias
	task.FunName = funName
	task.Args = args
	task.Result = nil
	task.Err = nil
	task.startTime = time.Now().UnixMilli()
	task.timer = nil
	if core.IsDebug {
		task.timeout = time.Hour
	} else {
		task.timeout = timeout
	}
	atomic.StoreInt32(&task.lockCount, 0)
	atomic.StoreInt32(&task.suspend, 0)
	atomic.StoreInt32(&task.dispose, 0)
	return task
}

type ActorMessage struct {
	UID       int64           // 玩家ID
	SourceCtx *ActorContext   // 来源上下文
	Alias     string          // 模块名称
	FunName   string          // 函数名称
	Args      []any           // 参数名称
	Result    []reflect.Value // 执行结果
	Err       error           // 错误
	lockCount int32           // 计次
	suspend   int32           // 是否挂起 0非挂起 1挂起 2挂起后等待回收
	startTime int64           // 开始时间
	timer     *time.Timer     // 倒计时
	timeout   time.Duration   // 超时时间
	dispose   int32           // 是否销毁
	sync.WaitGroup
}

func (that *ActorMessage) Free() {
	that.Args = nil
	that.Result = nil
	that.Err = nil
	atomic.StoreInt32(&that.lockCount, 0)
	atomic.StoreInt32(&that.suspend, 0)
	atomic.StoreInt32(&that.dispose, 1)
}

func (that *ActorMessage) Add(delta int) {
	atomic.AddInt32(&that.lockCount, 1)
	that.WaitGroup.Add(delta)
	atomic.StoreInt32(&that.lockCount, 1)
}

func (that *ActorMessage) Done() {
	if atomic.LoadInt32(&that.lockCount) > 0 {
		atomic.AddInt32(&that.lockCount, -1)
		that.WaitGroup.Done()
	} else {
		srcAlias := core.TernaryF(that.SourceCtx != nil, func() string { return that.SourceCtx.Alias() }, "nil")
		core.Logger.Errorf("[%s]message Done err UID=%d FunName=%s等待返回 src=%s", that.Alias, that.UID, that.FunName, srcAlias)
	}
	atomic.StoreInt32(&that.suspend, 0)
}

func (that *ActorMessage) Suspend(free bool) bool {
	if atomic.LoadInt32(&that.suspend) == 1 {
		if that.timer != nil {
			that.timer.Stop()
			that.timer = nil
		}
		if that.SourceCtx != nil {
			that.SourceCtx.suspendChan <- that
			that.SourceCtx.Send(that.SourceCtx.ActorID(), &actorSuspend{ctx: that.SourceCtx}, false)
		} else {
			that.Done()
		}
		atomic.StoreInt32(&that.suspend, 2) // 等待回收
		return true
	}
	if free {
		that.Free()
	}
	return false
}

// 是否完成
func (that *ActorMessage) IsFinish() bool {
	return atomic.LoadInt32(&that.lockCount) <= 0
}
