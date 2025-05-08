package test

import (
	"testing"

	"wgame_server/libray/core"
	"wgame_server/libray/define/PB"
)

func BenchmarkReflectCopy(b *testing.B) {
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

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ret := core.HF_ReflectCopy(pb)
		_ = ret
	}
}
