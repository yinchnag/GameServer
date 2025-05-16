package test

import (
	"wgame_server/libray/actor"
	"wgame_server/module/activity"
)

var actorSystem *actor.ActorSystem

func createServer() {
	actorSystem = actor.NewActorSystem()
	actorSystem.AllocActor(func() actor.IRceiver { return &activity.ActivityActor{} })
}
