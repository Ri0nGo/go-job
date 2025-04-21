package cache

import "context"

// EmailCodeRepo
// @Description: 邮箱验证码，屏蔽下层实现，可以通过缓存实现，也可以通过db实现，甚至可以通过文件来实现
type IEmailCodeCache interface {
	Set(ctx context.Context, biz, email, code string) error
	Verify(ctx context.Context, biz, email, code string) (err error)
}
