package extend

import "reflect"

var (
	cacheProtocols = make(map[string]*Protocol, 0) // 暂时使用map,压力测试后可能会换成切片
)

type Protocol struct {
	name string // 协议名称
}

func (that *Protocol) Init(val any) {
	that.reflectStructName(val)
}

func (that *Protocol) reflectStructName(val any) {
	if _, ok := cacheProtocols[that.name]; ok {
		valType := reflect.TypeOf(val)
		that.name = valType.Name()
		cacheProtocols[that.name] = that
	}
}
