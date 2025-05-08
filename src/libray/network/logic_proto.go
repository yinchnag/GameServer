package network

import (
	"fmt"
	"reflect"
	"strings"

	"wgame_server/libray/core"

	"google.golang.org/protobuf/proto"
)

const (
	MSG_ID_LOGIN = 1 // 登入协议
)

// 逻辑协议
type ILogicProto interface {
	Init(interface{}) ILogicProto
	HandleProtocol(session core.ISocketSession, msgData []byte, msgid int, meta proto.Message) bool
	GetProtoID() int
	GetMsg() proto.Message
}

type LogicProto struct {
	modName string // 协议所属模块
	meta    proto.Message
}

func (that *LogicProto) Init(proto interface{}) ILogicProto {
	metaName := fmt.Sprint(reflect.TypeOf(proto))
	metaArr := strings.Split(metaName, ".")

	that.modName = strings.Split(metaArr[1], "_")[1]
	return that
}

// 处理消息
func (that *LogicProto) HandleProtocol(session core.ISocketSession, msgData []byte, msgId int, meta proto.Message) bool {
	if msgData != nil {
		err := proto.Unmarshal(msgData, meta)
		if err != nil {
			core.Logger.Debugf("Failed to unmarshal proto msgid=%d", msgId)
			return false
		}
	}
	user := session.GetPlayer()
	if user != nil {
		user.ModInvoke(that.modName, meta)
	} else if msgId == MSG_ID_LOGIN { // 登入协议
		WGServer.Login(session, meta)
	} else { // 玩家还未登入
		session.SendError(1000)
	}
	return true
}

func (that *LogicProto) GetProtoID() int {
	return 0
}

func (that *LogicProto) GetMsg() proto.Message {
	return that.meta
}
