package failover

import (
	"context"
	"errors"
	"go-job/internal/iface/email"
	"log/slog"
	"sync/atomic"
)

/*
轮询实现邮件发送
*/

type EmailPollFailOver struct {
	idx  uint64 // 当前使用的节点
	svcs []email.IEmailService
}

func NewEmailPollFailOver(svcs []email.IEmailService) *EmailPollFailOver {
	return &EmailPollFailOver{
		svcs: svcs,
	}
}

func (e *EmailPollFailOver) Send(ctx context.Context, email []string, subject, content string) error {
	idx := atomic.AddUint64(&e.idx, 1)
	length := uint64(len(e.svcs))
	for i := idx; i < idx+length; i++ {
		svc := e.svcs[i%length]
		err := svc.Send(ctx, email, subject, content)
		if err == nil { // 只要有错误就切换服务商
			return nil
		}
		// todo 这里可以添加监控
		slog.Error("send email fail", "name", svc.Name(), "err", err)
	}

	// todo 这里可以发出告警
	return errors.New("all svc unavailable")
}

func (e *EmailPollFailOver) Name() string {
	return "email-failover"
}
