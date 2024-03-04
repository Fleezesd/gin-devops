package view

import (
	"github.com/fleezesd/gin-devops/src/web/middleware"
	"github.com/gin-gonic/gin"
)

func ConfigRoutes(r *gin.Engine) {
	base := r.Group("")
	{
		base.GET("/ping", ping)
		base.POST("/login", UserLogin)
	}

	// 下面 group 都需有认证
	afterLoginApiGroup := r.Group("/api")
	afterLoginApiGroup.Use(middleware.JWTAuthMiddleware())
	{
		afterLoginApiGroup.GET("/getUserInfo", GetUserInfoAfterLogin)
		afterLoginApiGroup.GET("/getPermCode", GetPermCode)
	}
	systemApiGroup := afterLoginApiGroup.Group("system")
	{
		systemApiGroup.GET("/getMenuList", GetMenuList)
	}

}
