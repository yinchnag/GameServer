package network

import (
	"wgame_server/libray/core"
	"wgame_server/libray/module"

	"google.golang.org/protobuf/reflect/protoreflect"
)

type IWgameS interface {
	LoadConfig()                                                       // 加载配置
	LoadDb()                                                           // 加载数据库
	ConnectRpc()                                                       // 连接rpc
	Login(session core.ISocketSession, meta protoreflect.ProtoMessage) // 玩家登入入口
	RegisterModule(mod module.IModule)                                 // 注册模块
}

// 服务器接口对象
var WGServer IWgameS
