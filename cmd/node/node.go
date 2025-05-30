package main

import (
	"flag"
	"fmt"
	"github.com/gin-gonic/gin"
	"go-job/node/pkg/auth"
	"go-job/node/pkg/config"
	"go-job/node/router"
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

	beforeRunWeb()
	runWeb()
}

func beforeRunWeb() {
	auth.InitJwtToken(config.App.Master.Key)
}

func runWeb() {
	engine := gin.Default()
	router.InitRouter(engine)
	engine.Run(fmt.Sprintf("%s:%d", config.App.Server.Ip, config.App.Server.Port))
}
