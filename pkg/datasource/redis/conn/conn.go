package conn

import (
	"context"
	"time"

	"sonymimic1/Golang_server/checkRTP/config"

	"github.com/redis/go-redis/v9"
)

// SetupRedisConnection is creating a new connection to our cache db.
func SetupRedisConnection(ctx context.Context, cfg config.RedisConfig) *redis.ClusterClient {

	setting := redis.ClusterOptions{
		ReadTimeout:  time.Duration(cfg.ReadTimeout) * time.Second,
		WriteTimeout: time.Duration(cfg.WriteTimeout) * time.Second,
		MaxRetries:   cfg.MaxRetries,
		DialTimeout:  time.Duration(cfg.DialTimeout) * time.Second,
		PoolSize:     cfg.PoolSize,
		PoolTimeout:  time.Duration(cfg.PoolTimeout) * time.Second,
		Addrs:        cfg.Hosts,
	}

	redisClient := redis.NewClusterClient(&setting)
	//檢查redisClient是否正常
	pong, err := redisClient.Ping(ctx).Result()
	if err != nil {
		println("redisClient ping fail !")
	} else {
		println(pong)
	}

	return redisClient
}

// CloseRedisConnection method is closing a connection between your app and your cache db
func CloseRedisConnection(client *redis.ClusterClient) {

	if err := client.Close(); err != nil {
		println("cachedb close fail !")
	}

}
