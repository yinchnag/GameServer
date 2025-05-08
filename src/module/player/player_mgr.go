package player

import (
	"sync"

	"wgame_server/libray/core"
	"wgame_server/libray/define/PB"
	"wgame_server/libray/entity"
	"wgame_server/libray/manager"
	"wgame_server/libray/module"
	"wgame_server/libray/network"
)

var (
	playerMgr     *PlayerMgr
	playerMgrOnce sync.Once
)

type PlayerMgr struct {
	PlayerMap sync.Map                    // 玩家MAP
	user      map[int64]*entity.PlayerObj // 在线玩家列表

	module.ModObj
}

// 单例(多线程安全)
func GetPlayerMgr() *PlayerMgr {
	playerMgrOnce.Do(func() {
		playerMgr = &PlayerMgr{}
	})
	return playerMgr
}

func NewPlayerMgr() *PlayerMgr {
	playerMgrOnce.Do(func() {
		playerMgr = &PlayerMgr{}
	})
	return playerMgr
}

func (that *PlayerMgr) Init(mod interface{}) module.IModule {
	that.user = make(map[int64]*entity.PlayerObj)
	that.ModObj.Init(mod)
	that.SetInvokerAll(that)
	return that
}

func (that *PlayerMgr) LaterLoad() {
	manager.GetSignalManager().AddListener(manager.SIGNAL_CONNECT_SERVER, that.Login)
}

func (that *PlayerMgr) newPlayer(uid int64, ws core.ISocketRecv) *entity.PlayerObj {
	user := &entity.PlayerObj{}
	user.Entity = entity.NewEntity(ws)
	user.EntityId = uid
	that.user[uid] = user
	manager.GetSignalManager().Notify(manager.SIGNAL_PLAYER_NEW, user)
	return user
}

func (that *PlayerMgr) Login(ws core.ISocketSession, meta *PB.C2S_Player_Login) {
	user := that.newPlayer(1, ws.(*network.WsConnectSession).GetSocketRecv())
	user.Load()
	// user.LaterLoad()
}
