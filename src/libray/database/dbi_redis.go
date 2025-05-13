package database

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"

	"wgame_server/libray/core"

	"github.com/gomodule/redigo/redis"
)

// 常数定义
const (
	Max_Redis_Idle_Conn   = 30  // 最大空闲连接数，提前等待着，过了超时时间关闭
	Max_Redis_Active_Conn = 64  // 最大连接数，即最多的tcp连接数，一般建议往大的配置
	Max_Redis_Idle_Time   = 180 // 空闲连接超时时间，但应该设置比redis服务器超时时间短
)

// 操作定义
const (
	c_EXPIRE_OPTION         = "EX"
	c_NOT_EXISTS_OPTION     = "NX"
	c_MATCH_OPTION          = "MATCH"
	c_COUNT_OPTION          = "COUNT"
	c_SET_COMMAND           = "SET"
	c_DEL_COMMAND           = "DEL"
	c_GET_COMMAND           = "GET"
	c_KEYS_COMMAND          = "KEYS"
	c_PING_COMMAND          = "PING"
	c_ECHO_COMMAND          = "ECHO"
	c_INFO_COMMAND          = "INFO"
	c_HSET_COMMAND          = "HSET"
	c_HGET_COMMAND          = "HGET"
	c_HMSET_COMMAND         = "HMSET"
	c_HDEL_COMMAND          = "HDEL"
	c_HLEN_COMMAND          = "HLEN"
	c_HKEYS_COMMAND         = "HKEYS"
	c_SCAN_COMMAND          = "SCAN"
	c_HSCAN_COMMAND         = "HSCAN"
	c_GET_RANGE_COMMAND     = "GETRANGE"
	c_SET_RANGE_COMMAND     = "SETRANGE"
	c_EXPIRE_COMMAND        = "EXPIRE"
	c_EXISTS_COMMAND        = "EXISTS"
	c_HEXISTS_COMMAND       = "HEXISTS"
	c_HGETALL_COMMAND       = "HGETALL"
	c_INCRBY_COMMAND        = "INCRBY"
	c_INCRBYFLOAT_COMMAND   = "INCRBYFLOAT"
	c_HINCRBY_COMMAND       = "HINCRBY"
	c_HINCRBYFLOAT_COMMAND  = "HINCRBYFLOAT"
	c_TTL_COMMAND           = "TTL"
	c_APPEND_COMMAND        = "APPEND"
	c_ZADD_COMMAND          = "ZADD"
	c_ZSCORE_COMMAND        = "ZSCORE"
	c_ZRANGEBYSCORE_COMMAND = "ZRANGEBYSCORE"
	c_ZREM_COMMAND          = "ZREM"
	c_ZCARD_COMMAND         = "ZCARD"
	c_SADD_COMMAND          = "SADD"  // 向集合添加一个或多个成员
	c_SCARD_COMMAND         = "SCARD" // 获取集合的成员数
	c_SREM_COMMAND          = "SREM"  // 移除集合中一个或多个成员
	c_SPOP_COMMAND          = "SPOP"  // 移除并返回集合中的一个随机元素
	c_SSCAN_COMMAND         = "SSCAN" // 迭代集合中的元素
)

func newPool(ip string, dbidx int, auth string) *redis.Pool {
	return &redis.Pool{
		MaxIdle:     Max_Redis_Idle_Conn,
		MaxActive:   Max_Redis_Active_Conn, // max number of connections
		IdleTimeout: Max_Redis_Idle_Time * time.Second,
		Wait:        true, // 如果超过最大连接，是报错，还是等待
		Dial: func() (redis.Conn, error) {
			conn, err := redis.Dial("tcp", ip)
			if err != nil {
				core.Logger.Error(err.Error())
				return nil, err
			}
			if auth != "" {
				conn.Do("AUTH", auth)
			}
			conn.Do("SELECT", dbidx)
			return conn, err
		},
	}
}

