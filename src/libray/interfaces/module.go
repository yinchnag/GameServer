package interfaces

// 函数执行顺便， Start->Load->LaterLoad
type IModule interface {
	Start()
	Load()
	LaterLoad()
	Save()
	Update()
}
