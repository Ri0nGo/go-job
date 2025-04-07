package main

import (
	"flag"
	"fmt"
	"go-job/master/pkg/config"
	"go-job/master/pkg/ioc"
)

func main() {
	var configPath string
	flag.StringVar(&configPath, "c", "./config/master.yaml", "set yaml file path")

	// 初始化配置
	err := config.InitConfig(configPath)
	if err != nil {
		panic(err)
	}

	server := ioc.InitWebServer()
	server.Run(fmt.Sprintf("%s:%d",
		config.App.Server.Ip, config.App.Server.Port))

}