// 分析Scan结果
func parseScanResults(results []interface{}) (int64, []string, error) {
	if len(results) != 2 {
		return 0, []string{}, nil
	}
	cursorIndex, err := strconv.ParseInt(string(results[0].([]byte)), 10, 64)
	if err != nil {
		return 0, nil, err
	}
	keyInterfaces := results[1].([]interface{})
	keys := make([]string, len(keyInterfaces))
	for index, keyInterface := range keyInterfaces {
		keys[index] = string(keyInterface.([]byte))
	}
	return cursorIndex, keys, nil
}

// 查询结果转字符串
func toString(reply interface{}, err error) (string, bool, error) {
	result, e := redis.String(reply, err)
	if e == redis.ErrNil {
		return result, false, nil
	}
	if e != nil {
		return result, false, e
	}
	return result, true, nil
}

// 分析转义
func toBool(reply interface{}, err error) (bool, error) {
	_, ok, e := toString(reply, err)
	return ok, e
}

type DbiRedisPipeline struct {
	conn redis.Conn
}

func (that *DbiRedisPipeline) Send(commandName string, args ...interface{}) error {
	return nil
}

func (that *DbiRedisPipeline) Exec(ctx context.Context) error {
	that.conn.Send("EXEC")
	return nil
}

type DbiRedis struct {
	pool   *redis.Pool // 连接池
	Host   string      // IP+端口，127.0.0.1:6379
	Index  int         // redis db index
	Auth   string      // redis密码
	Prefix string      // redis前缀
}

// 初始化
func (that *DbiRedis) Init(host string, dbidx int, auth string, prefix string) {
	that.Host = host
	that.Index = dbidx
	that.Auth = auth
	that.Prefix = prefix
	that.pool = newPool(host, dbidx, auth)
}

// 获得Redis连接
func (that *DbiRedis) GetRedisConn() redis.Conn {
	if that.pool != nil {
		return that.pool.Get()
	}
	return nil
}

// Ping消息
func (that *DbiRedis) Ping() (string, error) {
	conn := that.GetRedisConn()
	if conn == nil {
		return "", errors.New("no redis conn")
	}
	defer conn.Close()
	return redis.String(conn.Do(c_PING_COMMAND))
}

// Echo消息
func (that *DbiRedis) Echo() (string, error) {
	conn := that.GetRedisConn()
	if conn == nil {
		return "", errors.New("no redis conn")
	}
	defer conn.Close()
	return redis.String(conn.Do(c_ECHO_COMMAND + ` "testing"`))
}

// Info返回Redis信息和状态
func (that *DbiRedis) Info() (string, error) {
	conn := that.GetRedisConn()
	if conn == nil {
		return "", errors.New("no redis conn")
	}
	defer conn.Close()
	return redis.String(conn.Do(c_INFO_COMMAND))
}

// Scan数据
func (that *DbiRedis) Scan(startIndex int64, pattern string) (int64, []string, error) {
	conn := that.GetRedisConn()
	if conn == nil {
		return 0, nil, errors.New("no redis conn")
	}
	defer conn.Close()
	results, err := redis.Values(conn.Do(c_SCAN_COMMAND, startIndex, c_MATCH_OPTION, pattern))
	if err != nil {
		return 0, nil, err
	}
	return parseScanResults(results)
}

// Set sets a key/value pair
func (that *DbiRedis) Set(key string, value interface{}) (bool, error) {
	conn := that.GetRedisConn()
	if conn == nil {
		return false, errors.New("no redis conn")
	}
	defer conn.Close()
	return toBool(conn.Do(c_SET_COMMAND, that.Prefix+key, value))
}

// Append to a key's value
func (that *DbiRedis) Append(key string, value string) (int64, error) {
	conn := that.GetRedisConn()
	if conn == nil {
		return 0, errors.New("no redis conn")
	}
	defer conn.Close()
	return redis.Int64(conn.Do(c_APPEND_COMMAND, that.Prefix+key, value))
}

