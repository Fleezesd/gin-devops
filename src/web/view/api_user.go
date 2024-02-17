package view

import (
	"github.com/fleezesd/gin-devops/src/common"
	"github.com/fleezesd/gin-devops/src/models"
	"github.com/gin-gonic/gin"
)

func UserLogin(c *gin.Context) {
	var user models.UserLoginRequest
	if err := c.ShouldBindJSON(&user); err != nil {
		common.FailWithMessage(err.Error(), c)
		return
	}
	common.OkWithData(user, c)
}
