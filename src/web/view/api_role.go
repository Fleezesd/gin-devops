package view

import (
	"fmt"
	"strconv"

	"github.com/fleezesd/gin-devops/src/common"
	"github.com/fleezesd/gin-devops/src/config"
	"github.com/fleezesd/gin-devops/src/models"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
)

func GetRoleListAll(c *gin.Context) {

	sc := c.MustGet(common.GIN_CTX_CONFIG_CONFIG).(*config.ServerConfig)

	// 	数据库中拿到所有的menu列表

	roles, err := models.GetRoleAll()
	if err != nil {
		sc.Logger.Ctx(c.Request.Context()).Error("去数据库中拿所有的角色错误",
			zap.Error(err),
		)
		common.FailWithMessage(fmt.Sprintf("去数据库中拿所有的角色错误:%v", err.Error()), c)
		return
	}
	fmt.Println(roles)
	// 遍历 role 准备menuIds 列表
	for _, role := range roles {
		for _, menu := range role.Menus {
			menu.Key = menu.ID
			menu.Value = menu.ID
		}
		for _, api := range role.Apis {
			api.Key = api.ID
			api.Value = api.ID
		}

	}

	common.OkWithDetailed(roles, "ok", c)
}

func CreateRole(c *gin.Context) {

	sc := c.MustGet(common.GIN_CTX_CONFIG_CONFIG).(*config.ServerConfig)
	ctx := c.Request.Context()
	// 校验一下 menu字段
	var reqRole models.Role
	err := c.ShouldBindJSON(&reqRole)
	if err != nil {
		sc.Logger.Ctx(ctx).Error("解析新增角色请求失败", zap.Any("角色", reqRole), zap.Error(err))
		common.FailWithMessage(err.Error(), c)
		return
	}

	// 在这里校验字段，是否必填，范围是否正确
	err = validate.Struct(reqRole)
	if err != nil {
		if errors, ok := err.(validator.ValidationErrors); ok {
			common.FailWithWithDetailed(

				gin.H{
					"翻译前": err.Error(),
					"翻译后": errors.Translate(trans),
				},
				"请求出错",
				c,
			)
			return
		}
		common.FailWithMessage(err.Error(), c)
		return

	}

	menus := make([]*models.Menu, 0)
	// 遍历角色menu 列表 找到角色
	for _, menuId := range reqRole.MenuIds {
		dbMenu, err := models.GetMenuById(menuId)
		if err != nil {
			sc.Logger.Ctx(ctx).Error("根据id找菜单错误", zap.Any("菜单", reqRole), zap.Error(err))
			common.FailWithMessage(err.Error(), c)
			return
		}
		menus = append(menus, dbMenu)

	}
	reqRole.Menus = menus
	err = reqRole.CreateOne()
	if err != nil {
		sc.Logger.Ctx(ctx).Error("创建角色错误", zap.Any("角色", reqRole), zap.Error(err))
		common.FailWithMessage(err.Error(), c)
		return
	}

	common.OkWithMessage("创建成功", c)

}

