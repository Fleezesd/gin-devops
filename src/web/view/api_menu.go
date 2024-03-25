package view

import (
	"strconv"

	"github.com/fleezesd/gin-devops/src/common"
	"github.com/fleezesd/gin-devops/src/config"
	"github.com/fleezesd/gin-devops/src/models"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func GetMenuList(c *gin.Context) {
	// 拿到用户的role列表 遍历role列表 拿到menuList
	var (
		userName = c.MustGet(common.GIN_CTX_JWT_USER_NAME).(string)
		ctx      = c.Request.Context()
	)
	sc := c.MustGet(common.GIN_CTX_CONFIG_CONFIG).(*config.ServerConfig)
	dbUser, err := models.GetUserByUserName(userName)
	if err != nil {
		sc.Logger.Ctx(ctx).Error("通过token解析的用户名,获取用户失败! 用户名不存在!",
			zap.Error(err),
		)
		common.FailWithMessage(err.Error(), c)
		return
	}

	var (
		fatherMenuMap  = make(map[uint]*models.Menu)
		uniqueChildMap = make(map[uint]*models.Menu)
	)

	roles := dbUser.Roles
	for _, role := range roles {
		sc.Logger.Ctx(ctx).Info("role的menuList详情",
			zap.Any("menu", role.Menus),
		)
		// 去重 menu
		for _, menu := range role.Menus {
			menu.Meta = &models.MenuMeta{}
			menu.Meta.Icon = menu.Icon
			menu.Meta.Title = menu.Title
			menu.Key = menu.ID
			menu.Value = menu.ID
			if menu.ParentMenu == "" {
				fatherMenuMap[menu.ID] = menu
			}

			// 判断是否重复子菜单
			_, ok := uniqueChildMap[menu.ID]
			if ok {
				continue
			}
			// 否则塞入
			uniqueChildMap[menu.ID] = menu

			// 子菜单录入父菜单
			fatherMenuId, _ := strconv.Atoi(menu.ParentMenu)
			fatherMenu, err := models.GetMenuById(fatherMenuId)
			if err != nil {
				sc.Logger.Ctx(ctx).Error("menu寻找错误", zap.Error(err))
				continue
			}
			load, ok := fatherMenuMap[fatherMenu.ID]
			if !ok {
				fatherMenuMap[fatherMenu.ID] = fatherMenu
				fatherMenu.Children = make([]*models.Menu, 0)
				fatherMenu.Children = append(fatherMenu.Children, menu)
			} else {
				load.Children = append(load.Children, menu)
			}
		}
	}

	finalMenus := make([]*models.Menu, 0)

	for _, menu := range fatherMenuMap {
		finalMenus = append(finalMenus, menu)
	}

	common.OkWithDetailed(finalMenus, "获取菜单列表成功", c)
}
