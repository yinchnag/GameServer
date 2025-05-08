package base

import (
	"fmt"
	"reflect"

	"wgame_server/libray/interfaces"
)

var ReceiverFunc []string = []string{}

var (
	_ interfaces.IRceiver = &BaseReceiver{}
	_ interfaces.IModule  = &BaseReceiver{}
)

// 模块对象基础类
// 继承(组合)该类的对象需主动调用Init函数
type BaseReceiver struct {
	invald map[string]reflect.Value // 所有反射获得的函数
	name   string                   // 对象名称
}

// 初始化对象
// 在任何情况下创建出对象后该函数一定是第一个调用，否则会出现painc
func (that *BaseReceiver) Init(val any) {
	that.invald = make(map[string]reflect.Value)
	that.reflectFunc(val)
}

func (that *BaseReceiver) reflectFunc(val any) {
	// valValue := reflect.ValueOf(val)
	valType := reflect.TypeOf(val)
	that.name = valType.Name()
	for i := 0; i < valType.NumMethod(); i++ {
		method := valType.Method(i)
		fmt.Printf("bind method: %s\n", method.Name)
		that.invald[method.Name] = method.Func
	}
}

func (that *BaseReceiver) GetName() string {
	return that.name
}

func (that *BaseReceiver) HandlerEvent() {
}

func (that *BaseReceiver) Destory() {
	that.invald = nil
}

// 模块被启动时第一个调用
func (that *BaseReceiver) Start() {
}

// 模块帧函数
func (that *BaseReceiver) Update() {
}
