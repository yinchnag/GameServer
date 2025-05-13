package database

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
)

func newDbiRedisV8(ip string, dbidx int, auth string) *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     ip,
		Password: auth,
		DB:       dbidx,
		// PoolSize:     Max_Redis_Active_Conn,
		// MinIdleConns: Max_Redis_Idle_Conn,
		// IdleTimeout:  Max_Redis_Idle_Time * time.Second,
	})
}

type DbiRedisV8 struct {
	clt *redis.Client
}

func (this *DbiRedisV8) Init(ip string, dbidx int, auth string) {
	this.clt = newDbiRedisV8(ip, dbidx, auth)
}

func (that *DbiRedisV8) Get(key string) (string, error) {
	ctx := context.Background()
	val, err := that.clt.Get(ctx, key).Result()
	if err != nil {
		return "", err
	}
	return val, nil
}

func (that *DbiRedisV8) Set(key string, value interface{}) (bool, error) {
	ctx := context.Background()
	res, err := that.clt.Set(ctx, key, value, 0).Result()

	return res == "OK", err
}

func (that *DbiRedisV8) Append(ctx context.Context, key string, value string) (int64, error) {
	len, err := that.clt.Append(ctx, key, value).Result()
	return len, err
}

func (that *DbiRedisV8) GetRange(ctx context.Context, key string, start, end int64) (string, error) {
	return that.clt.GetRange(ctx, key, start, end).Result()
}

func (that *DbiRedisV8) SetRange(ctx context.Context, key string, start int64, value string) (int64, error) {
	return that.clt.SetRange(ctx, key, start, value).Result()
}

func (that *DbiRedisV8) Expire(ctx context.Context, key string, timeout time.Duration) (bool, error) {
	return that.clt.Expire(ctx, key, timeout).Result()
}

func (that *DbiRedisV8) SetNx(ctx context.Context, key string, value string) (bool, error) {
	return that.clt.SetNX(ctx, key, value, 0).Result()
}

func (that *DbiRedisV8) SetEx(ctx context.Context, key string, value interface{}, timeout time.Duration) (bool, error) {
	result, err := that.clt.SetEX(ctx, key, value, timeout).Result()
	return result == "OK", err
}

func (that *DbiRedisV8) Exists(ctx context.Context, keys ...string) (bool, error) {
	result, err := that.clt.Exists(ctx, keys...).Result()
	return result > 0, err
}

func (that *DbiRedisV8) Del(ctx context.Context, hasPrefix bool, keys ...string) (int64, error) {
	return that.clt.Del(ctx, keys...).Result()
}

func (that *DbiRedisV8) Keys(ctx context.Context, pattern string) ([]string, error) {
	return that.clt.Keys(ctx, pattern).Result()
}

func (that *DbiRedisV8) IncrBy(ctx context.Context, key string, increment int64) (int64, error) {
	return that.clt.IncrBy(ctx, key, increment).Result()
}

func (that *DbiRedisV8) Incr(ctx context.Context, key string) (int64, error) {
	return that.clt.Incr(ctx, key).Result()
}

func (that *DbiRedisV8) IncrByFloat(key string, increment float64) (float64, error) {
	return that.clt.IncrByFloat(context.Background(), key, increment).Result()
}

func (that *DbiRedisV8) DecrBy(ctx context.Context, key string, decrement int64) (int64, error) {
	return that.clt.DecrBy(ctx, key, decrement).Result()
}

func (that *DbiRedisV8) Decr(ctx context.Context, key string) (int64, error) {
	return that.clt.Decr(ctx, key).Result()
}

func (that *DbiRedisV8) DecrByFloat(ctx context.Context, key string, decrement float64) (float64, error) {
	return that.clt.IncrByFloat(ctx, key, decrement).Result()
}

func (that *DbiRedisV8) HScan(ctx context.Context, key string, startIndex uint64, pattern string, count int64) ([]string, uint64, error) {
	return that.clt.HScan(ctx, key, startIndex, pattern, count).Result()
}

func (that *DbiRedisV8) HSet(ctx context.Context, key string, field string, value string) (bool, error) {
	count, err := that.clt.HSet(ctx, key, field, value).Result()
	return count > 0, err
}

func (that *DbiRedisV8) HMSet(ctx context.Context, key string, item map[string]interface{}) error {
	return that.clt.HMSet(ctx, key, item).Err()
}

func (that *DbiRedisV8) HKeys(ctx context.Context, key string) ([]string, error) {
	return that.clt.HKeys(ctx, key).Result()
}

