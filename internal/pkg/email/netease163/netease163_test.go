package netease163

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

func TestNetEase163EmailService(t *testing.T) {
	testEmails := []string{
		"920728039@qq.com",
	}
	configPath := getMasterConfigPath()
	config.InitConfig(configPath)
	smtpCfg, ok := config.App.SMTP["netease163"]
	assert.Equal(t, ok, true)
	svc := NewNetEase163EmailService("test163", smtpCfg.Key,
		smtpCfg.Sender,
		smtpCfg.SMTPHost,
		smtpCfg.SMTPPort)
	err := svc.Send(context.Background(), testEmails, "163发送邮件啦", time.Now().Format(time.DateTime)+": 这是一封测试邮件，")
	if err != nil {
		t.Error(err)
	}
}
