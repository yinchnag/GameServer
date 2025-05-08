package module

import (
	"bytes"
	"reflect"
	"runtime"
	"strconv"
	"sync"
	"sync/atomic"
	"time"

	"wgame_server/libray/core"

	"github.com/timandy/routine"
)

const CHAN_TASK_SIZE = 5000 // 玩家管道任务长度

const BenchMark_TimeOut = 1000 // 压测耗时较长标准
// ----------------------------------------------------
// 通道任务池
var chanTaskPool = sync.Pool{
	New: func() interface{} {
		return new(ChanTask)
	},
}

// 从池中分配通道任务
func NewChanTask(modName string, funName string, args []interface{}) *ChanTask {
	task := chanTaskPool.Get().(*ChanTask)
	task.Module = modName
	task.FunName = funName
	task.Args = args
	atomic.StoreInt32(&task.Failed, 0)
	return task
}

// 协程buff池
var littleBuf = sync.Pool{
	New: func() interface{} {
		buf := make([]byte, 64)
		return &buf
	},
}

// 获取协程id(调用10000次50ms)
func GetGoroutineID() int64 {
	goroutineId := routine.Goid()
	if goroutineId == 0 {
		bp := littleBuf.Get().(*[]byte)
		defer littleBuf.Put(bp)
		buff := *bp
		runtime.Stack(buff, false)
		buff = bytes.TrimPrefix(buff, []byte("goroutine "))
		buff = buff[:bytes.IndexByte(buff, ' ')]
		goroutineId, _ = strconv.ParseInt(string(buff), 10, 64)
	}
	return goroutineId
}

type ChanTask struct {
	Module    string          // 模块名称
	FunName   string          // 函数名称
	Result    []reflect.Value // 执行结果
	Args      []interface{}   // 参数
	Err       error           // 错误
	Stack     []string        // 当前堆栈
	Failed    int32           // 状态
	LockCount int32           // 计次
	Future    *core.Future    // 异步实例

	sync.WaitGroup
}

func (that *ChanTask) SetTimeout(timeout time.Duration) {
	if that.Future != nil {
		that.Future.SetTimeout(timeout)
	}
}

// 释放通道任务
func (that *ChanTask) Free() {
	chanTaskPool.Put(that)
}

// 追加
func (that *ChanTask) Add(delta int) {
	atomic.AddInt32(&that.LockCount, 1)
	that.WaitGroup.Add(delta)
}

// 完成
func (that *ChanTask) Done() {
	if atomic.LoadInt32(&that.LockCount) > 0 {
		atomic.AddInt32(&that.LockCount, -1)
		that.WaitGroup.Done()
	}
}

// 是否完成
func (that *ChanTask) IsFinish() bool {
	return atomic.LoadInt32(&that.LockCount) <= 0
}