// GetRange to get a key's value's range
func (that *DbiRedis) GetRange(key string, start int, end int) (string, error) {
	conn := that.GetRedisConn()
	if conn == nil {
		return "", errors.New("no redis conn")
	}
	defer conn.Close()
	return redis.String(conn.Do(c_GET_RANGE_COMMAND, that.Prefix+key, start, end))
}

// SetRange to set a key's value's range
func (that *DbiRedis) SetRange(key string, start int, value string) (int64, error) {
	conn := that.GetRedisConn()
	if conn == nil {
		return 0, errors.New("no redis conn")
	}
	defer conn.Close()
	return redis.Int64(conn.Do(c_SET_RANGE_COMMAND, that.Prefix+key, start, value))
}

// Expire sets a key's timeout in seconds
func (that *DbiRedis) Expire(key string, timeout int64) (bool, error) {
	conn := that.GetRedisConn()
	if conn == nil {
		return false, errors.New("no redis conn")
	}
	defer conn.Close()
	count, err := redis.Int64(conn.Do(c_EXPIRE_COMMAND, that.Prefix+key, timeout))
	return count > 0, err
}

// SetNx sets a key/value pair if the key does not exist
func (that *DbiRedis) SetNx(key string, value string) (bool, error) {
	conn := that.GetRedisConn()
	if conn == nil {
		return false, errors.New("no redis conn")
	}
	defer conn.Close()
	return toBool(conn.Do(c_SET_COMMAND, that.Prefix+key, value, c_NOT_EXISTS_OPTION))
}

// SetEx sets a key/value pair with a timeout in seconds
func (that *DbiRedis) SetEx(key string, value string, timeout int64) (bool, error) {
	conn := that.GetRedisConn()
	if conn == nil {
		return false, errors.New("no redis conn")
	}
	defer conn.Close()
	return toBool(conn.Do(c_SET_COMMAND, that.Prefix+key, value, c_EXPIRE_OPTION, timeout))
}

// Get retrieves a key's value
func (that *DbiRedis) Get(key string) (string, bool, error) {
	conn := that.GetRedisConn()
	if conn == nil {
		return "", false, errors.New("no redis conn")
	}
	defer conn.Close()
	return toString(conn.Do(c_GET_COMMAND, that.Prefix+key))
}

// Exists checks how many keys exist
func (that *DbiRedis) Exists(keys ...string) (bool, error) {
	conn := that.GetRedisConn()
	if conn == nil {
		return false, errors.New("no redis conn")
	}
	defer conn.Close()
	args := make([]interface{}, len(keys))
	for i, key := range keys {
		args[i] = that.Prefix + key
	}
	count, err := redis.Int64(conn.Do(c_EXISTS_COMMAND, args...))
	return count > 0, err
}

// Del deletes keys
func (that *DbiRedis) Del(hasPrefix bool, keys ...string) (int64, error) {
	conn := that.GetRedisConn()
	if conn == nil {
		return 0, errors.New("no redis conn")
	}
	defer conn.Close()
	args := make([]interface{}, len(keys))
	for i, key := range keys {
		if hasPrefix {
			args[i] = key
		} else {
			args[i] = that.Prefix + key
		}
	}
	return redis.Int64(conn.Do(c_DEL_COMMAND, args...))
}

// Keys retrieves keys that match a pattern
func (that *DbiRedis) Keys(pattern string) ([]string, error) {
	conn := that.GetRedisConn()
	if conn == nil {
		return nil, errors.New("no redis conn")
	}
	defer conn.Close()
	return redis.Strings(conn.Do(c_KEYS_COMMAND, pattern))
}

// IncrBy increments the key's value by the increment provided
func (that *DbiRedis) IncrBy(key string, increment int64) (int64, error) {
	conn := that.GetRedisConn()
	if conn == nil {
		return 0, errors.New("no redis conn")
	}
	defer conn.Close()
	return redis.Int64(conn.Do(c_INCRBY_COMMAND, that.Prefix+key, increment))
}

