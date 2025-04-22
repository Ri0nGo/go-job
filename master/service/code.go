package service

import (
	"context"
	"fmt"
	"go-job/internal/pkg/email"
	"go-job/master/repo"
	"log/slog"
	"math/rand"
)

// EmailCodeService
// @Description: 邮箱验证码服务
type IEmailCodeService interface {
	Send(ctx context.Context, biz, email string) error
	Verify(ctx context.Context, biz, email, code string) error
}

type EmailCodeService struct {
	emailSvc  email.IEmailService
	emailRepo repo.IEmailCodeRepo
}

func (s *EmailCodeService) Send(ctx context.Context, biz, emailStr string) error {
	code := s.generateCode()
	if err := s.emailRepo.Set(ctx, biz, emailStr, code); err != nil {
		slog.Error("email set error", "err", err)
		return err
	}
	tpl := email.GetEmailTpl(email.EmailBindVerfyCodeTpl)
	return s.emailSvc.Send(ctx, []string{emailStr}, tpl.Subject, fmt.Sprintf(tpl.Content, emailStr, code))
}

func (s *EmailCodeService) generateCode() string {
	code := rand.Intn(1000000)
	return fmt.Sprintf("%06d", code)
}

func (s *EmailCodeService) Verify(ctx context.Context, biz, email, code string) error {
	return s.emailRepo.Verify(ctx, biz, email, code)
}

func NewEmailCodeService(emailSvc email.IEmailService, emailRepo repo.IEmailCodeRepo) IEmailCodeService {
	return &EmailCodeService{
		emailSvc:  emailSvc,
		emailRepo: emailRepo,
	}
}
