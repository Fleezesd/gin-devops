package view

import (
	"github.com/fleezesd/gin-devops/src/web/middleware"
	"github.com/gin-gonic/gin"
)

func ConfigRoutes(r *gin.Engine) {
	base := r.Group("/")
	{
		base.GET("/ping", ping)
		base.POST("/login", UserLogin)
	}

	// 下面 group 都需有认证
	afterLoginApiGroup := r.Group("/api")
	afterLoginApiGroup.Use(middleware.JWTAuthMiddleware())
	afterLoginApiGroup.Use(middleware.CasbinAuthMiddleware())
	{
		afterLoginApiGroup.GET("/getUserInfo", GetUserInfoAfterLogin)
		afterLoginApiGroup.GET("/getPermCode", GetPermCode)
	}
	systemApiGroup := afterLoginApiGroup.Group("system")
	{
		// 用户相关 account
		systemApiGroup.POST("/createAccount", CreateAccount)
		systemApiGroup.POST("/accountExist", AccountExist)
		systemApiGroup.POST("/updateAccount", UpdateAccount)
		systemApiGroup.POST("/changePassword", ChangePassword)
		systemApiGroup.GET("/getAccountList", GetAccountList)
		systemApiGroup.GET("/getAllUserAndRoles", GetAllUserAndRoles)
		systemApiGroup.DELETE("/deleteAccount/:id", DeleteAccount)

		// 角色相关 role
		systemApiGroup.GET("/getRoleListAll", GetRoleListAll)
		systemApiGroup.POST("/createRole", CreateRole)
		systemApiGroup.POST("/updateRole", UpdateRole)
		systemApiGroup.POST("/setRoleStatus", SetRoleStatus)
		systemApiGroup.DELETE("/deleteRole/:id", DeleteRole)

		// 菜单相关 menu
		systemApiGroup.GET("/getMenuList", GetMenuList)
		systemApiGroup.GET("/getMenuListAll", GetMenuListAll)
		systemApiGroup.POST("/updateMenu", UpdateMenu)
		systemApiGroup.POST("/createMenu", CreateMenu)
		systemApiGroup.DELETE("/deleteMenu/:id", DeleteMenu)

		// api 相关
		systemApiGroup.GET("/getApiList", GetApiList)
		systemApiGroup.DELETE("/deleteApi/:id", DeleteApi)
		systemApiGroup.POST("/createApi", CreateApi)
		systemApiGroup.POST("/updateApi", UpdateApi)
	}

}
