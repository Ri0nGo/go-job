package main

import (
	"flag"
	"fmt"
	"go-job/master/pkg/config"
	"go-job/master/pkg/ioc"
	"go-job/master/pkg/job"
	"log/slog"
)

func main() {
	var configPath string
	flag.StringVar(&configPath, "c", "./config/master.yaml", "set yaml file path")
	flag.Parse()

	// 初始化配置
	err := config.InitConfig(configPath)
	if err != nil {
		panic(err)
	}

	RunApp()
}

func RunApp() {
	container := ioc.InitWebServer()
	bootstrap(container)
	container.Engine.Run(fmt.Sprintf("%s:%d",
		config.App.Server.Ip, config.App.Server.Port))
}

func bootstrap(c *ioc.WebContainer) {
	err := job.InitGlobalData(c.MysqlDB, c.JobSvc, c.NotifyStore)
	if err != nil {
		slog.Error("init job data to node error", "err", err)
	}

}
