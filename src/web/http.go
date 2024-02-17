package web

import (
	"net/http"
	"time"

	"github.com/fleezesd/gin-devops/src/common"
	"github.com/fleezesd/gin-devops/src/config"
	"github.com/fleezesd/gin-devops/src/web/middleware"
	"github.com/fleezesd/gin-devops/src/web/view"
	"github.com/gin-contrib/requestid"
	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
)

// StartHttp 单独启动Gin
func StartHttp(sc *config.ServerConfig) error {
	// 配置模式
	gin.SetMode(gin.ReleaseMode)
	gin.DisableConsoleColor()

	r := gin.New()

	// 记录耗时 传递变量中间件
	m := make(map[string]interface{})
	m[common.GIN_CTX_CONFIG_CONFIG] = sc
	r.Use(middleware.ConfigMiddleware(m))
	// zap 中间件
	r.Use(middleware.NewGinZapLogger(sc.Logger))
	r.Use(ginzap.RecoveryWithZap(sc.Logger, true))
	// requestId 中间件
	r.Use(requestid.New())

	// 注册路由
	view.ConfigRoutes(r)
	// Http 读写超时配置
	s := http.Server{
		Addr:           sc.HttpAddr,
		Handler:        r,
		ReadTimeout:    time.Second * 5,
		WriteTimeout:   time.Second * 5,
		MaxHeaderBytes: 1 << 20,
	}
	return s.ListenAndServe()
}