// Incr increments the key's value
func (that *DbiRedis) Incr(key string) (int64, error) {
	return that.IncrBy(key, 1)
}

// IncrByFloat increments the key's value by the increment provided
func (that *DbiRedis) IncrByFloat(key string, increment float64) (float64, error) {
	conn := that.GetRedisConn()
	if conn == nil {
		return 0, errors.New("no redis conn")
	}
	defer conn.Close()
	return redis.Float64(conn.Do(c_INCRBYFLOAT_COMMAND, that.Prefix+key, increment))
}

// DecrBy decrements the key's value by the decrement provided
func (that *DbiRedis) DecrBy(key string, decrement int64) (int64, error) {
	return that.IncrBy(key, -decrement)
}

// Decr decrements the key's value
func (that *DbiRedis) Decr(key string) (int64, error) {
	return that.IncrBy(key, -1)
}

// DecrByFloat decrements the key's value by the decrement provided
func (that *DbiRedis) DecrByFloat(key string, decrement float64) (float64, error) {
	return that.IncrByFloat(key, -decrement)
}

// HScan incrementally iterate over key's fields and values
func (that *DbiRedis) HScan(key string, startIndex int64, pattern string, count int) (int64, []string, error) {
	conn := that.GetRedisConn()
	if conn == nil {
		return 0, nil, errors.New("no redis conn")
	}
	defer conn.Close()
	results, err := redis.Values(conn.Do(c_HSCAN_COMMAND, that.Prefix+key, startIndex, c_MATCH_OPTION, pattern, c_COUNT_OPTION, count))
	if err != nil {
		return 0, nil, err
	}
	return parseScanResults(results)
}

// HSet sets a key's field/value pair
func (that *DbiRedis) HSet(key string, field string, value string) (bool, error) {
	conn := that.GetRedisConn()
	if conn == nil {
		return false, errors.New("no redis conn")
	}
	defer conn.Close()
	code, err := redis.Int(conn.Do(c_HSET_COMMAND, that.Prefix+key, field, value))
	return code > 0, err
}

// HMSet sets a key's field/value pair map
func (that *DbiRedis) HMSet(key string, item map[string]interface{}) error {
	conn := that.GetRedisConn()
	if conn == nil {
		return errors.New("no redis conn")
	}
	defer conn.Close()
	reply, err := conn.Do(c_HMSET_COMMAND, redis.Args{}.Add(that.Prefix+key).AddFlat(item)...)
	if err != nil {
		return err
	}
	if reply != "OK" {
		return fmt.Errorf("reply string is wrong!: %s", reply)
	}
	return nil
}

// HKeys retrieves a hash's keys
func (that *DbiRedis) HKeys(key string) ([]string, error) {
	conn := that.GetRedisConn()
	if conn == nil {
		return nil, errors.New("no redis conn")
	}
	defer conn.Close()
	return redis.Strings(conn.Do(c_HKEYS_COMMAND, that.Prefix+key))
}

// HExists determine's a key's field's existence
func (that *DbiRedis) HExists(key string, field string) (bool, error) {
	conn := that.GetRedisConn()
	if conn == nil {
		return false, errors.New("no redis conn")
	}
	defer conn.Close()
	return redis.Bool(conn.Do(c_HEXISTS_COMMAND, that.Prefix+key, field))
}

// HExists determine's a key's field's existence
func (that *DbiRedis) HLen(key string) (int, error) {
	conn := that.GetRedisConn()
	if conn == nil {
		return 0, errors.New("no redis conn")
	}
	defer conn.Close()
	return redis.Int(conn.Do(c_HLEN_COMMAND, that.Prefix+key))
}

