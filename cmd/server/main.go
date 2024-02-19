package main

import (
	"flag"

	"github.com/fleezesd/gin-devops/src/common"
	"github.com/fleezesd/gin-devops/src/config"
	"github.com/fleezesd/gin-devops/src/models"
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
	sc.Logger = logger
	logger.Info("读取Server配置",
		zap.String("httpAddr", sc.HttpAddr),
		zap.String("logLevel", sc.LogLevel),
		zap.String("logFilePath", sc.LogFilePath),
	)

	// 初始化数据库
	if err := models.InitDB(sc.Mysql.DSN); err != nil {
		logger.Error("初始化gorm db错误",
			zap.Error(err),
		)
		return
	}
	logger.Info("初始化gorm db连接成功")

	// 启动Http Gin
	err = web.StartHttp(sc)
	if err != nil {
		panic(err)
	}
}
