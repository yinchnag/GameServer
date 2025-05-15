package interfaces

import "reflect"

type IRceiver interface {
	Init(any)
	GetName() string
	Receive(any)
	Invoker(int64, string, ...any) ([]reflect.Value, error)
	HandlerEvent()
	Destory()
}
