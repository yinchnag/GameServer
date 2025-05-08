package interfaces

type IRceiver interface {
	Init(any)
	GetName() string
	HandlerEvent()
	Dispay()
}
