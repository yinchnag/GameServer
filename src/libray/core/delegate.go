package core

import (
	"fmt"
	"reflect"
	"sync"
	"sync/atomic"
)

type Delegate struct {
	sync.Mutex                        //互斥锁
	invokers   sync.Map               //事件通知列表
	con        map[int][]reflect.Kind //事件参数对照表
	lockCount  int32                  //锁次数，检查循环锁
}

func (that *Delegate) Init() {
	that.con = make(map[int][]reflect.Kind)
}

// 加锁通知列表
func (that *Delegate) lockInvokers() {
	if atomic.AddInt32(&that.lockCount, 1) == 1 {
		that.Lock()
	}
}

// 解锁通知列表
func (that *Delegate) unLockInvokers() {
	if atomic.AddInt32(&that.lockCount, -1) == 0 {
		that.Unlock()
	}
}

func (that *Delegate) AddListener(event int, cb interface{}) error {
	if reflect.TypeOf(cb).Kind() != reflect.Func {
		return fmt.Errorf("%s is not of type reflect.Func", reflect.TypeOf(cb).Kind())
	}
	that.lockInvokers()
	defer that.unLockInvokers()
	funcArr := that.getInvokers(event)
	handler, parameter := that.getHandlerAndParameter(cb)
	if !that.checkContrasting(event, parameter) {
		return fmt.Errorf("event %d parameter error,of %s", event, reflect.TypeOf(cb).Kind())
	}
	index := FindSlice(funcArr, func(value reflect.Value, key int) bool {
		return value == handler
	})
	if index == -1 {
		funcArr = append(funcArr, handler)
		that.invokers.Store(event, funcArr)
	}
	return nil
}

func (that *Delegate) RemoveListener(event int, cb interface{}) error {
	if reflect.TypeOf(cb).Kind() != reflect.Func {
		return fmt.Errorf("%s is not of type reflect.Func", reflect.TypeOf(cb).Kind())
	}
	funcArr := that.getInvokers(event)
	funcArr = SliceRemoveByVal(funcArr, reflect.ValueOf(cb))
	that.invokers.Store(event, funcArr)
	return nil
}

// 事件通知
func (that *Delegate) Notify(event int, args ...interface{}) {
	actual, ok := that.invokers.Load(event)
	if !ok {
		return
	}
	funcArr, _ := actual.([]reflect.Value)
	if len(funcArr) < 1 {
		return
	}

	defer func() {
		if err := recover(); err != nil {
			Logger.Errorf("notify event %d raise %v", event, err)
		}
	}()
	funcType := funcArr[0].Type()
	options := make([]reflect.Value, len(args))
	for i, v := range args {
		if v == nil {
			options[i] = reflect.New(funcType.In(i)).Elem()
		} else {
			options[i] = reflect.ValueOf(v)
		}
	}
	for _, item := range funcArr {
		item.Call(options)
	}
}
func (that *Delegate) checkContrasting(event int, parameter []reflect.Kind) bool {
	con, ok := that.con[event]
	if !ok {
		that.con[event] = append(that.con[event], parameter...)
		return true
	}
	if len(con) != len(parameter) {
		Logger.Errorf("event %d parameter error,of len :%d", event, len(parameter))
		return false
	}
	for i := 0; i < len(con); i++ {
		if con[i] != parameter[i] {
			Logger.Errorf("event %d parameter error,of %s", event, parameter[i])
			return false
		}
	}
	return true
}

func (that *Delegate) getHandlerAndParameter(cb interface{}) (handler reflect.Value, parameter []reflect.Kind) {
	handler = reflect.ValueOf(cb)
	refType := reflect.TypeOf(cb)
	for i := 0; i < refType.NumIn(); i++ {
		parameter = append(parameter, refType.In(i).Kind())
	}
	return handler, parameter
}

func (that *Delegate) getInvokers(event int) []reflect.Value {
	actual, _ := that.invokers.LoadOrStore(event, make([]reflect.Value, 0))
	funcArr, _ := actual.([]reflect.Value)
	return funcArr
}
