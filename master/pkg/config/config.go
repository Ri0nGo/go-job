package config

import (
	"github.com/spf13/viper"
	"log"
	"log/slog"
)

func InitConfig(filePath string) error {
	setting, err := loadConfig(filePath)
	if err != nil {
		return err
	}
	slog.Info("config info", "detail", setting)
	App = setting
	return nil
}

func loadConfig(filePath string) (*Application, error) {
	viper.SetConfigFile(filePath)
	err := viper.ReadInConfig()
	if err != nil {
		return nil, err
	}
	var app Application
	err = viper.Unmarshal(&app)
	if err != nil {
		log.Fatalf("Error unmarshaling config, %s", err)
	}
	return &app, nil
}