// HGet retrieves a key's field's value
func (that *DbiRedis) HGet(key string, field string) (string, bool, error) {
	conn := that.GetRedisConn()
	if conn == nil {
		return "", false, errors.New("no redis conn")
	}
	defer conn.Close()
	return toString(conn.Do(c_HGET_COMMAND, that.Prefix+key, field))
}

// HGetAll retrieves the key
func (that *DbiRedis) HGetAll(key string, tag map[string]*DbiTableTag) (map[string]interface{}, error) {
	conn := that.GetRedisConn()
	if conn == nil {
		return nil, errors.New("no redis conn")
	}
	defer conn.Close()
	reply, err := conn.Do(c_HGETALL_COMMAND, that.Prefix+key)
	return that.ConvertMap(reply, err, tag)
}

// HGetAll retrieves the key
func (that *DbiRedis) HGetAllRaw(key string) (map[string]interface{}, error) {
	conn := that.GetRedisConn()
	if conn == nil {
		return nil, errors.New("no redis conn")
	}
	defer conn.Close()
	reply, err := conn.Do(c_HGETALL_COMMAND, that.Prefix+key)
	return that.ConvertMapRaw(reply, err)
}

// HDel deletes a key's fields
func (that *DbiRedis) HDel(key string, fields ...interface{}) (int64, error) {
	conn := that.GetRedisConn()
	if conn == nil {
		return 0, errors.New("no redis conn")
	}
	defer conn.Close()
	args := append([]interface{}{that.Prefix + key}, fields...)
	return redis.Int64(conn.Do(c_HDEL_COMMAND, args...))
}

// HIncrBy increments the key's field's value by the increment provided
func (that *DbiRedis) HIncrBy(key string, field string, increment int64) (int64, error) {
	conn := that.GetRedisConn()
	if conn == nil {
		return 0, errors.New("no redis conn")
	}
	defer conn.Close()
	return redis.Int64(conn.Do(c_HINCRBY_COMMAND, that.Prefix+key, field, increment))
}

// HIncr increments the key's field's value
func (that *DbiRedis) HIncr(key string, field string) (int64, error) {
	return that.HIncrBy(key, field, 1)
}

// HIncrByFloat increments the key's field's value by the increment provided
func (that *DbiRedis) HIncrByFloat(key string, field string, increment float64) (float64, error) {
	conn := that.GetRedisConn()
	if conn == nil {
		return 0, errors.New("no redis conn")
	}
	defer conn.Close()
	return redis.Float64(conn.Do(c_HINCRBYFLOAT_COMMAND, that.Prefix+key, field, increment))
}

// HDecr decrements the key's field's value
func (that *DbiRedis) HDecr(key string, field string) (int64, error) {
	return that.HIncrBy(key, field, -1)
}

// HDecrBy decrements the key's field's value by the decrement provided
func (that *DbiRedis) HDecrBy(key string, field string, decrement int64) (int64, error) {
	return that.HIncrBy(key, field, -decrement)
}

// HDecrByFloat decrements the key's field's value by the decrement provided
func (that *DbiRedis) HDecrByFloat(key string, field string, decrement float64) (float64, error) {
	return that.HIncrByFloat(key, field, -decrement)
}

// 查询key过期时间
func (that *DbiRedis) GetTTL(key string) (int, error) {
	conn := that.GetRedisConn()
	if conn == nil {
		return 0, errors.New("no redis conn")
	}
	defer conn.Close()
	return redis.Int(conn.Do(c_TTL_COMMAND, that.Prefix+key))
}

// 加入有序集合（自顶向下，积分从小到大排列）
// score,member [score,member]
func (that *DbiRedis) ZAdd(key string, fields ...interface{}) (interface{}, error) {
	conn := that.GetRedisConn()
	if conn == nil {
		return 0, errors.New("no redis conn")
	}
	defer conn.Close()
	args := append([]interface{}{that.Prefix + key}, fields...)
	reply, err := conn.Do(c_ZADD_COMMAND, args...)
	return reply, err
}

