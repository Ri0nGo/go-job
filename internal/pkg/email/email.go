package email

import "go-job/master/pkg/config"

func InitEmailService() IEmailService {
	return initEmailService()
}

func initEmailService() IEmailService {
	return NewQQEmailService(config.App.SMTP.Key,
		config.App.SMTP.Sender,
		config.App.SMTP.SMTPHost,
		config.App.SMTP.SMTPPort)
}
