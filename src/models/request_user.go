package models

type UserLoginRequest struct {
	Username string `json:"username"  validate:"required,min=3,max=20"` // 用户名
	Password string `json:"password" validate:"required,min=3,max=20"`  // 密码
}
