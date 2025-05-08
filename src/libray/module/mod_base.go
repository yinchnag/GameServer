package module

import (
	"fmt"
	"reflect"
	"runtime/debug"
	"strings"

	"wgame_server/libray/core"
)

var DefMod *ModObj = &ModObj{}

type IModule interface {
	Init(interface{}) IModule
	Load()
	LaterLoad()
	OnStart()
	Save()
	Update(int64)
	HasSaveDirty() bool
	SetSaveDirty(bool)
	OnRefresh()
	SendInfo()
	OnDestory()
	Invoker(funcName string, args ...interface{}) ([]reflect.Value, error)
	GetNumOut(funcName string) int
	GetName() string

	GetRobot() IModule
}

type ModObj struct {
	modName  string                   // 模块名称
	invokers map[string]reflect.Value // 反射列表
}

func (that *ModObj) SetInvokerAll(inst interface{}) {
	metaName := fmt.Sprint(reflect.TypeOf(inst))
	metaArr := strings.Split(metaName, ".")
	that.modName = metaArr[1]
	// 基础函数不导出
	filterArr := []string{}
	baseVal := reflect.ValueOf(that)
	baseType := baseVal.Type()
	baseNum := baseType.NumMethod()
	for i := 0; i < baseNum; i++ {
		name := baseType.Method(i).Name
		if name == "Save" || name == "SendInfo" {
			continue
		}
		method := baseVal.Method(i)
		if method.CanInterface() {
			filterArr = append(filterArr, name)
		}
	}

	// 导出
	instVal := reflect.ValueOf(inst)
	instType := instVal.Type()
	methodNum := instType.NumMethod()
	for i := 0; i < methodNum; i++ {
		name := instType.Method(i).Name
		idx := core.FindSlice(filterArr, func(value string, _ int) bool {
			return name == value
		})
		if idx != -1 {
			continue
		}
		method := instVal.Method(i)
		if method.CanInterface() {
			that.SetInvoker(name, method)
		}
	}
}

// 注册函数
func (that *ModObj) SetInvoker(funcName string, cb interface{}) {
	kind := reflect.TypeOf(cb).Kind()
	if kind == reflect.Func {
		that.invokers[funcName] = reflect.ValueOf(cb)
	} else {
		that.invokers[funcName] = cb.(reflect.Value)
	}
}

// 初始化基础模块
func (that *ModObj) Init(mod interface{}) IModule {
	that.invokers = make(map[string]reflect.Value)
	return that
}

func (that *ModObj) Load() {
}

func (that *ModObj) LaterLoad() {
}

func (that *ModObj) OnStart()           {}
func (that *ModObj) Save()              {}
func (that *ModObj) Update(int64)       {}
func (that *ModObj) HasSaveDirty() bool { return false }
func (that *ModObj) SetSaveDirty(bool)  {}
func (that *ModObj) OnRefresh()         {}
func (that *ModObj) SendInfo()          {}

func (that *ModObj) OnDestory() {}

// 获取返回值数量
func (that *ModObj) GetNumOut(funcName string) int {
	handler, ok := that.invokers[funcName]
	if !ok {
		return 0
	}
	funcType := handler.Type()
	argOutNum := funcType.NumOut()
	return argOutNum
}

// 获得模块名称
func (that *ModObj) GetName() string {
	return that.modName
}

// 调用注册函数
func (that *ModObj) Invoker(funcName string, args ...interface{}) (_ []reflect.Value, reterr error) {
	handler, ok := that.invokers[funcName]
	if !ok {
		errMsg := fmt.Errorf("can't find %s_%s invokers len:%d", that.GetName(), funcName, len(that.invokers))
		core.Logger.Error(errMsg.Error())
		return nil, errMsg
	}
	defer func() {
		if err := recover(); err != nil {
			reterr = fmt.Errorf("failed to invoker %s_%s %v", that.GetName(), funcName, err)
			core.Logger.Error(reterr)
			if core.IsDebug {
				stack := string(debug.Stack())
				core.Logger.Error(stack)
			}
		}
	}()
	funcType := handler.Type()
	options := make([]reflect.Value, len(args))
	for i, v := range args {
		if v == nil {
			options[i] = reflect.New(funcType.In(i)).Elem()
		} else {
			options[i] = reflect.ValueOf(v)
		}
	}
	retArr := handler.Call(options)
	if len(retArr) > 0 {
		last := retArr[len(retArr)-1].Interface()
		opterr, ok := last.(error)
		if ok {
			return nil, opterr // 捕捉函数返回错误
		}
	}
	return retArr, nil
}

func (that *ModObj) GetRobot() IModule {
	return nil
}
