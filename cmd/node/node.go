package main

import (
	"flag"
	"fmt"
	"go-job/node/pkg/auth"
	"go-job/node/pkg/config"
	"go-job/node/pkg/ioc"
	"go-job/node/pkg/startup"
	"log/slog"
)

func main() {
	var configPath string
	flag.StringVar(&configPath, "c", "./config/node.yaml", "set yaml file path")
	flag.Parse()

	// 初始化配置
	err := config.InitConfig(configPath)
	if err != nil {
		panic(err)
	}

	RunApp()
}

func beforeRunWeb(container *ioc.WebContainer) {
	auth.InitJwtToken(config.App.Master.Key)
	if err := startup.SyncJobFromMaster(container.JobSvc); err != nil {
		slog.Error("sync job from master error", "err", err)
	}
}

func RunApp() {
	container := ioc.InitWebServer()
	beforeRunWeb(container)
	container.Engine.Run(fmt.Sprintf("%s:%d", config.App.Server.Ip, config.App.Server.Port))
}
