package module

import (
	"reflect"
	"runtime/debug"
	"sync"
	"time"

	"wgame_server/libray/core"
)

var (
	modMgr     *ModMgr   // 单例
	modMgrOnce sync.Once // 单次加锁
)

// 单例(多线程安全)
func GetModMgr() *ModMgr {
	modMgrOnce.Do(func() {
		modMgr = &ModMgr{}
	})
	return modMgr
}

// 模块管理器
type ModMgr struct {
	ModLoader
}

// 初始化(务必调用)
func (that *ModMgr) Init() {
	that.ModLoader.Init()
}

// 初始化管理器
func (that *ModMgr) GetName() string {
	return reflect.TypeOf(that).Elem().Name()
}

// 默认运行更新
func (that *ModMgr) RunUpdateLoop() {
	defer func() {
		err := recover()
		if err != nil {
			core.Logger.Errorln(err, string(debug.Stack()))
		}
	}()

	// 逻辑循环
	that.SetGoroutineID("全局管理器")
	that.Load()
	that.LaterLoad()
	that.OnStart()
	logicStart := time.Now()
	ticker := time.NewTicker(time.Second)
	for {
		select {
		case task := <-that.TaskChan:
			that.DoTask(task)
		case dt := <-ticker.C:
			that.Update(dt.Sub(logicStart).Milliseconds(), true)
			logicStart = dt
		}
	}
}