func (that *DbiRedisV8) HExists(ctx context.Context, key string, field string) (bool, error) {
	return that.clt.HExists(ctx, key, field).Result()
}

func (that *DbiRedisV8) HLen(ctx context.Context, key string) (int64, error) {
	return that.clt.HLen(ctx, key).Result()
}

func (that *DbiRedisV8) HGet(ctx context.Context, key string, field string) (string, error) {
	return that.clt.HGet(ctx, key, field).Result()
}

func (that *DbiRedisV8) HGetAll(ctx context.Context, key string, tag map[string]*DbiTableTag) (map[string]string, error) {
	return that.clt.HGetAll(ctx, key).Result()
}

func (that *DbiRedisV8) HGetAllRaw(ctx context.Context, key string) (map[string]string, error) {
	return that.clt.HGetAll(ctx, key).Result()
}

func (that *DbiRedisV8) HDel(ctx context.Context, key string, fields ...string) (int64, error) {
	return that.clt.HDel(ctx, key, fields...).Result()
}

func (that *DbiRedisV8) HIncrBy(ctx context.Context, key string, field string, increment int64) (int64, error) {
	return that.clt.HIncrBy(ctx, key, field, increment).Result()
}

func (that *DbiRedisV8) HIncr(ctx context.Context, key string, field string) (int64, error) {
	return that.clt.HIncrBy(ctx, key, field, 1).Result()
}

func (that *DbiRedisV8) HIncrByFloat(ctx context.Context, key string, field string, increment float64) (float64, error) {
	return that.clt.HIncrByFloat(ctx, key, field, increment).Result()
}

func (that *DbiRedisV8) HDecr(ctx context.Context, key string, field string) (int64, error) {
	return that.clt.HIncrBy(ctx, key, field, -1).Result()
}

func (that *DbiRedisV8) HDecrBy(ctx context.Context, key string, field string, decrement int64) (int64, error) {
	return that.clt.HIncrBy(ctx, key, field, -decrement).Result()
}

func (that *DbiRedisV8) HDecrByFloat(ctx context.Context, key string, field string, decrement float64) (float64, error) {
	return that.HIncrByFloat(ctx, key, field, -decrement)
}

func (that *DbiRedisV8) GetTTL(ctx context.Context, key string) (time.Duration, error) {
	return that.clt.TTL(ctx, key).Result()
}

func (that *DbiRedisV8) ZAdd(ctx context.Context, key string, fields ...interface{}) (interface{}, error) {
	members := []*redis.Z{}
	for field := range fields {
		members = append(members, &redis.Z{
			Score:  0,
			Member: field,
		})
	}
	return that.clt.ZAdd(ctx, key, members...).Result()
}

func (that *DbiRedisV8) ZScore(ctx context.Context, key string, member string) (float64, error) {
	return that.clt.ZScore(ctx, key, member).Result()
}

func (that *DbiRedisV8) ZRangeByScore(ctx context.Context, key string, min string, max string, limit ...int64) ([]string, error) {
	return that.clt.ZRangeByScore(ctx, key, &redis.ZRangeBy{Min: min, Max: max, Offset: limit[0], Count: limit[1]}).Result()
}

func (that *DbiRedisV8) ZRem(ctx context.Context, key string, members ...string) (int64, error) {
	src := []interface{}{}
	for i := 0; i < len(members); i++ {
		src[i] = members[i]
	}
	return that.clt.ZRem(ctx, key, src...).Result()
}

func (that *DbiRedisV8) ZCard(ctx context.Context, key string) (int64, error) {
	return that.clt.ZCard(ctx, key).Result()
}

func (that *DbiRedisV8) SAdd(ctx context.Context, key string, members ...string) (int64, error) {
	src := []interface{}{}
	for i := 0; i < len(members); i++ {
		src[i] = members[i]
	}
	return that.clt.SAdd(ctx, key, src...).Result()
}

func (that *DbiRedisV8) SCard(ctx context.Context, key string) (int64, error) {
	return that.clt.SCard(ctx, key).Result()
}

func (that *DbiRedisV8) SRem(ctx context.Context, key string, members ...string) (int64, error) {
	src := []interface{}{}
	for i := 0; i < len(members); i++ {
		src[i] = members[i]
	}
	return that.clt.SRem(ctx, key, src...).Result()
}

func (that *DbiRedisV8) SScan(ctx context.Context, key string, cursor uint64, match string, count int64) ([]string, uint64, error) {
	return that.clt.SScan(ctx, key, cursor, match, count).Result()
}
