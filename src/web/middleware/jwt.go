package middleware

import (
	"fmt"
	"strings"
	"time"

	"github.com/fleezesd/gin-devops/src/common"
	"github.com/fleezesd/gin-devops/src/config"
	"github.com/fleezesd/gin-devops/src/models"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func JWTAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1. 遍历header 寻找 Authorization
		authHeaderString := c.Request.Header.Get("Authorization")
		if authHeaderString == "" {
			common.Req401WithWithDetailed(gin.H{"reload": true}, "未登录或非法访问没有Authorization的header", c)
			c.Abort()
			return
		}
		// authHeaderString格式校验： Bearer {{ token }}
		parts := strings.SplitN(authHeaderString, " ", 2)
		if !(len(parts) == 2 && parts[0] == "Bearer") {
			common.Req401WithWithDetailed(gin.H{"reload": true}, "请求头中的auth格式错误", c)
			c.Abort()
			return
		}

		// 2. 拿到token 做jwt解析
		sc := c.MustGet(common.GIN_CTX_CONFIG_CONFIG).(*config.ServerConfig)
		userClaims, err := models.ParseToken(parts[1], sc)
		if err != nil {
			common.Req401WithWithDetailed(gin.H{"reload": true}, fmt.Sprintf("parseToken 解析token包含的信息错误:%v", err.Error()), c)
			c.Abort()
			return
		}

		// 3. 验证token 是否过期
		// 临期逻辑 建立新token 续期
		if userClaims.RegisteredClaims.ExpiresAt.Unix()-time.Now().Unix() < int64(sc.JWT.BufferDuration/time.Second) {
			sc.Logger.Info("jwt 临期了，刷新jwt",
				zap.String("user", userClaims.Username),
				zap.Any("老token过期时间", userClaims.RegisteredClaims.ExpiresAt),
				zap.Any("临期窗口", sc.JWT.BufferDuration),
			)
			// 新生成 Token
			newToken, err := models.GenJwtToken(userClaims.User, sc)
			if err != nil {
				common.Result500(7, gin.H{}, fmt.Sprintf("parseToken 解析token包含的信息错误:%v", err.Error()), c)
				c.Abort()
			}
			c.Header("new-token", newToken)	// 前端配合 下次请求带新token
		}
		// 4. claim 对象传递下去
		c.Set(common.GIN_CTX_JWT_USER, userClaims.User)
		c.Next()
	}
}
