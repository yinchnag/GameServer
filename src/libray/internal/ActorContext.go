package internal

import (
	"fmt"
	"server/src/library/interfaces"
)

type ActorContext struct {
	receivers map[string]interfaces.IRceiver
}

func (that *ActorContext) Plush(receiver interfaces.IRceiver) {
	receiver.Init(receiver)
	name := receiver.GetName()
	if _, ok := that.receivers[name]; !ok {
		that.receivers[name] = receiver
		fmt.Printf("Plush Receiver: %s \n", name)
	}
}
