package server

import (
	"sync/atomic"

	"wgame_server/libray/actor"
	"wgame_server/libray/core"
	"wgame_server/libray/define/PB"
	"wgame_server/libray/entity"
	"wgame_server/libray/manager"
	"wgame_server/libray/module"
	"wgame_server/libray/network"
	"wgame_server/module/activity"
	"wgame_server/module/player"

	"google.golang.org/protobuf/reflect/protoreflect"
)

type ZoneServer struct {
	wsServer    *network.WsServer // ws服务器
	alive       int32             // 是否激活
	modMgr      *module.ModMgr
	actorSystem *actor.ActorSystem
}

func (that *ZoneServer) Start() {
	that.startActor()
	network.WGServer = that
	if atomic.LoadInt32(&that.alive) == 1 {
		return
	}
	atomic.StoreInt32(&that.alive, 1)

	manager.GetSignalManager().AddListener(manager.SIGNAL_START_SERVER, func() {
		that.RegisterModule(player.NewPlayerMgr())
	})

	manager.GetSignalManager().AddListener(manager.SIGNAL_PLAYER_NEW, func(user *entity.PlayerObj) {
		user.AddModule(&player.PlayerMod{})
	})

	that.modMgr = module.GetModMgr()
	manager.GetSignalManager().Notify(manager.SIGNAL_START_SERVER)
	that.wsServer = &network.WsServer{}
	that.wsServer.Init()
	that.wsServer.Start(":8000", 1000)
}

func (that *ZoneServer) startActor() {
	that.actorSystem = actor.NewActorSystem()
	that.actorSystem.AllocActor(func() actor.IRceiver { return &activity.ActivityActor{} })
}

func (that *ZoneServer) LoadConfig() {
}

func (that *ZoneServer) LoadDb() {
}

func (that *ZoneServer) ConnectRpc() {
}

func (that *ZoneServer) Login(session core.ISocketSession, meta protoreflect.ProtoMessage) {
	player.GetPlayerMgr().Login(session, meta.(*PB.C2S_Player_Login))
}

func (that *ZoneServer) RegisterModule(manager module.IModule) {
	if manager != nil {
		that.modMgr.AddModule(manager, that.modMgr)
		manager.LaterLoad()
		core.Logger.Infof("WsServer 加载模块 %s", manager.GetName())
	}
}
