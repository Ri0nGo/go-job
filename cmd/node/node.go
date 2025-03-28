package main

import (
	"flag"
	"go-job/node/pkg/config"
)

func main() {
	var configPath string
	flag.StringVar(&configPath, "c", "./config/node.yaml", "set yaml file path")

	// 初始化配置
	err := config.InitConfig(configPath)
	if err != nil {
		panic(err)
	}

}
