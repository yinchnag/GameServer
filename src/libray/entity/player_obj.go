package entity

import (
	"sync"

	"wgame_server/libray/core"
	"wgame_server/libray/module"
)

type PlayerObj struct {
	ModMap   sync.Map // 服务器 模块映射
	RobotMap sync.Map // 机器人 模块映射
	*Entity
}

func NewPlayer(entity *Entity) *PlayerObj {
	return &PlayerObj{Entity: entity}
}

func (that *PlayerObj) Init() {
}

func (that *PlayerObj) Load() {
}

func (that *PlayerObj) LaterLoad() {
}

func (that *PlayerObj) AddModule(mod module.IModule) module.IModule {
	if mod == nil {
		core.Logger.Error("player add module is nil")
		return mod
	}
	mod = mod.Init(that)
	that.ModMap.Store(mod.GetName(), mod)

	robot := mod.GetRobot()
	if !core.IsNil(robot) {
		that.RobotMap.Store(robot.GetName(), robot)
	}
	return mod
}
