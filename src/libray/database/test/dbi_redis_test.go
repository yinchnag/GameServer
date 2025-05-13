package test

import (
	"strconv"
	"testing"

	"wgame_server/libray/database"
)

func GetDbiRedis() *database.DbiRedis {
	redis := &database.DbiRedis{}
	redis.Init("127.0.0.1:6379", 0, "Yc942628", "")
	return redis
}

func TestRedisConn(t *testing.T) {
	// 测试Redis连接池
	redis := GetDbiRedis()
	result, err := redis.Ping()
	if err != nil {
		t.Error(err)
	} else {
		t.Log(result)
	}
}

func TestRedisEcho(t *testing.T) {
	// 测试Redis连接池
	redis := GetDbiRedis()
	result, err := redis.Echo()
	if err != nil {
		t.Error(err)
	} else {
		t.Log(result)
	}
}

func TestRedisInfo(t *testing.T) {
	// 测试Redis连接池
	redis := GetDbiRedis()
	result, err := redis.Info()
	if err != nil {
		t.Error(err)
	} else {
		t.Log(result)
	}
}

// 设置字符串
func TestRedisSet(t *testing.T) {
	redis := GetDbiRedis()
	result, err := redis.Set("test", "test")
	if err != nil {
		t.Error(err)
	} else {
		t.Log(result)
	}
}

// 修改字符串某处
func TestRedisAppend(t *testing.T) {
	redis := GetDbiRedis()
	result, err := redis.Append("test", "123")
	if err != nil {
		t.Error(err)
	} else {
		t.Log(result)
	}
}

// 截取字符串
func TestRedisGetRange(t *testing.T) {
	redis := GetDbiRedis()
	result, err := redis.GetRange("test", 4, 8)
	if err != nil {
		t.Error(err)
	} else {
		t.Log(result)
	}
}

// 为一个已存在的key值设置过期时间
func TestRedisExpire(t *testing.T) {
	redis := GetDbiRedis()
	result, err := redis.Expire("test", 5000)
	if err != nil {
		t.Error(err)
	} else {
		t.Log(result)
	}
}

// 键值不存在则设置，设置成功返回1，失败返回0
func TestRedisSetNx(t *testing.T) {
	redis := GetDbiRedis()
	result, err := redis.SetNx("test", "test")
	if err != nil {
		t.Error(err)
	} else {
		t.Log(result)
	}
}

// 创建一个有过期时间的键值对，成功返回1，失败返回0
func TestRedisSetEx(t *testing.T) {
	redis := GetDbiRedis()
	result, err := redis.SetEx("test", "test", 5)
	if err != nil {
		t.Error(err)
	} else {
		t.Log(result)
	}
}

// 获得一个键值
func TestRedisGet(t *testing.T) {
	redis := GetDbiRedis()
	result, ok, err := redis.Get("test")
	if err != nil || !ok {
		t.Error(err)
	} else {
		t.Log(result)
	}
}

// 判断一个键值是否存在
func TestRedisExists(t *testing.T) {
	redis := GetDbiRedis()
	result, err := redis.Exists("test")
	if err != nil {
		t.Error(err)
	} else {
		t.Log(result)
	}
}

// 删除一个键值
func TestRedisDel(t *testing.T) {
	redis := GetDbiRedis()
	result, err := redis.Del(false, "test")
	if err != nil {
		t.Error(err)
	} else {
		t.Log(result)
	}
}

// 根据正则表达式获得符合规则的键值对
func TestRedisKeys(t *testing.T) {
	redis := GetDbiRedis()
	result, err := redis.Keys("*")
	if err != nil {
		t.Error(err)
	} else {
		t.Log(result)
	}
}

// 增加一个键值对的值，成功返回增加后的值，失败返回0
func TestRedisIncrBy(t *testing.T) {
	redis := GetDbiRedis()
	result, err := redis.IncrBy("num", 2)
	if err != nil {
		t.Error(err)
	} else {
		t.Log(result)
	}
}

func BenchmarkRedisSet(b *testing.B) {
	redis := GetDbiRedis()

	key := make([]byte, 0, 32) // 预分配足够大的缓冲区
	for i := 0; i < b.N; i++ {
		key = strconv.AppendInt(key[:0], int64(i), 10) // 复用内存
		redis.Set("test"+string(key), "test")
	}
}
