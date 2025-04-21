package email

import "context"

// EmailService
// @Description: 发生邮件的抽象，针对不同的SMTP服务商
type IEmailService interface {
	Send(ctx context.Context, email []string, subject, content string) error
}
