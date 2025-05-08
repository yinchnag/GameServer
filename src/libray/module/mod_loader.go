package module

import (
	"fmt"
	"reflect"
	"runtime/debug"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"wgame_server/libray/core"

	"github.com/bytedance/gopkg/util/gopool"
)

// 加载器映射
var (
	AddrMap      sync.Map // 地址映射实例
	GoroutineMap sync.Map // 协程映射地址
)

type ModLoader struct {
	ModMap      sync.Map       // 模块列表
	TaskChan    chan *ChanTask // 任务通道
	TaskCount   int32          // 剩余任务数量
	Stack       atomic.Value   // 栈信息
	LoadFlag    int32          // 加载标记,当>0表示加载完成
	goroutineID int64          // 协程ID
	addr        string         // 协程地址
	saveDirty   int32          // 存档脏标记
	taskFlag    int32          // 执行标记,当>0表明当前有任务阻塞，直接调用
}

// 初始化
func (that *ModLoader) Init() {
	if that.TaskChan == nil {
		that.TaskChan = make(chan *ChanTask, CHAN_TASK_SIZE)
	}
}

// 读档
func (that *ModLoader) Load() {
	that.ModMap.Range(func(_, value interface{}) bool {
		mod, ok := value.(IModule)
		if ok {
			mod.Load()
		}
		return true
	})
}

// 读档后处理
func (that *ModLoader) LaterLoad() {
	that.ModMap.Range(func(_, value interface{}) bool {
		mod, ok := value.(IModule)
		if ok {
			mod.LaterLoad()
		}
		return true
	})
	atomic.StoreInt32(&that.LoadFlag, 1)
}

// 启动模块
func (that *ModLoader) OnStart() {
	that.ModMap.Range(func(_, value interface{}) bool {
		mod, ok := value.(IModule)
		if ok {
			mod.OnStart()
		}
		return true
	})
}

// 存档
func (that *ModLoader) Save() {
	that.ModMap.Range(func(_, value interface{}) bool {
		mod, ok := value.(IModule)
		if ok {
			mod.Save()
		}
		return true
	})
}

// 逻辑帧更新
func (that *ModLoader) Update(dt int64, online bool) {
	that.ModMap.Range(func(_, value interface{}) bool {
		mod, ok := value.(IModule)
		if ok {
			if online {
				mod.Update(dt)
			}
			if mod.HasSaveDirty() {
				mod.Save()
				mod.SetSaveDirty(false)
			}
		}
		return true
	})
	if atomic.LoadInt32(&that.saveDirty) == 1 {
		that.Save()
		atomic.StoreInt32(&that.saveDirty, 0)
	}
}

// 每日五点更新
func (that *ModLoader) OnRefresh(isLogin bool) {
	that.ModMap.Range(func(_, value interface{}) bool {
		mod, ok := value.(IModule)
		if ok {
			mod.OnRefresh()
		}
		return true
	})
}

// 同步消息
func (that *ModLoader) SendInfo() {
	that.ModMap.Range(func(_, value interface{}) bool {
		mod, ok := value.(IModule)
		if ok {
			mod.SendInfo()
		}
		return true
	})
}

// 获取任务通道
func (that *ModLoader) GetTaskChan() chan *ChanTask {
	if that.TaskChan == nil {
		that.TaskChan = make(chan *ChanTask, CHAN_TASK_SIZE)
	}
	return that.TaskChan
}

// 执行多线程任务
func (that *ModLoader) DoTask(task *ChanTask) {
	interval := that.getTaskInerval()
	task.SetTimeout(interval) // 重置倒计时
	atomic.AddInt32(&that.taskFlag, 1)
	if core.IsDebug {
		that.Stack.Store(task.Stack)
	}
	task.Result, task.Err = modInvokeInternal(that, task.Module, task.FunName, task.Args...)
	atomic.AddInt32(&that.taskFlag, -1)
	task.Done()
	if core.IsDebug {
		that.Stack.Store([]string{})
	}
}

// 获取任务间隔时间
func (that *ModLoader) getTaskInerval() time.Duration {
	return time.Millisecond * 1000
}

