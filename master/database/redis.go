package database

import (
	"context"
	"github.com/redis/go-redis/v9"
	"go-job/master/pkg/config"
	"log/slog"
)

func NewRedisClient() redis.Cmdable {
	client := redis.NewClient(&redis.Options{
		Addr:     config.App.Redis.Addr,
		Password: config.App.Redis.Auth, // 没有密码，默认值
		DB:       config.App.Redis.DB,   // 默认DB 0
	})
	ping := client.Ping(context.Background())
	if err := ping.Err(); err != nil {
		panic(err)
	}
	slog.Info("redis connect success")
	return client
}
