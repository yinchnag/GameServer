package test

import (
	"context"
	"testing"
	"time"

	"wgame_server/libray/database"
)

func GetDbiRedisV8() *database.DbiRedisV8 {
	redis := &database.DbiRedisV8{}
	redis.Init("localhost:6379", 0, "123456")
	return redis
}

func TestRedisV8Set(t *testing.T) {
	// 测试Redis连接池
	redis := GetDbiRedisV8()
	result, err := redis.Set("testt", "test")
	if err != nil {
		t.Error(err)
	} else {
		t.Log(result)
	}
}

func TestRedisV8Append(t *testing.T) {
	redis := GetDbiRedisV8()
	ctx := context.Background()
	result, err := redis.Append(ctx, "testt", "123")
	if err != nil {
		t.Error(err)
	} else {
		t.Log(result)
	}
}

func TestRedisV8SetEx(t *testing.T) {
	redis := GetDbiRedisV8()
	ctx := context.Background()
	result, err := redis.SetEx(ctx, "testt", "123", 10*time.Second)
	if err != nil {
		t.Error(err)
	} else {
		t.Log(result)
	}
}
