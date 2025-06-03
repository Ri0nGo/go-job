package email

import (
	"github.com/redis/go-redis/v9"
	"go-job/internal/iface/email"
	"go-job/internal/pkg/email/failover"
	"go-job/internal/pkg/email/netease163"
	"go-job/internal/pkg/email/qq"
	emailRatelimit "go-job/internal/pkg/email/ratelimit"
	"go-job/internal/pkg/ratelimit"
	"go-job/master/pkg/config"
	"time"
)

const (
	prefix            = "email-limiter-"
	qqEmailSmtpKey    = "qq"
	netease163SmtpKey = "netease163"
)

type EmailRateLimitConfig struct {
	Interval time.Duration
	Rate     int
}

var EmailLimitConfig = map[string]map[string]EmailRateLimitConfig{
	qqEmailSmtpKey: {
		"min": {
			Interval: time.Minute,
			Rate:     30,
		},
		"day": {
			Interval: time.Hour * 24,
			Rate:     150,
		},
	},
	netease163SmtpKey: {
		"min": {
			Interval: time.Minute,
			Rate:     5,
		},
		"day": {
			Interval: time.Hour * 24,
			Rate:     100,
		},
	},
}

func InitEmailService(cmd redis.Cmdable) email.IEmailService {
	return initEmailService(cmd)
}

// initEmailService 初始化邮箱服务，采用装饰器模式来增强短信服务的可用性
// 区分不同服务的限流策略
func initEmailService(cmd redis.Cmdable) email.IEmailService {
	var svcs []email.IEmailService
	if smtpCfg, ok := config.App.SMTP[qqEmailSmtpKey]; ok {
		svc := qq.NewQQEmailService(qqEmailSmtpKey,
			smtpCfg.Key,
			smtpCfg.Sender,
			smtpCfg.SMTPHost,
			smtpCfg.SMTPPort)
		limiter := emailRatelimit.NewEmailRatelimit(svc)
		for period, conf := range EmailLimitConfig[qqEmailSmtpKey] {
			key := prefix + period + ":" + svc.Name()
			rl := ratelimit.NewRedisSlidingWindowLimiter(cmd, conf.Interval, conf.Rate)
			limiter.RegistryLimiter(key, rl)
		}
		svcs = append(svcs, limiter)
	}
	if smtpCfg, ok := config.App.SMTP[netease163SmtpKey]; ok {
		svc := netease163.NewNetEase163EmailService(netease163SmtpKey,
			smtpCfg.Key,
			smtpCfg.Sender,
			smtpCfg.SMTPHost,
			smtpCfg.SMTPPort)
		limiter := emailRatelimit.NewEmailRatelimit(svc)
		for period, conf := range EmailLimitConfig[netease163SmtpKey] {
			key := prefix + period + ":" + svc.Name()
			rl := ratelimit.NewRedisSlidingWindowLimiter(cmd, conf.Interval, conf.Rate)
			limiter.RegistryLimiter(key, rl)
		}
		svcs = append(svcs, limiter)
	}

	if len(svcs) == 0 {
		panic("smtp svc is empty")
	}

	return failover.NewEmailPollFailOver(svcs) // 故障转移
}
