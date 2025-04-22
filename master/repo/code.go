package repo

import (
	"context"
	"go-job/master/repo/cache"
)

// IEmailCodeRepo
// @Description: 解耦处理，后续可以支持注入其他的验证码存储方式，比如Memcached等
type IEmailCodeRepo interface {
	Set(ctx context.Context, biz, email, code string) error
	Verify(ctx context.Context, biz, email, code string) error
}

type EmailCodeRepo struct {
	codeCache cache.IEmailCodeCache
}

func (e *EmailCodeRepo) Set(ctx context.Context, biz, email, code string) error {
	return e.codeCache.Set(ctx, biz, email, code)
}

func (e *EmailCodeRepo) Verify(ctx context.Context, biz, email, code string) error {
	return e.codeCache.Verify(ctx, biz, email, code)
}

func NewEmailCodeRepo(codeCache cache.IEmailCodeCache) IEmailCodeRepo {
	return &EmailCodeRepo{codeCache}
}
