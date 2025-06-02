package email

import (
	"github.com/redis/go-redis/v9"
	"go-job/internal/pkg/ratelimit"
	"go-job/master/pkg/config"
	"time"
)

// 1秒内最多发送50封邮件
const (
	defaultLimitInterval = time.Second * 1
	defaultLimitRate     = 50
)

func InitEmailService(cmd redis.Cmdable) IEmailService {
	return initEmailService(cmd)
}

func initEmailService(cmd redis.Cmdable) IEmailService {
	limiter := ratelimit.NewRedisSlidingWindowLimiter(cmd, defaultLimitInterval, defaultLimitRate)
	svc := NewQQEmailService(config.App.SMTP.Key,
		config.App.SMTP.Sender,
		config.App.SMTP.SMTPHost,
		config.App.SMTP.SMTPPort)
	return NewEmailLimiter(limiter, svc)
}
