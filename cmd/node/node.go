package main

import (
	"flag"
	"fmt"
	"github.com/gin-gonic/gin"
	"go-job/node/pkg/config"
	"go-job/node/router"
)

func main() {
	var configPath string
	flag.StringVar(&configPath, "c", "./config/node.yaml", "set yaml file path")

	// 初始化配置
	err := config.InitConfig(configPath)
	if err != nil {
		panic(err)
	}

	runWeb()
}

func runWeb() {
	engine := gin.Default()
	router.RegistryRoute(engine)
	engine.Run(fmt.Sprintf("%s:%d", config.App.Config.Ip, config.App.Config.Port))
}
