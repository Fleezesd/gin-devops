package view

import (
	"fmt"

	"github.com/fleezesd/gin-devops/src/common"
	"github.com/fleezesd/gin-devops/src/config"
	"github.com/fleezesd/gin-devops/src/models"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
)

func UserLogin(c *gin.Context) {
	var user models.UserLoginRequest
	sc := c.MustGet(common.GIN_CTX_CONFIG_CONFIG).(*config.ServerConfig)
	// 校验json字段
	if err := c.ShouldBindJSON(&user); err != nil {
		common.FailWithMessage(err.Error(), c)
		return
	}
	// 校验validate字段
	if err := validate.Struct(user); err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			common.FailWithWithDetailed(gin.H{
				"翻译前": err.Error(),
				"翻译后": err.Translate(trans),
			}, "请求校验出错", c)
			return
		}
	}

	// 检测用户
	dbUser, err := models.CheckUserPassword(&user)
	if err != nil {
		sc.Logger.Error("登陆失败! 用户名不存在或者密码错误!",
			zap.Error(err),
		)
		common.FailWithMessage(fmt.Sprintf("用户名不存在或者密码错误:%v", err.Error()), c)
		return
	}
	// 生成jwt
	models.TokenNext(dbUser, c)
}

// GetUserInfoAfterLogin 登录后获取用户信息 来自于 jwt Header
func GetUserInfoAfterLogin(c *gin.Context) {
	// 拿到 UserClaim
	user := c.MustGet(common.GIN_CTX_JWT_USER).(*models.User)
	common.OkWithDetailed(user, "ok", c)
}