// 获取积分
func (that *DbiRedis) ZScore(key string, member interface{}) (int64, error) {
	conn := that.GetRedisConn()
	if conn == nil {
		return 0, errors.New("no redis conn")
	}
	defer conn.Close()
	return redis.Int64(conn.Do(c_ZSCORE_COMMAND, that.Prefix+key, member))
}

// 获取区间积分
// min取值 (0 0 -inf
// max取值 (0 0 +inf
func (that *DbiRedis) ZRangeByScore(key string, min interface{}, max interface{}, limit ...int) ([]int64, error) {
	conn := that.GetRedisConn()
	if conn == nil {
		return nil, errors.New("no redis conn")
	}
	defer conn.Close()
	if len(limit) >= 2 {
		return redis.Int64s(conn.Do(c_ZRANGEBYSCORE_COMMAND, that.Prefix+key, min, max, "LIMIT", limit[0], limit[1]))
	} else {
		return redis.Int64s(conn.Do(c_ZRANGEBYSCORE_COMMAND, that.Prefix+key, min, max))
	}
}

// 删除有序集合成员
// member [member]
func (that *DbiRedis) ZRem(key string, fields ...interface{}) (int, error) {
	conn := that.GetRedisConn()
	if conn == nil {
		return 0, errors.New("no redis conn")
	}
	defer conn.Close()
	args := append([]interface{}{that.Prefix + key}, fields...)
	return redis.Int(conn.Do(c_ZREM_COMMAND, args...))
}

// 获取有序集合长度
func (that *DbiRedis) ZCard(key string) (int, error) {
	conn := that.GetRedisConn()
	if conn == nil {
		return 0, errors.New("no redis conn")
	}
	defer conn.Close()
	return redis.Int(conn.Do(c_ZCARD_COMMAND, that.Prefix+key))
}

// 向集合添加一个或多个成员
func (that *DbiRedis) SAdd(key string, fields ...interface{}) (interface{}, error) {
	conn := that.GetRedisConn()
	if conn == nil {
		return 0, errors.New("no redis conn")
	}
	defer conn.Close()
	args := append([]interface{}{that.Prefix + key}, fields...)
	reply, err := conn.Do(c_SADD_COMMAND, args...)
	return reply, err
}

// 获取集合的成员数
func (that *DbiRedis) SCard(key string) (int, error) {
	conn := that.GetRedisConn()
	if conn == nil {
		return 0, errors.New("no redis conn")
	}
	defer conn.Close()
	return redis.Int(conn.Do(c_SCARD_COMMAND, that.Prefix+key))
}

// 删除集合成员
// member [member]
func (that *DbiRedis) SRem(key string, fields ...interface{}) (int, error) {
	conn := that.GetRedisConn()
	if conn == nil {
		return 0, errors.New("no redis conn")
	}
	defer conn.Close()
	args := append([]interface{}{that.Prefix + key}, fields...)
	return redis.Int(conn.Do(c_SREM_COMMAND, args...))
}

// 移除并返回集合中的一个随机元素
func (that *DbiRedis) SPop(key string) (string, bool, error) {
	conn := that.GetRedisConn()
	if conn == nil {
		return "", false, errors.New("no redis conn")
	}
	defer conn.Close()
	return toString(conn.Do(c_SPOP_COMMAND, that.Prefix+key))
}

// 迭代集合中的元素
func (that *DbiRedis) SScan(key string, startIndex int64, pattern string) (int64, []string, error) {
	conn := that.GetRedisConn()
	if conn == nil {
		return 0, nil, errors.New("no redis conn")
	}
	defer conn.Close()
	results, err := redis.Values(conn.Do(c_SSCAN_COMMAND, that.Prefix+key, startIndex, c_MATCH_OPTION, pattern))
	if err != nil {
		return 0, nil, err
	}
	return parseScanResults(results)
}

