package server

import (
	"wgame_server/libray/core"
	"wgame_server/libray/define/PB"
	"wgame_server/libray/network"

	"google.golang.org/protobuf/proto"
)

type C2S_Player_Login struct {
	msg PB.C2S_Player_Login
	network.LogicProto
}

func (that *C2S_Player_Login) new() network.ILogicProto {
	that.LogicProto.Init(that)
	return that
}

func (that *C2S_Player_Login) GetProtoID() int {
	return int(PB.C2S_Player_c2s_player_login)
}

func (that *C2S_Player_Login) GetMsg() proto.Message {
	return core.HF_ReflectCopy(&that.msg).(proto.Message)
}

func init() {
	network.C2S_MSG_MAP = append(network.C2S_MSG_MAP, (&C2S_Player_Login{}).new())
}
