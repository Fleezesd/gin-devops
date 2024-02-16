package middleware

import (
	"time"

	"github.com/fleezesd/gin-devops/src/common"
	"github.com/fleezesd/gin-devops/src/config"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func TimeCost() gin.HandlerFunc {
	return func(c *gin.Context) {
		//请求前获取当前时间
		nowTime := time.Now()

		//请求处理
		c.Next()

		sc := c.MustGet(common.GIN_CTX_CONFIG_CONFIG).(*config.ServerConfig)
		sc.Logger.Info("耗时中间件打印结果",
			zap.String("url", c.Request.URL.String()),
			zap.Duration("耗时", time.Since(nowTime)),
		)
	}
}
