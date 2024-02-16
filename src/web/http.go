package web

import (
	"net/http"
	"time"

	"github.com/fleezesd/gin-devops/src/common"
	"github.com/fleezesd/gin-devops/src/config"
	"github.com/fleezesd/gin-devops/src/web/middleware"
	"github.com/fleezesd/gin-devops/src/web/view"
	"github.com/gin-gonic/gin"
)

// StartHttp 单独启动Gin
func StartHttp(sc *config.ServerConfig) error {
	// 配置模式
	gin.SetMode(gin.ReleaseMode)
	gin.DisableConsoleColor()

	r := gin.Default()
	// 记录耗时 传递变量中间件
	m := make(map[string]interface{})
	m[common.GIN_CTX_CONFIG_CONFIG] = sc
	r.Use(middleware.TimeCost(), middleware.ConfigMiddleware(m))
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
