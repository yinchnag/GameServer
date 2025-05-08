package test

import (
	"fmt"
	"testing"

	"wgame_server/libray/core"
	"wgame_server/libray/define/PB"
	"wgame_server/libray/module"
	"wgame_server/module/player"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
)

func TestProtoEnumToString(t *testing.T) {
	if PB.S2C_Player_s2c_player_error.String() == "PT_Test" {
		t.Log("ok")
	} else {
		t.Error(PB.S2C_Player_s2c_player_error.String())
		t.Fail()
	}
}

func TestProtoMessageToString(t *testing.T) {
	msg := &PB.S2C_Player_Login{}
	t.Log(msg.String())
}

func TestProtoHandle(t *testing.T) {
	pb := &PB.S2C_Player_Login{}
	pb.ServerId = 1
	pb.Account = []byte{1, 2, 3}
	pb.UserId = []byte{4, 5, 6}
	pb.Uid = 7
	pb.Uname = []byte{8, 9, 10}
	pb.Iconid = 11
	pb.Exp = 12
	pb.Level = 13
	pb.Regtime = 14
	pb.Sex = 15
	pb.Vip = 16
	pb.VipExp = 17
	pb.PowerUptime = 18
	pb.Fight = 19
	pb.GameServerTime = 20
	pb.Token = []byte{21, 22, 23}
	pb.RenameNum = 21
	pb.ZoneOffset = 22
	pb.OpenTime = 23
	pb.FavoriteDay1 = 24
	pb.FavoriteDay2 = 25
	pb.FavoriteCount = 26
	pb.LastOfflineTime = 24
	pb.FavoritePlayers = []int64{27, 28, 29}
	test := &TestMod{}
	test.Init(test)
	test.Invoker("Proto", core.HF_ReflectCopy(ProtoToBytes(pb)))
}

func ProtoToBytes(meta protoreflect.ProtoMessage) protoreflect.ProtoMessage {
	return meta
}

type TestMod struct {
	module.ModObj
}

func (that *TestMod) Init(mod interface{}) {
	that.ModObj.Init(mod)
	that.SetInvokerAll(that)
}

func (that *TestMod) Proto(meta *PB.S2C_Player_Login) {
	fmt.Println(meta.Account)
	fmt.Println(meta.Exp)
	fmt.Println(meta.FavoriteDay1)
}

func TestProtoInit(t *testing.T) {
	login := player.C2S_Player_Login{}
	meta := login.Init(&login)
	pb := PB.C2S_Player_Login{}
	pb.Account = []byte("123456")
	pb.Password = []byte("123456")
	bytedata, _ := proto.Marshal(&pb)
	meta.HandleProtocol(nil, bytedata, 1, login.GetMsg())
}
