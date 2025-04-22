package cache

import (
	"context"
	_ "embed"
	"errors"
	"fmt"
	"github.com/redis/go-redis/v9"
)

var (
	ErrCodeVerifyTooMany = errors.New("验证太频繁")
	ErrCodeSendTooMany   = errors.New("发送验证码太频繁")
	ErrInputCodeInvalid  = errors.New("验证码错误")

	//go:embed lua/set_email_code.lua
	luaSetEmailCode string

	//go:embed lua/verify_email_code.lua
	luaVerifyEmailCode string
)

// EmailCodeCache 通过缓存来实现邮箱验证码的存储和校验
type EmailCodeCache struct {
	redisCache redis.Cmdable
}

func (c *EmailCodeCache) Set(ctx context.Context, biz, email, code string) error {
	res, err := c.redisCache.Eval(ctx, luaSetEmailCode, []string{c.key(biz, email)}, code).Int()
	if err != nil {
		return err
	}
	switch res {
	case 0:
		return nil
	case -1:
		return ErrCodeSendTooMany
	default:
		return errors.New(fmt.Sprintf("key: %s don't expire time", c.key(biz, email)))
	}
}

func (c *EmailCodeCache) key(biz, email string) string {
	return fmt.Sprintf("email_code:%s:%s", biz, email)

}

func (c *EmailCodeCache) Verify(ctx context.Context, biz, email, code string) error {
	res, err := c.redisCache.Eval(ctx, luaVerifyEmailCode, []string{c.key(biz, email)}, code).Int()
	if err != nil {
		return nil
	}
	switch res {
	case 0:
		return nil
	case -1:
		return ErrCodeVerifyTooMany
	default:
		return ErrInputCodeInvalid
	}
}

func NewEmailCodeCache(redis redis.Cmdable) IEmailCodeCache {
	return &EmailCodeCache{
		redisCache: redis,
	}
}
