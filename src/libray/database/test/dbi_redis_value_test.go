package test

import (
	"testing"

	"wgame_server/libray/database"
)

type Item struct {
	ItemId int64 `json:"itemid"`
	Count  int64 `json:"count"`
	database.DbiRedisValue[Item]
}

func TestRedisValue(t *testing.T) {
	item := &Item{}
	cp := item.Interface(`{"itemid":1,"count":2}`)

	_ = cp
	t.Log(item)
}