func UpdateRole(c *gin.Context) {
	sc := c.MustGet(common.GIN_CTX_CONFIG_CONFIG).(*config.ServerConfig)
	ctx := c.Request.Context()
	// 校验一下 menu字段
	var reqRole models.Role
	err := c.ShouldBindJSON(&reqRole)
	if err != nil {
		sc.Logger.Ctx(ctx).Error("解析更新角色请求失败", zap.Any("角色", reqRole), zap.Error(err))
		common.FailWithMessage(err.Error(), c)
		return
	}

	err = validate.Struct(reqRole)
	if err != nil {
		if errors, ok := err.(validator.ValidationErrors); ok {
			common.FailWithWithDetailed(

				gin.H{
					"翻译前": err.Error(),
					"翻译后": errors.Translate(trans),
				},
				"请求出错",
				c,
			)
			return
		}
		common.FailWithMessage(err.Error(), c)
		return

	}

	_, err = models.GetRoleById(int(reqRole.ID))
	if err != nil {
		sc.Logger.Ctx(ctx).Error("根据id找角色错误", zap.Any("角色", reqRole), zap.Error(err))
		common.FailWithMessage(err.Error(), c)
		return
	}

	menus := make([]*models.Menu, 0)
	// 遍历角色menu 列表 找到角色
	for _, menuId := range reqRole.MenuIds {
		dbMenu, err := models.GetMenuById(menuId)
		if err != nil {
			sc.Logger.Error("根据id找菜单错误", zap.Any("菜单", reqRole), zap.Error(err))
			common.FailWithMessage(err.Error(), c)
			return
		}
		menus = append(menus, dbMenu)

	}

	apis := make([]*models.Api, 0)
	// 遍历角色menu 列表 找到角色
	for _, apiId := range reqRole.ApiIds {
		dbApi, err := models.GetApiById(apiId)
		if err != nil {
			sc.Logger.Ctx(ctx).Error("根据id找api错误", zap.Any("api", reqRole), zap.Error(err))
			common.FailWithMessage(err.Error(), c)
			return
		}
		apis = append(apis, dbApi)

	}

	err = reqRole.UpdateMenus(menus)
	if err != nil {
		sc.Logger.Ctx(ctx).Error("更新角色和关联的菜单错误", zap.Any("角色", reqRole), zap.Error(err))
		common.FailWithMessage(err.Error(), c)
		return
	}

	err = reqRole.UpdateApis(apis, sc)
	if err != nil {
		sc.Logger.Ctx(ctx).Error("更新角色和关联的api错误", zap.Any("角色", reqRole), zap.Error(err))
		common.FailWithMessage(err.Error(), c)
		return
	}

	common.OkWithMessage("更新成功", c)
	sc.Logger.Info("更新角色和关联的菜单成功", zap.Any("角色", reqRole))
}

func SetRoleStatus(c *gin.Context) {

	sc := c.MustGet(common.GIN_CTX_CONFIG_CONFIG).(*config.ServerConfig)
	ctx := c.Request.Context()
	// 校验一下 menu字段
	var reqRole models.SetRoleStatusReq
	err := c.ShouldBindJSON(&reqRole)
	if err != nil {
		sc.Logger.Ctx(ctx).Error("解析更新角色请求失败", zap.Any("角色", reqRole), zap.Error(err))
		common.FailWithMessage(err.Error(), c)
		return
	}

	// 在这里校验字段，是否必填，范围是否正确
	err = validate.Struct(reqRole)
	if err != nil {

		// 这里为什么要判断错误是否是 ValidationErrors
		if errors, ok := err.(validator.ValidationErrors); ok {
			common.FailWithWithDetailed(

				gin.H{
					"翻译前": err.Error(),
					"翻译后": errors.Translate(trans),
				},
				"请求出错",
				c,
			)
			return
		}
		common.FailWithMessage(err.Error(), c)
		return

	}

	dbRole, err := models.GetRoleById(reqRole.Id)
	if err != nil {
		sc.Logger.Ctx(ctx).Error("根据id找角色错误", zap.Any("角色", reqRole), zap.Error(err))
		common.FailWithMessage(err.Error(), c)
		return
	}
	dbRole.Status = reqRole.Status
	err = dbRole.UpdateMenus(dbRole.Menus)
	if err != nil {
		sc.Logger.Ctx(ctx).Error("更新角色和关联的菜单错误", zap.Any("角色", reqRole), zap.Error(err))
		common.FailWithMessage(err.Error(), c)
		return
	}

	common.OkWithMessage("更新成功", c)
	sc.Logger.Ctx(ctx).Info("更新角色和关联的菜单成功", zap.Any("角色", reqRole))

}

func DeleteRole(c *gin.Context) {
	sc := c.MustGet(common.GIN_CTX_CONFIG_CONFIG).(*config.ServerConfig)
	ctx := c.Request.Context()
	// 校验一下 menu字段
	id := c.Param("id")
	sc.Logger.Info("删除角色", zap.Any("id", id))

	// db中根据id找到这个user
	intVar, _ := strconv.Atoi(id)
	dbRole, err := models.GetRoleById(intVar)
	if err != nil {
		sc.Logger.Ctx(ctx).Error("根据id找角色错误", zap.Any("角色", id), zap.Error(err))
		common.FailWithMessage(err.Error(), c)
		return
	}

	err = dbRole.DeleteOne()
	if err != nil {
		sc.Logger.Ctx(ctx).Error("根据id删除角色错误", zap.Any("角色", id), zap.Error(err))
		common.FailWithMessage(err.Error(), c)
		return
	}
	common.OkWithMessage("删除成功", c)
}
