package main

import (
	"flag"
	"fmt"

	"github.com/fleezesd/gin-devops/src/config"
	"github.com/fleezesd/gin-devops/src/web"
)

func main() {

	var (
		configFile string
	)
	flag.StringVar(&configFile, "configFile", "./server.yml", "Configuration file path.")
	flag.Parse()

	// 读取配置
	serverConfig, err := config.LoadServer(configFile)
	if err != nil {
		panic(err)
	}
	fmt.Printf("server配置:%v\n", serverConfig)

	// 启动Http Gin
	err = web.StartHttp(serverConfig)
	if err != nil {
		panic(err)
	}
}
