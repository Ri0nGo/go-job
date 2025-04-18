package email

import (
	"context"
	"go-job/master/pkg/config"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func getWorkSpace() string {
	wk := os.Getenv("WORKSPACE")
	if wk == "" {
		dir, _ := os.Getwd()
		absPath, _ := filepath.Abs(dir)
		return absPath
	}
	return wk
}

func getMasterConfigPath() string {
	return filepath.Join(getWorkSpace(), "config", "master-prod.yaml")

}

func TestNewQQEmailService(t *testing.T) {
	testEmails := []string{
		"*****@qq.com",
	}
	configPath := getMasterConfigPath()
	config.InitConfig(configPath)
	em := NewQQEmailService(config.App.SMTP.Key, config.App.SMTP.Sender,
		config.App.SMTP.SMTPHost, config.App.SMTP.SMTPPort)
	err := em.Send(context.Background(), testEmails, "发送邮件啦", time.Now().Format(time.DateTime)+": 这是一封测试邮件，")
	if err != nil {
		t.Error(err)
	}
}
