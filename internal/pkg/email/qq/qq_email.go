package qq

import (
	"context"
	"go-job/internal/iface/email"
	"gopkg.in/gomail.v2"
)

/*
发送邮箱依赖于QQ邮箱
*/

type QQEmailService struct {
	name     string
	key      string // 授权码
	sender   string
	smtpHost string
	smtpPort int
}

func (e *QQEmailService) Send(ctx context.Context, emails []string, subject, content string) error {
	msg := gomail.NewMessage()
	msg.SetHeader("From", e.sender)
	msg.SetHeader("To", emails...)
	msg.SetHeader("Subject", subject)
	msg.SetBody("text/html", content)
	n := gomail.NewDialer(e.smtpHost, e.smtpPort, e.sender, e.key)
	if err := n.DialAndSend(msg); err != nil {
		return err
	}
	return nil
}

func (e *QQEmailService) Name() string {
	return e.name
}

func NewQQEmailService(name, key, sender, host string, port int) email.IEmailService {
	return &QQEmailService{
		name:     name,
		key:      key,
		sender:   sender,
		smtpHost: host,
		smtpPort: port,
	}
}
