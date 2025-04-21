package cache

import "context"

// EmailCodeCache 通过缓存来实现验证码的存储和校验
type EmailCodeCache struct {
}

func (Q *EmailCodeCache) Set(ctx context.Context, biz, email, code string) error {
	//TODO implement me
	panic("implement me")
}

func (Q *EmailCodeCache) Verify(ctx context.Context, biz, email, code string) (err error) {
	//TODO implement me
	panic("implement me")
}
