package middleware

import (
	"github.com/fleezesd/gin-devops/src/common"
	"github.com/fleezesd/gin-devops/src/config"
	"github.com/fleezesd/gin-devops/src/models"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func CasbinAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 取得用户的角色
		sc := c.MustGet(common.GIN_CTX_CONFIG_CONFIG).(*config.ServerConfig)
		userName := c.MustGet(common.GIN_CTX_JWT_USER_NAME).(string)
		dbUser, err := models.GetUserByUserName(userName)
		if err != nil {
			sc.Logger.Ctx(c.Request.Context()).Error("[Casbin]获取用户失败! 用户名不存在!",
				zap.Error(err),
			)
			common.FailWithMessage(err.Error(), c)
			c.Abort()
		}
		// 拿到用户的role和menu 进行auth验证
		roles := dbUser.Roles
		pass := false
		for _, role := range roles {
			ok, err := models.Enforcer.Enforce(role.RoleValue, c.Request.URL.Path, c.Request.Method)
			if err != nil {
				sc.Logger.Ctx(c.Request.Context()).Error("[Casbin]验证用户权限失败!",
					zap.Error(err),
				)
				common.FailWithMessage(err.Error(), c)
				c.Abort()
			}
			if ok {
				pass = true
				break
			}
		}
		if !pass {
			sc.Logger.Ctx(c.Request.Context()).Error("[Casbin]用户没有权限!",
				zap.String("username", userName),
				zap.String("path", c.Request.URL.Path),
				zap.String("method", c.Request.Method),
			)
			common.Req403WithWithMessage("用户没有权限", c)
			c.Abort()
		}
		sc.Logger.Ctx(c.Request.Context()).Info("[Casbin]用户校验通过!")
		c.Next()
	}
}
