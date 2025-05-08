package player

import (
	"wgame_server/libray/core"
	"wgame_server/libray/define/PB"
	"wgame_server/libray/module"
	"wgame_server/module/player/robot"
)

type PlayerMod struct {
	module.ModObj
}

func (that *PlayerMod) Init(mod interface{}) module.IModule {
	that.ModObj.Init(mod)
	that.SetInvokerAll(that)
	return that
}

func (that *PlayerMod) GetRobot() module.IModule {
	robot := &robot.PlayerRobot{}
	robot.ModObj.Init(robot)
	robot.SetInvokerAll(robot)
	return robot
}

// 登录入口
//
//	export PlayerMod_Login
func (that *PlayerMod) Login(msg *PB.C2S_Player_Login) {
	core.Logger.Debugf("PlayerMod_Login %v", msg)
}
