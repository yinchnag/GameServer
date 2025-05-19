package bag

import "wgame_server/libray/actor"

type BagMod struct {
	data map[string]int
	actor.ActorReceiver
}

func (that *BagMod) GetItem(itemId int) {
}
