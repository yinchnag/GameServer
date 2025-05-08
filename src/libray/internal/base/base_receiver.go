package base

import (
	"fmt"
	"reflect"

	"wgame_server/libray/interfaces"
)

// 反射函数信息 -- 运行时数据，不能存盘
// 该结构体用于存储函数的反射信息
type FuncInfo struct {
	funcName   string         // 函数名称
	method     reflect.Value  // 函数反射对象
	inputArgs  []reflect.Kind // 输入参数类型
	outputArgs []reflect.Kind // 输出参数类型
}

var (
	_ interfaces.IRceiver = &BaseReceiver{}
	_ interfaces.IModule  = &BaseReceiver{}
)

// 模块对象基础类
// 继承(组合)该类的对象需主动调用Init函数
type BaseReceiver struct {
	invokers map[string]*FuncInfo // 所有反射获得的函数
	name     string               // 对象名称
}

// 初始化对象
// 在任何情况下创建出对象后该函数一定是第一个调用，否则会出现painc
func (that *BaseReceiver) Init(val any) {
	that.invokers = make(map[string]*FuncInfo)
	that.SetInvokerAll(val)
}

// 设置对象中所有函数的反射
func (that *BaseReceiver) SetInvokerAll(val any) {
	instVal := reflect.ValueOf(val)
	instType := instVal.Type()
	that.name = instType.Name()
	methodNum := instType.NumMethod()
	for i := 0; i < methodNum; i++ {
		name := instType.Method(i).Name
		method := instVal.Method(i)
		if method.CanInterface() {
			kind := method.Kind()
			_ = kind
			that.SetInvoker(name, method)
		}
	}
}

// 判断是否为函数，如果是，则记录
func (that *BaseReceiver) SetInvoker(funcName string, method reflect.Value) {
	if method.Kind() != reflect.Func {
		return
	}
	info := &FuncInfo{
		funcName: funcName,
		method:   method,
	}
	typeOf := method.Type()
	for i := 0; i < typeOf.NumIn(); i++ { // 获取输入参数类型
		info.inputArgs = append(info.inputArgs, typeOf.In(i).Kind())
	}
	for i := 0; i < typeOf.NumOut(); i++ { // 获取输出参数类型
		info.outputArgs = append(info.outputArgs, typeOf.Out(i).Kind())
	}
	that.invokers[funcName] = info
	fmt.Printf("bind method: %s\n", funcName)
}

func (that *BaseReceiver) Invoker(uid int64, funcName string, args ...any) (_ []reflect.Value, err error) {
	if invoker, ok := that.invokers[funcName]; ok {
		defer func() {
			if err := recover(); err != nil {
				fmt.Printf("调用函数 %s 失败, err: %v\n", funcName, err)
				err = fmt.Errorf("调用函数 %s 失败, err: %v", funcName, err)
			}
		}()
		if len(args) != len(invoker.inputArgs) {
			return nil, fmt.Errorf("参数数量不匹配")
		}
		inArgs := make([]reflect.Value, 1)
		for i, arg := range args {
			inArgs[i] = reflect.ValueOf(arg)
		}
		outArgs := invoker.method.Call(inArgs)
		return outArgs, nil
	}
	return nil, fmt.Errorf("没有找到函数 %s", funcName)
}

func (that *BaseReceiver) GetName() string {
	return that.name
}

func (that *BaseReceiver) HandlerEvent() {
}

func (that *BaseReceiver) Destory() {
	that.invokers = nil
}

// 模块被启动时第一个调用
func (that *BaseReceiver) Start() {
}

// 模块帧函数
func (that *BaseReceiver) Update() {
}
