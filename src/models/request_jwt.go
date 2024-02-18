package models

import (
	"time"

	"github.com/fleezesd/gin-devops/src/common"
	"github.com/fleezesd/gin-devops/src/config"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"go.uber.org/zap"
)

func TokenNext(dbUser *User, c *gin.Context) {
	sc := c.MustGet(common.GIN_CTX_CONFIG_CONFIG).(*config.ServerConfig)
	token, err := GenJwtToken(dbUser, sc)
	if err != nil {
		sc.Logger.Error("生成jwt失败", zap.Error(err))
		common.FailWithMessage("生成jwt失败", c)
		return
	}
	// 构造返回数据
	userResp := UserLoginResponse{
		User:  dbUser,
		Token: token,
	}
	common.OkWithDetailed(userResp, "用户登录成功", c)
}

// GenJwtToken 生成Jwt
func GenJwtToken(dbUser *User, sc *config.ServerConfig) (string, error) {
	c := UserCustomClaims{
		User: dbUser,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    sc.JWT.Issuer,
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(sc.JWT.ExpiresDuration)),
		},
	}
	//使用sha256签名方法创建签名对象
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, c)

	//使用指定的secret签名 获得加盐后的token
	return token.SignedString([]byte(sc.JWT.SigningKey))
}

// ParseToken 解析Token
func ParseToken(jwtLongToken string, sc *config.ServerConfig) (*UserCustomClaims, error) {
	tokenClaims, err := jwt.ParseWithClaims(
		jwtLongToken,
		&UserCustomClaims{},
		func(token *jwt.Token) (i interface{}, e error) {
			return []byte(sc.JWT.SigningKey), nil
		},
	)
	if err != nil {
		sc.Logger.Error("根据长 tokenString 解析 jwt错误",
			zap.Error(err))
		return nil, err
	}
	if claims, ok := tokenClaims.Claims.(*UserCustomClaims); ok && tokenClaims.Valid {
		return claims, nil
	}
	return nil, err
}
