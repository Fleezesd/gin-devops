package main

import (
	"context"
	"flag"

	"github.com/fleezesd/gin-devops/src/common"
	"github.com/fleezesd/gin-devops/src/config"
	"github.com/fleezesd/gin-devops/src/models"
	"github.com/fleezesd/gin-devops/src/web"
	"github.com/prometheus/common/version"
	"github.com/spf13/cobra"
	"github.com/uptrace/opentelemetry-go-extra/otelzap"
	"github.com/uptrace/uptrace-go/uptrace"
	"go.uber.org/zap"
)

var sc *config.ServerConfig

var rootCmd = &cobra.Command{
	Use:               "server",
	Short:             "gin devops server",
	SilenceUsage:      true,
	DisableAutoGenTag: true,
	Long:              `devops restful api server`,
	PreRunE:           PreRunE,
	RunE:              RunE,
	Version:           version.Info(),
}

func PreRunE(cmd *cobra.Command, args []string) (err error) {
	var (
		configFile string
		ctx        = context.Background()
	)
	flag.StringVar(&configFile, "configFile", "./server.yml", "Configuration file path.")
	flag.Parse()

	// 读取配置
	sc, err = config.LoadServer(configFile)
	if err != nil {
		panic(err)
	}

	// 日志配置 + oTelZap
	logger := otelzap.New(common.NewZapLogger(sc.LogLevel, sc.LogFilePath))

	sc.Logger = logger
	logger.Ctx(ctx).Info("读取Server配置",
		zap.String("httpAddr", sc.HttpAddr),
		zap.String("logLevel", sc.LogLevel),
		zap.String("logFilePath", sc.LogFilePath),
	)

	// 初始化数据库
	if err = models.InitDB(sc.Mysql.DSN); err != nil {
		logger.Error("初始化gorm db错误",
			zap.Error(err),
		)
		return
	}
	logger.Ctx(ctx).Info("初始化gorm db连接成功")

	// 同步表结构
	if err = models.MigrateTable(); err != nil {
		logger.Ctx(ctx).Error("gorm db同步表结构错误",
			zap.Error(err),
		)
		return
	}
	logger.Ctx(ctx).Info("gorm db同步表结构成功")

	return
}

func RunE(cmd *cobra.Command, args []string) error {
	// 测试注册和密码加盐 后续删除
	models.MockUserRegister(sc)

	ctx := context.Background()
	uptrace.ConfigureOpentelemetry(
		uptrace.WithDSN("http://gin_devops_secret_token@localhost:14318?grpc=14317"),
		uptrace.WithServiceName("gin-devops"),
		uptrace.WithServiceVersion("1.0.0"),
	)
	defer uptrace.Shutdown(ctx)

	// 启动Http Gin
	return web.StartHttp(sc)
}
