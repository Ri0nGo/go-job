package ratelimit

import "context"

type Limiter interface {
	// Limit bool 表示是否需要限流
	// error 表示限流器是否有错误
	Limit(ctx context.Context, key string) (bool, error)
}
