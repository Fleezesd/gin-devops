package main

import (
	"flag"

	"github.com/fleezesd/gin-devops/src/common"
	"github.com/fleezesd/gin-devops/src/config"
	"github.com/fleezesd/gin-devops/src/web"
	"go.uber.org/zap"
)

func main() {

	var (
		configFile string
	)
	flag.StringVar(&configFile, "configFile", "./server.yml", "Configuration file path.")
	flag.Parse()

	// 读取配置
	sc, err := config.LoadServer(configFile)
	if err != nil {
		panic(err)
	}

	// 日志配置
	logger := common.NewZapLogger(sc.LogLevel, sc.LogFilePath)
	logger.Info("读取Server配置",
		zap.String("httpAddr", sc.HttpAddr),
	)

	// 启动Http Gin
	err = web.StartHttp(sc)
	if err != nil {
		panic(err)
	}
}
