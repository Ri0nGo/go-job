package qq

import (
	"context"
	"github.com/stretchr/testify/assert"
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
		"920728039@qq.com",
	}
	configPath := getMasterConfigPath()
	config.InitConfig(configPath)

	smtpCfg, ok := config.App.SMTP["qq"]
	assert.Equal(t, ok, true)
	svc := NewQQEmailService("testqq", smtpCfg.Key,
		smtpCfg.Sender,
		smtpCfg.SMTPHost,
		smtpCfg.SMTPPort)
	err := svc.Send(context.Background(), testEmails, "发送邮件啦", time.Now().Format(time.DateTime)+": 这是一封测试邮件，")
	if err != nil {
		t.Error(err)
	}
}
