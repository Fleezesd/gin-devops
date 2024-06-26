package web

import (
	"net/http"
	"time"

	"github.com/fleezesd/gin-devops/src/common"
	"github.com/fleezesd/gin-devops/src/config"
	"github.com/fleezesd/gin-devops/src/web/middleware"
	"github.com/fleezesd/gin-devops/src/web/view"
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/requestid"
	"github.com/gin-gonic/gin"
	ginprometheus "github.com/zsais/go-gin-prometheus"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
)

// StartHttp 单独启动Gin
func StartHttp(sc *config.ServerConfig) error {
	// 配置模式
	gin.SetMode(gin.DebugMode)
	gin.DisableConsoleColor()

	r := gin.New()

	// 记录耗时 传递变量中间件
	m := make(map[string]interface{})
	m[common.GIN_CTX_CONFIG_CONFIG] = sc
	r.Use(middleware.ConfigMiddleware(m))

	// oTel 中间件
	r.Use(otelgin.Middleware("gin-devops"))
	// zap 中间件
	r.Use(middleware.NewGinZapLogger(sc.Logger))
	// requestId 中间件
	r.Use(requestid.New())
	// prometheus 中间件
	p := ginprometheus.NewPrometheus("gin_devops")
	p.Use(r)
	// cors 跨域中间件
	r.Use(cors.Default())

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
