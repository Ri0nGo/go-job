package ratelimit

import (
	"context"
	_ "embed"
	"github.com/redis/go-redis/v9"
	"time"
)

//go:embed lua/slide_window.lua
var luaSlideWindowScript string

type RedisSlidingWindowLimiter struct {
	cmd      redis.Cmdable
	interval time.Duration // 间隔
	rate     int           // 单位间隔内允许的最大个数
}

func NewRedisSlidingWindowLimiter(cmd redis.Cmdable, interval time.Duration, rate int) *RedisSlidingWindowLimiter {
	return &RedisSlidingWindowLimiter{
		cmd:      cmd,
		interval: interval,
		rate:     rate,
	}
}

func (l *RedisSlidingWindowLimiter) Limit(ctx context.Context, key string) (bool, error) {
	return l.cmd.Eval(ctx, luaSlideWindowScript, []string{key},
		l.interval.Milliseconds(), l.rate, time.Now().UnixMilli()).Bool()
}
