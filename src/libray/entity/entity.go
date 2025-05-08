package entity

import (
	"wgame_server/libray/core"
	"wgame_server/libray/network"

	"google.golang.org/protobuf/proto"
)

type Entity struct {
	EntityId int64
	core.ISocketRecv
}

func NewEntity(conn core.ISocketRecv) *Entity {
	return &Entity{
		ISocketRecv: conn,
	}
}

// 发送消息
func (that *Entity) SendMessage(msgId int, msg proto.Message) {
	msgBodies, err := proto.Marshal(msg)
	if err != nil {
		core.Logger.Errorf("Failed to marshal proto msgid=%d", msgId)
		return
	}
	that.SendBytes(network.HF_EncodeMsgPB(uint16(msgId), msgBodies))
	core.Logger.Debugf("SendBytes: msgid=%d", msgId)
}
