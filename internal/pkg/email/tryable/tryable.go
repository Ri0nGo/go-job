package tryable

import (
	"context"
	"fmt"
	"go-job/internal/iface/email"
	"log/slog"
)

type TryAble struct {
	maxCnt uint32
	svc    email.IEmailService
}

func (t *TryAble) Send(ctx context.Context, email []string, subject, content string) error {
	for i := 0; i < int(t.maxCnt); i++ {
		err := t.svc.Send(ctx, email, subject, content)
		if err == nil {
			return nil
		}
		slog.Error("send email fail with tryable", "err", err)
	}
	return fmt.Errorf("try %d times all failed", t.maxCnt)
}

func (t *TryAble) Name() string {
	return t.svc.Name()
}

func NewTryAble(maxCnt uint32, svc email.IEmailService) *TryAble {
	return &TryAble{
		maxCnt: maxCnt,
		svc:    svc,
	}
}
