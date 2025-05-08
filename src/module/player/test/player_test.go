package test

import (
	"testing"

	"wgame_server/libray/entity"
	"wgame_server/libray/manager"
	"wgame_server/module/player"
)

func TestPlayerSigna(t *testing.T) {
	manager.GetSignalManager().AddListener(manager.SIGNAL_PLAYER_NEW, func(user *entity.PlayerObj) {
		user.AddModule(&player.PlayerMod{})
	})
	manager.GetSignalManager().Notify(manager.SIGNAL_PLAYER_NEW, &entity.PlayerObj{})
}
