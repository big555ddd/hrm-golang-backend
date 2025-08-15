package config

import (
	"app/internal/logger"
	"context"

	"github.com/redis/go-redis/v9"
)

var RedisClient *redis.Client

func Redis() {
	RedisClient = redis.NewClient(&redis.Options{
		Addr:     confString("REDIS_ADDR", "localhost:6379"),
		Password: confString("REDIS_PASSWORD", ""),
		DB:       int(confInt64("REDIS_DB", 0)),
	})

	if err := RedisClient.Ping(context.Background()).Err(); err != nil {
		panic("failed to connect to redis: ")
	}
	logger.Info("redis connected success")
}

func GetRedis() *redis.Client {
	return RedisClient
}
