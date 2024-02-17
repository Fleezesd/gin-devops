package view

import (
	"github.com/gin-gonic/gin"
)

func ConfigRoutes(r *gin.Engine) {
	base := r.Group("/basic-api")
	{
		base.GET("/ping", ping)
		base.POST("/login", UserLogin)
	}

}
