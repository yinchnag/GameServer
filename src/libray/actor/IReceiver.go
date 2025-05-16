package actor

import "reflect"

type IRceiver interface {
	init(*ActorContext, any)
	GetName() string
	GetNumOut(string) int
	Receive(any)
	Invoker(int64, string, ...any) ([]reflect.Value, error)
	HandlerEvent()
	Destory()
}
