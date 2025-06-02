package email

import (
	"context"
	"errors"
	"go-job/internal/pkg/ratelimit"
)

/*
邮件限流器
*/

var errSendEmailToMany = errors.New("send email to many emails")

type EmailLimiter struct {
	key     string
	limiter ratelimit.Limiter
	svc     IEmailService
}

func NewEmailLimiter(limiter ratelimit.Limiter, svc IEmailService) *EmailLimiter {
	return &EmailLimiter{
		key:     "email-limiter",
		limiter: limiter,
		svc:     svc,
	}
}

func (l *EmailLimiter) SetKey(prefix string) *EmailLimiter {
	l.key = prefix
	return l
}

func (l *EmailLimiter) Send(ctx context.Context, email []string, subject, content string) error {
	isLimit, err := l.limiter.Limit(ctx, l.key)
	if err != nil {
		return err
	}
	if isLimit {
		return errSendEmailToMany
	}
	return l.Send(ctx, email, subject, content)
}
