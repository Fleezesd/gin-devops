package models

import "github.com/golang-jwt/jwt/v5"

type UserLoginRequest struct {
	Username string `json:"username"  validate:"required,min=3,max=20"` // 用户名
	Password string `json:"password"  validate:"required,min=3,max=20"` // 密码
	//Email    string `json:"email"     validate:"required,email"`        // 邮箱
	//Gender   string `json:"gender"    validate:"oneof=male female"`     // 性别
}

// UserCustomClaims 根据token 解析成的对象
type UserCustomClaims struct {
	*User
	jwt.RegisteredClaims // 内嵌的标准声明
}

// UserLoginResponse login接口返回对象
type UserLoginResponse struct {
	*User
	Token     string `json:"token"`
	ExpiresAt int64  `json:"expiresAt"` // 为了方便前端 直接拿到过期时间
}

type AccountExistRequest struct {
	Account string `json:"account"`
}

type ChangePasswordRequest struct {
	PasswordOld string `json:"passwordOld"`
	PasswordNew string `json:"passwordNew"`
}
