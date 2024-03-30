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

type DefineUserOrGroup struct {
	Label string `json:"label"`
	Value string `json:"value"`
}

func UserLogin(c *gin.Context) {
	var (
		user models.UserLoginRequest
		ctx  = c.Request.Context()
		sc   = c.MustGet(common.GIN_CTX_CONFIG_CONFIG).(*config.ServerConfig)
	)
	if err := c.ShouldBindJSON(&user); err != nil {
		sc.Logger.Ctx(ctx).Error("登陆失败! 请求参数错误!",
			zap.Error(err),
		)
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
	dbUser, err := models.CheckUserPassword(sc.Logger, ctx, &user)
	if err != nil {
		sc.Logger.Ctx(ctx).Error("登陆失败! 用户名不存在或者密码错误!",
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
	var (
		sc  = c.MustGet(common.GIN_CTX_CONFIG_CONFIG).(*config.ServerConfig)
		ctx = c.Request.Context()
	)
	// 拿到 UserClaim
	userName := c.MustGet(common.GIN_CTX_JWT_USER_NAME).(string)
	dbUser, err := models.GetUserByUserName(userName)
	if err != nil {
		sc.Logger.Ctx(ctx).Error("获取用户失败! 用户名不存在!",
			zap.Error(err),
		)
		common.FailWithMessage(err.Error(), c)
		return
	}
	common.OkWithDetailed(dbUser, "ok", c)
}

func GetPermCode(c *gin.Context) {
	common.OkWithDetailed([]string{"2000", "4000", "6000"}, "ok", c)
}

func CreateAccount(c *gin.Context) {
	var (
		sc  = c.MustGet(common.GIN_CTX_CONFIG_CONFIG).(*config.ServerConfig)
		ctx = c.Request.Context()
	)
	// 校验一下 menu字段
	var reqUser models.User
	err := c.ShouldBindJSON(&reqUser)
	if err != nil {
		sc.Logger.Ctx(ctx).Error("解析新增用户请求失败", zap.Any("user", reqUser), zap.Error(err))
		common.FailWithMessage(err.Error(), c)
		return
	}

	// 在这里校验字段，是否必填，范围是否正确
	if err := validate.Struct(reqUser); err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			common.FailWithWithDetailed(gin.H{
				"翻译前": err.Error(),
				"翻译后": err.Translate(trans),
			}, "请求校验出错", c)
			return
		}
	}

	// 根据 rolesFront 去db中查询 role ，把role给他关联一下
	reqUser.Roles = make([]*models.Role, 0)
	for _, roleValue := range reqUser.RolesFront {
		dbRole, err := models.GetRoleByRoleValue(roleValue)
		if err != nil {
			sc.Logger.Ctx(ctx).Error("根据RolesFront去db中寻找角色失败", zap.Error(err))
			common.FailWithMessage(err.Error(), c)
			return
		}
		reqUser.Roles = append(reqUser.Roles, dbRole)
	}

	// 密码我们要 加密处理
	reqUser.Password = common.BcryptHash(reqUser.Password)
	err = reqUser.CreateOne()
	if err != nil {
		sc.Logger.Ctx(ctx).Error("创建用户错误", zap.Error(err))
		common.FailWithMessage(err.Error(), c)
		return
	}
	sc.Logger.Ctx(ctx).Info("创建用户成功", zap.Any("user", reqUser))
	common.OkWithMessage("创建成功", c)
}

func AccountExist(c *gin.Context) {
	sc := c.MustGet(common.GIN_CTX_CONFIG_CONFIG).(*config.ServerConfig)
	// 校验一下 menu字段
	var reqUser models.AccountExistRequest
	err := c.ShouldBindJSON(&reqUser)
	if err != nil {
		sc.Logger.Error("解析编辑用户请求失败", zap.Any("用户", reqUser), zap.Error(err))
		common.FailWithMessage(err.Error(), c)
		return
	}

	// 在这里校验字段，是否必填，范围是否正确
	err = validate.Struct(reqUser)
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
	// 先去db中根据id找到这个user
	dbUser, _ := models.GetUserByUserName(reqUser.Account)
	if dbUser != nil {
		sc.Logger.Info("用户已存在", zap.Any("用户", reqUser))
		common.FailWithMessage("用户已存在", c)
		return
	}
	common.OkWithMessage("用户名不存在 可用", c)
}

func UpdateAccount(c *gin.Context) {
	sc := c.MustGet(common.GIN_CTX_CONFIG_CONFIG).(*config.ServerConfig)
	// 校验一下 menu字段
	var reqUser models.User
	err := c.ShouldBindJSON(&reqUser)
	if err != nil {
		sc.Logger.Error("解析编辑用户请求失败", zap.Any("用户", reqUser), zap.Error(err))
		common.FailWithMessage(err.Error(), c)
		return
	}

	// 在这里校验字段，是否必填，范围是否正确
	err = validate.Struct(reqUser)
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

	sc.Logger.Info("编辑用户请求字段打印", zap.Any("用户", reqUser))

	// 先 去db中根据id找到这个user

	_, err = models.GetUserById(int(reqUser.ID))
	if err != nil {
		sc.Logger.Error("根据id找user错误", zap.Any("菜单", reqUser), zap.Error(err))
		common.FailWithMessage(err.Error(), c)
		return
	}

	// 根据 rolesFront 去db中查询 role ，把role给他关联一下

	reqUser.Roles = make([]*models.Role, 0)
	for _, roleValue := range reqUser.RolesFront {
		dbRole, err := models.GetRoleByRoleValue(roleValue)
		if err != nil {
			sc.Logger.Error("根据RolesFront去db中找角色失败", zap.Any("用户", reqUser), zap.Error(err))
			common.FailWithMessage(err.Error(), c)
			return
		}
		reqUser.Roles = append(reqUser.Roles, dbRole)
	}
	// update 更新这个个体
	sc.Logger.Info("更新用户打印值", zap.Any("用户", reqUser))
	err = reqUser.UpdateOne(reqUser.Roles)
	if err != nil {
		sc.Logger.Error("更新用户错误", zap.Any("用户", reqUser), zap.Error(err))
		common.FailWithMessage(err.Error(), c)
		return
	}

	common.OkWithMessage("创建成功", c)
}

func ChangePassword(c *gin.Context) {
	// 拿到这个人用户 对应的role列表
	// 	遍历 role列表 找到 Menu list
	// 在拼装父子结构 返回的是数组 第一层 father 第2层children

	// 我得拿到 userCliams

	userName := c.MustGet(common.GIN_CTX_JWT_USER_NAME).(string)
	sc := c.MustGet(common.GIN_CTX_CONFIG_CONFIG).(*config.ServerConfig)
	dbUser, err := models.GetUserByUserName(userName)
	if err != nil {
		sc.Logger.Error("通过token解析到的userName去数据库中找User失败",
			zap.Error(err),
		)
		common.FailWithMessage(fmt.Sprintf("通过token解析到的userName去数据库中找User失败:%v", err.Error()), c)
		return
	}

	// 校验一下 menu字段
	var reqChange models.ChangePasswordRequest
	err = c.ShouldBindJSON(&reqChange)
	if err != nil {
		sc.Logger.Error("解析修改密码请求失败", zap.Any("用户", reqChange), zap.Error(err))
		common.FailWithMessage(err.Error(), c)
		return
	}

	// 在这里校验字段，是否必填，范围是否正确
	err = validate.Struct(reqChange)
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

	// 先校验旧密码是否对
	// 对比password
	err = common.BcryptCheck(reqChange.PasswordOld, dbUser.Password)
	if err != nil {
		sc.Logger.Error("旧密码错误", zap.Any("用户", reqChange), zap.Error(err))
		common.FailWithMessage("旧密码错误", c)
		return
	}

	// 变更密码
	// 密码我们要 加密处理
	dbUser.Password = common.BcryptHash(reqChange.PasswordNew)
	err = dbUser.UpdateOne(dbUser.Roles)

	if err != nil {
		sc.Logger.Error("解析修改密码请求失败", zap.Any("用户", reqChange), zap.Error(err))
		common.FailWithMessage(err.Error(), c)
		return
	}
	common.OkWithMessage("密码修改成功", c)

}

func GetAccountList(c *gin.Context) {

	sc := c.MustGet(common.GIN_CTX_CONFIG_CONFIG).(*config.ServerConfig)

	// 	数据库中拿到所有的menu列表

	users, err := models.GetUserAll()
	if err != nil {
		sc.Logger.Error("去数据库中拿所有的用户错误",
			zap.Error(err),
		)
		common.FailWithMessage(fmt.Sprintf("去数据库中拿所有的用户错误:%v", err.Error()), c)
		return
	}
	resp := &ResponseResourceCommon{
		Total: len(users),
		Items: users,
	}
	common.OkWithDetailed(resp, "ok", c)
}

func GetAllUserAndRoles(c *gin.Context) {

	sc := c.MustGet(common.GIN_CTX_CONFIG_CONFIG).(*config.ServerConfig)

	// 	数据库中拿到所有的menu列表

	users, err := models.GetUserAll()
	if err != nil {
		sc.Logger.Error("去数据库中拿所有的用户错误",
			zap.Error(err),
		)
		common.FailWithMessage(fmt.Sprintf("去数据库中拿所有的用户错误:%v", err.Error()), c)
		return
	}

	roles, err := models.GetRollAll()
	if err != nil {
		sc.Logger.Error("去数据库中拿所有的角色错误",
			zap.Error(err),
		)
		common.FailWithMessage(fmt.Sprintf("去数据库中拿所有的角色错误:%v", err.Error()), c)
		return
	}

	res := []DefineUserOrGroup{}
	for _, user := range users {
		user := user

		key := fmt.Sprintf("%s@%s", "用户", user.Username)

		one := DefineUserOrGroup{
			Label: key,
			Value: key,
		}
		res = append(res, one)

	}

	for _, role := range roles {
		role := role

		key := fmt.Sprintf("%s@%s", "组", role.RoleValue)

		one := DefineUserOrGroup{
			Label: key,
			Value: key,
		}
		res = append(res, one)

	}
	common.OkWithDetailed(res, "ok", c)

}

func DeleteAccount(c *gin.Context) {
	sc := c.MustGet(common.GIN_CTX_CONFIG_CONFIG).(*config.ServerConfig)
	// 校验一下 menu字段
	id := c.Param("id")
	sc.Logger.Info("删除用户", zap.Any("id", id))

	// 先去db中根据id找到这个user
	intVar, _ := strconv.Atoi(id)
	dbUser, err := models.GetUserById(intVar)
	if err != nil {
		sc.Logger.Error("根据id找user错误", zap.Any("用户", id), zap.Error(err))
		common.FailWithMessage(err.Error(), c)
		return
	}

	err = dbUser.DeleteOne()
	if err != nil {
		sc.Logger.Error("根据id删除user错误", zap.Any("用户", id), zap.Error(err))
		common.FailWithMessage(err.Error(), c)
		return
	}
	common.OkWithMessage("删除成功", c)
}
