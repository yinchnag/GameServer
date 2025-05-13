package database

import "wgame_server/libray/core"

type DbiServer struct {
	conf        *core.JS_DatabaseConf
	redisClient DbiRedis
}

func (that *DbiServer) Init(conf *core.JS_DatabaseConf) {
}
