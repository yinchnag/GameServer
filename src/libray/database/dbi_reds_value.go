package database

import "github.com/bytedance/sonic"

type IRedisValue[T any] interface {
	Interface(string) T
}

type DbiRedisValue[T any] struct {
	prefix []string
}

func (that *DbiRedisValue[T]) AddPrefix(prefix ...string) {
	that.prefix = prefix
}

func (that *DbiRedisValue[T]) New() T {
	return *new(T)
}

func (that *DbiRedisValue[T]) Interface(val string) T {
	t := new(T)
	// s, _ := sonic.Marshal(t)
	sonic.Unmarshal([]byte(val), &t)
	return *t
}

func AddCacheValue[T any](value DbiRedisValue[T]) {
}
