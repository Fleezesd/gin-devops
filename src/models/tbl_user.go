package models

import "gorm.io/gorm"

type User struct {
	gorm.Model        // gorm 预定义结构体
	UserId     int    `json:"userId" gorm:"comment:用户id"`
	Username   string `json:"username" gorm:"index;comment:用户登录名"`
	Password   string `json:"password" gorm:"comment:用户登录密码"`
	RealName   string `json:"realName" gorm:"comment:用户昵称"`
	Desc       string `json:"desc" gorm:"comment:用户描述"`
	HomePath   string `json:"homePath" gorm:"comment:登陆后的默认首页"`
	Enable     int    `json:"enable" gorm:"default:1;comment:用户是否被冻结 1正常 2冻结"`
}