// 添加多线程任务
func (that *ModLoader) AddTask(task *ChanTask) ([]reflect.Value, error) {
	if len(that.TaskChan) >= CHAN_TASK_SIZE {
		return nil, fmt.Errorf("AddTask: Chan task overlap %d", CHAN_TASK_SIZE)
	}
	task.Add(1)
	that.TaskChan <- task // 通道切换任务消耗15ms左右
	startTime := core.ServerTime().UnixMilli()
	defer func() {
		atomic.AddInt32(&that.TaskCount, -1)
	}()

	// 无参数返回则不等待执行完成
	mod := that.GetModule(task.Module)
	numOut := mod.GetNumOut(task.FunName)
	if numOut < 1 {
		gopool.Go(func() {
			task.Future = core.Async(func() {
				task.Wait()
				task.Free()
			})
			_, err := task.Future.AwaitTimeout(time.Millisecond * 1000)
			if err != nil {
				task.Done()
				core.Logger.Errorf("调用非阻塞任务 %s_%s 超时,总耗时%d", task.Module, task.FunName, core.ServerTime().UnixMilli()-startTime)
			}
		})
		return nil, nil
	}

	// 阻塞等待
	task.Future = core.Async(func() {
		task.Wait()
		task.Free()
	})
	_, err := task.Future.AwaitTimeout(time.Millisecond * 1000)
	if err != nil {
		task.Done()
		core.Logger.Errorf("阻塞任务 %s_%s 超时,总耗时%d", task.Module, task.FunName, core.ServerTime().UnixMilli()-startTime)
		if core.IsDebug {
			for v := range that.TaskChan {
				core.Logger.Debugf("阻塞任务等待任务--%v", v.FunName)
			}
		}
	}
	return task.Result, task.Err
}

// 获取模块
func (that *ModLoader) GetModule(name string) IModule {
	value, ok := that.ModMap.Load(name)
	if ok {
		return value.(IModule)
	}
	return DefMod
}

// 设置模块
func (that *ModLoader) AddModule(mod IModule, host interface{}) IModule {
	if mod == nil || host == nil {
		core.Logger.Error("addModule invalid")
		return mod
	}
	// metaName := fmt.Sprint(reflect.TypeOf(mod))
	// metaArr := strings.Split(metaName, ".")
	mod.Init(host)
	that.ModMap.Store(mod.GetName(), mod)
	return mod
}

// 设置模块
func (that *ModLoader) RemoveModule(name string) {
	that.ModMap.Delete(name)
}

// 获取模块
func (that *ModLoader) ForEach(cb func(manager IModule)) {
	that.ModMap.Range(func(_, value interface{}) bool {
		manager, ok := value.(IModule)
		if ok && cb != nil {
			cb(manager)
		}
		return true
	})
}

// 设置模块
func (that *ModLoader) SetSaveDirty() {
	atomic.StoreInt32(&that.saveDirty, 1)
}

// 设置协程id
func (that *ModLoader) SetGoroutineID(name string) {
	if that.addr == "" {
		that.addr = "0x" + strconv.FormatInt(int64(reflect.ValueOf(that).Pointer()), 16)
		AddrMap.Store(that.addr, that)
	}
	goroutineID := GetGoroutineID()
	if that.goroutineID == goroutineID {
		return
	}
	if that.goroutineID > 0 {
		GoroutineMap.Delete(that.goroutineID)
	}
	that.goroutineID = goroutineID
	GoroutineMap.Store(that.goroutineID, that.addr)
	core.Logger.Debugf("【%s】设置调度协程addr=%s,id=%d", name, that.addr, that.goroutineID)
}

// 设置模块
func (that *ModLoader) GetGoroutineID() int64 {
	return that.goroutineID
}

// 销毁
func (that *ModLoader) OnDestory() {
	that.ModMap.Range(func(_, value interface{}) bool {
		mod, ok := value.(IModule)
		if ok {
			mod.OnDestory()
		}
		return true
	})
}

