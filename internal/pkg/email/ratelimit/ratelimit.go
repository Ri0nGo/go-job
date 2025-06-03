package ratelimit

import (
	"context"
	"errors"
	"go-job/internal/iface/email"
	"go-job/internal/pkg/ratelimit"
)

/*
邮件限流器
*/

var errSendEmailToMany = errors.New("send email to many emails")

type EmailRatelimit struct {
	limiter map[string]ratelimit.Limiter
	svc     email.IEmailService
}

func NewEmailRatelimit(svc email.IEmailService) *EmailRatelimit {
	return &EmailRatelimit{
		limiter: make(map[string]ratelimit.Limiter),
		svc:     svc,
	}
}

func (l *EmailRatelimit) RegistryLimiter(key string, limiter ratelimit.Limiter) *EmailRatelimit {
	if _, ok := l.limiter[key]; ok {
		panic("registry limiter duplicated, key: " + key)
	}
	l.limiter[key] = limiter
	return l
}

func (l *EmailRatelimit) Send(ctx context.Context, email []string, subject, content string) error {
	for key, limiter := range l.limiter {
		isLimit, err := limiter.Limit(ctx, key)
		if err != nil {
			return err
		}
		if isLimit {
			return errSendEmailToMany
		}
	}

	return l.svc.Send(ctx, email, subject, content)
}

func (l *EmailRatelimit) Name() string {
	return l.svc.Name()
}
