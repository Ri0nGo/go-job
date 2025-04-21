package service

import (
	"context"
	"fmt"
	"go-job/internal/pkg/email"
	"go-job/master/repo"
	"math/rand"
)

// EmailCodeService
// @Description: 邮箱验证码服务
type IEmailCodeService interface {
	Send(ctx context.Context, biz, email string) error
	Verify(ctx context.Context, biz, email, code string) (err error)
}

type EmailCodeService struct {
	emailSvc  email.EmailService
	emailRepo repo.IEmailCodeRepo
}

func (s *EmailCodeService) Send(ctx context.Context, biz, emailStr string) error {
	code := s.generateCode()
	if err := s.emailRepo.Set(ctx, biz, emailStr, code); err != nil {
		return err
	}
	tpl := email.GetEmailTpl(email.EmailBindVerfyCodeTpl)
	return s.emailSvc.Send(ctx, code, tpl.Subject, fmt.Sprintf(tpl.Content, emailStr, code))
}

func (s *EmailCodeService) generateCode() string {
	code := rand.Intn(1000000)
	return fmt.Sprintf("%06d", code)
}

func (s *EmailCodeService) Verify(ctx context.Context, biz, email, code string) (err error) {
	//TODO implement me
	panic("implement me")
}