func (that *DbiRedis) Pipeline() {
	conn := that.GetRedisConn()
	if conn == nil {
		return
	}
	defer conn.Close()
}

// StringMap is a helper that converts an array of strings (alternating key, value)
// into a map[string]string. The HGETALL and CONFIG GET commands return replies in this format.
// Requires an even number of values in result.
func (that *DbiRedis) ConvertMap(result interface{}, err error, tag map[string]*DbiTableTag) (map[string]interface{}, error) {
	values, err := redis.Values(result, err)
	if err != nil {
		return nil, err
	}

	if len(values)%2 != 0 {
		return nil, fmt.Errorf("redigo: StringMap expects even number of values result, got %d", len(values))
	}

	m := make(map[string]interface{}, len(values)/2)
	for i := 0; i < len(values); i += 2 {
		key, ok := values[i].([]byte)
		if !ok {
			return nil, fmt.Errorf("redigo: StringMap key[%d] not a bulk string value, got %T", i, values[i])
		}
		tagKey := string(key)
		tagVal := tag[tagKey]
		if tagVal == nil {
			continue
		}
		value := values[i+1]
		switch tagVal.Kind {
		case reflect.Int8:
			value8, _ := redis.Int(value, nil)
			value = int8(value8)
		case reflect.Int16:
			value16, _ := redis.Int(value, nil)
			value = int16(value16)
		case reflect.Int32:
			value32, _ := redis.Int(value, nil)
			value = int32(value32)
		case reflect.Int64:
			value, _ = redis.Int64(value, nil)
		case reflect.Int:
			value, _ = redis.Int(value, nil)
		case reflect.Uint8:
			value8, _ := redis.Int(value, nil)
			value = uint8(value8)
		case reflect.Uint16:
			value16, _ := redis.Int(value, nil)
			value = uint16(value16)
		case reflect.Uint32:
			value32, _ := redis.Int(value, nil)
			value = uint32(value32)
		case reflect.Uint64:
			value, _ = redis.Uint64(value, nil)
		case reflect.Float32:
			value64, _ := redis.Float64(value, nil)
			value = float32(value64)
		case reflect.Float64:
			value, _ = redis.Float64(value, nil)
		case reflect.String:
			value, _ = redis.String(value, nil)
		default:
			continue
		}
		m[tagKey] = value
	}
	return m, nil
}

// StringMap is a helper that converts an array of strings (alternating key, value)
// into a map[string]string. The HGETALL and CONFIG GET commands return replies in this format.
// Requires an even number of values in result.
func (that *DbiRedis) ConvertMapRaw(result interface{}, err error) (map[string]interface{}, error) {
	values, err := redis.Values(result, err)
	if err != nil {
		return nil, err
	}

	if len(values)%2 != 0 {
		return nil, fmt.Errorf("redigo: StringMap expects even number of values result, got %d", len(values))
	}

	m := make(map[string]interface{}, len(values)/2)
	for i := 0; i < len(values); i += 2 {
		key, ok := values[i].([]byte)
		if !ok {
			return nil, fmt.Errorf("redigo: StringMap key[%d] not a bulk string value, got %T", i, values[i])
		}
		tagKey := string(key)
		tagVal, _ := redis.String(values[i+1], nil)
		m[tagKey] = tagVal
	}
	return m, nil
}

// 连接字符串
func (that *DbiRedis) Redis_JoinKey(params ...interface{}) string {
	tmp := make([]string, len(params))
	for i := 0; i < len(params); i++ {
		tmp[i] = fmt.Sprint(params[i])
	}
	return strings.Join(tmp, ":")
}

// 获取加锁字符串
func (that *DbiRedis) Redis_JoinLockKey(params ...interface{}) string {
	tmp := make([]string, len(params))
	for i := 0; i < len(params); i++ {
		tmp[i] = fmt.Sprint(params[i])
	}
	return "lock:" + strings.Join(tmp, ":")
}