// 判定循环死锁条件如下：
// 栈列表不算引用计数，栈列表最后一位必须为运行中loader，且有计数
func (that *ModLoader) checkDeadLock(top string) ([]string, bool) {
	val, ok := AddrMap.Load(top)
	if !ok {
		core.Logger.Errorf("堆栈地址%s映射加载器丢失", top)
		return nil, false
	}

	// 两层任务以上加入堆栈队列
	taskStack, ok := val.(*ModLoader).Stack.Load().([]string)
	if ok && len(taskStack) > 0 {
		// 如果递归出现则死锁
		for _, v := range taskStack {
			if v == that.addr {
				return nil, true
			}
		}
		taskStack = append(taskStack, that.addr)
	} else {
		taskStack = []string{top, that.addr}
	}

	// 当前加载器无运行中任务，则直接加入队列
	if atomic.LoadInt32(&that.taskFlag) < 1 {
		return taskStack, false
	}
	for i := 0; i < len(taskStack)-1; i++ {
		addr := taskStack[i]
		val, ok = AddrMap.Load(addr)
		if !ok {
			continue
		}

		// 加上自身计数,引用次数超过1次代表循环死锁
		loader := val.(*ModLoader)
		if atomic.LoadInt32(&loader.TaskCount) > 1 {
			core.Logger.Infof("检测到任务次数=%d, GoroutineID=%d, addr=%s循环死锁", loader.TaskCount, that.goroutineID, that.addr)
			return nil, true
		}
	}
	return taskStack, false
}

// 检查是否可直接调用
func (that *ModLoader) CheckCall(modName string, funcName string) ([]string, bool) {
	curGoroutineID := GetGoroutineID()
	if curGoroutineID == that.goroutineID || core.IsTesting {
		return nil, true
	}
	if that.goroutineID < 1 {
		core.Logger.Error("未初始化协程")
		return nil, true
	}
	if !core.IsDebug {
		return nil, false // 非调试模式直接加入
	}

	// 非阻塞则直接加入
	mod := that.GetModule(modName)
	numOut := mod.GetNumOut(funcName)
	if numOut < 1 {
		return nil, false
	}

	// 栈顶则直接加入
	var top string
	val, ok := GoroutineMap.Load(curGoroutineID)
	if ok {
		top = val.(string)
	} else {
		stack := string(debug.Stack())
		offset := strings.Index(stack, "(*ModLoader).DoTask")
		if offset < 0 {
			return nil, false
		}
		top = stack[offset+20 : offset+32]
	}

	// 循环死锁则直接执行任务，打断死锁条件
	return that.checkDeadLock(top)
}

// 执行切换协程
func (that *ModLoader) doSwitchCoroutine(cb func()) {
	if cb != nil {
		cb()
	}
}

// 切换协程
func (that *ModLoader) SwitchCoroutine(cb func()) {
	ModInvoke(that, "__core", "SwitchCoroutine", cb)
}

// ----------------------------------------------------
// 内置调用模块函数（不可与ModInvoke合并，避免递归问题）
func modInvokeInternal(that *ModLoader, modName string, funcName string, args ...interface{}) ([]reflect.Value, error) {
	if GetGoroutineID() != that.GetGoroutineID() {
		if !core.IsTesting {
			core.Logger.Warnf("当期任务协程未切换,%s_%s", modName, funcName)
		}
	}
	if modName == "__core" && len(args) > 0 {
		cb, _ := args[0].(func())
		that.doSwitchCoroutine(cb)
		return nil, nil
	}
	mod := that.GetModule(modName)
	return mod.Invoker(funcName, args...)
}

// 调用模块函数（同步模型，强制等待切回当前协程）
func ModInvoke(that *ModLoader, modName string, funcName string, args ...interface{}) ([]reflect.Value, error) {
	// 当相同线程或者线程被阻塞的时候，直接调用将可避开多线程冲突
	// 线程被阻塞将阻塞update逻辑和消息进入
	stack, check := that.CheckCall(modName, funcName)
	if check {
		return modInvokeInternal(that, modName, funcName, args...)
	} else {
		task := NewChanTask(modName, funcName, args)
		task.Stack = stack // 携带堆栈
		return that.AddTask(task)
	}
}

// 调用模块函数（同步模型，可多线程安全访问的数据）
func ModInvokeSafe(that *ModLoader, modName string, funcName string, args ...interface{}) ([]reflect.Value, error) {
	mod := that.GetModule(modName)
	return mod.Invoker(funcName, args...)
}
