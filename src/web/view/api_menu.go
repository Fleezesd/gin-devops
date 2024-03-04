package view

import (
	"github.com/fleezesd/gin-devops/src/common"
	"github.com/fleezesd/gin-devops/src/config"
	"github.com/fleezesd/gin-devops/src/models"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func GetMenuList(c *gin.Context) {
	// 拿到用户的role列表 遍历role列表 拿到menuList
	userName := c.MustGet(common.GIN_CTX_JWT_USER_NAME).(string)
	sc := c.MustGet(common.GIN_CTX_CONFIG_CONFIG).(*config.ServerConfig)
	dbUser, err := models.GetUserByUserName(userName)
	if err != nil {
		sc.Logger.Error("通过token解析的用户名,获取用户失败! 用户名不存在!",
			zap.Error(err),
		)
		common.FailWithMessage(err.Error(), c)
		return
	}
	menuList := []*models.Menu{}
	roles := dbUser.Roles
	for _, role := range roles {
		menuList = append(menuList, role.Menus...)
	}
	common.OkWithDetailed(menuList, "获取菜单列表成功", c)
}
