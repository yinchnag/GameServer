package core

import (
	"reflect"
)

type (
	IEntity interface {
		GetPlayer() IPlayer
	}
	ISocketConnect interface {
		SendBytes([]byte) error // 发送消息
		SendError(int)
		IEntity
	}
	ISocketRecv interface {
		ISocketConnect
	}
	// 套接字会话
	ISocketSession interface {
		ISocketConnect
	}
	IPlayer interface {
		SetOnline(bool) // 设置在线状态
		Login(ISocketConnect)
		ModInvoke(string, ...interface{}) ([]reflect.Value, error)
	}
)
