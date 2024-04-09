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

func GetApiList(c *gin.Context) {
	sc := c.MustGet(common.GIN_CTX_CONFIG_CONFIG).(*config.ServerConfig)

	apis, err := models.GetApiAll()
	if err != nil {
		sc.Logger.Ctx(c.Request.Context()).Error("去数据库中拿所有的api接口错误",
			zap.Error(err),
		)
		common.FailWithMessage(fmt.Sprintf("去数据库中拿所有的api接口错误:%v", err.Error()), c)
		return
	}

	fatherApiMap := make(map[uint]*models.Api)
	for _, api := range apis {
		api.Key = api.ID
		api.Value = api.ID
		if api.Pid == 0 {
			// 说明这个菜单是父级
			fatherApiMap[api.ID] = api
			continue
		}

		// 说明menu是子集
		fatherApi, err := models.GetApiById(api.Pid)
		if err != nil {
			sc.Logger.Error("通过Pid找Api错误", zap.Error(err))
			continue
		}
		fatherApi.Key = fatherApi.ID
		fatherApi.Value = fatherApi.ID

		load, ok := fatherApiMap[fatherApi.ID]

		if !ok {
			//之前还没设置过 这个父级
			fatherApi.Children = make([]*models.Api, 0)
			fatherApi.Children = append(fatherApi.Children, api)
			fatherApiMap[fatherApi.ID] = fatherApi
		} else {
			// 存在的话 我们直接把menu塞入 Children
			load.Children = append(load.Children, api)
		}

	}

	finalApis := make([]*models.Api, 0)
	// 最终遍历	 fatherMenuMap
	for _, m := range fatherApiMap {
		m := m
		finalApis = append(finalApis, m)
	}
	common.OkWithDetailed(finalApis, "ok", c)
}

func DeleteApi(c *gin.Context) {
	ctx := c.Request.Context()
	sc := c.MustGet(common.GIN_CTX_CONFIG_CONFIG).(*config.ServerConfig)
	// 校验一下 menu字段
	id := c.Param("id")
	sc.Logger.Ctx(ctx).Info("删除api权限", zap.Any("id", id))

	intVar, _ := strconv.Atoi(id)
	dbApi, err := models.GetApiById(intVar)
	if err != nil {
		sc.Logger.Ctx(ctx).Error("根据id找api错误", zap.Any("api", id), zap.Error(err))
		common.FailWithMessage(err.Error(), c)
		return
	}

	err = dbApi.DeleteOne()
	if err != nil {
		sc.Logger.Ctx(ctx).Error("根据id删除api错误", zap.Any("api", id), zap.Error(err))
		common.FailWithMessage(err.Error(), c)
		return
	}
	common.OkWithMessage("删除成功", c)
}

func CreateApi(c *gin.Context) {
	sc := c.MustGet(common.GIN_CTX_CONFIG_CONFIG).(*config.ServerConfig)
	ctx := c.Request.Context()
	var reqApi models.Api
	err := c.ShouldBindJSON(&reqApi)
	if err != nil {
		sc.Logger.Ctx(ctx).Error("解析新增reqApi请求失败", zap.Any("reqApi", reqApi), zap.Error(err))
		common.FailWithMessage(err.Error(), c)
		return
	}

	// 在这里校验字段，是否必填，范围是否正确
	err = validate.Struct(reqApi)
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

	err = reqApi.CreateOne()
	if err != nil {
		sc.Logger.Ctx(ctx).Error("创建Api错误", zap.Any("reqApi", reqApi), zap.Error(err))
		common.FailWithMessage(err.Error(), c)
		return
	}

	common.OkWithMessage("创建成功", c)
}
func UpdateApi(c *gin.Context) {
	sc := c.MustGet(common.GIN_CTX_CONFIG_CONFIG).(*config.ServerConfig)
	ctx := c.Request.Context()
	var reqApi models.Api
	err := c.ShouldBindJSON(&reqApi)
	if err != nil {
		sc.Logger.Ctx(ctx).Error("解析更新reqApi请求失败", zap.Any("reqApi", reqApi), zap.Error(err))
		common.FailWithMessage(err.Error(), c)
		return
	}
	err = validate.Struct(reqApi)
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

	_, err = models.GetApiById(int(reqApi.ID))
	if err != nil {
		sc.Logger.Ctx(ctx).Error("根据id找Api错误", zap.Any("Api", reqApi), zap.Error(err))
		common.FailWithMessage(err.Error(), c)
		return
	}
	err = reqApi.UpdateOne()
	if err != nil {
		sc.Logger.Ctx(ctx).Error("根据id更新Api错误", zap.Any("Api", reqApi), zap.Error(err))
		common.FailWithMessage(err.Error(), c)
		return
	}

	common.OkWithMessage("更新成功", c)
}
