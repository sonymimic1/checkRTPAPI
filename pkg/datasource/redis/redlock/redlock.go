package redlock

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

const (
	//透過lua腳本執行，達到原子性操作
	luaDelete = `
			if redis.call('GET', KEYS[1]) == ARGV[1] then
				return redis.call('DEL', KEYS[1])
			else
				return 0
			end
		`
)

type RedLock struct {
	redisClient *redis.ClusterClient
}

func NewRediLock(redisClient *redis.ClusterClient) *RedLock {
	return &RedLock{
		redisClient: redisClient,
	}
}

func (lock *RedLock) RediLock(ctx context.Context, key string, value string, ttl time.Duration) (bool, error) {

	// 隨機值 uuid. or used snowflake.
	//lock.value = uuid.New().String()

	//設定key值並做retry機制
	{
		retry := 0

		for {
			set, err := lock.redisClient.SetNX(ctx, key, value, ttl).Result()

			if err != nil {
				panic(err.Error())
			}
			if set == true {
				return true, nil
			}

			if retry >= 5 {
				return false, nil
			}
			retry++
			time.Sleep(time.Millisecond * 50)
		}
	}

}

func (lock *RedLock) RediUnLock(ctx context.Context, key string, value string) (bool, error) {

	result, err := lock.redisClient.Eval(ctx, luaDelete, []string{key}, []string{value}).Result()
	if err != nil {
		return false, err
	}
	if result != int64(1) {
		return false, err
	}
	return true, nil
}
