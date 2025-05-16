package activity

import (
	"strconv"

	"wgame_server/libray/actor"
)

type ActivityActor struct {
	data map[string]int
	actor.ActorReceiver
}

func (that *ActivityActor) Init() {
	that.data = map[string]int{}
	for i := 0; i < 1000; i++ {
		that.data[strconv.Itoa(i)] = i
	}
}

// 获得值
//
//	export ActivityActorGetInt
func (that *ActivityActor) GetInt(val string) int {
	ret, ok := that.data[val]
	if ok {
		return -1
	}
	return ret
}
