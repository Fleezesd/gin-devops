package models

import "gorm.io/gorm"

type Api struct {
	gorm.Model
	Path     string  `json:"path" gorm:"type:varchar(50);comment:路由路径"`
	Method   string  `json:"method" gorm:"type:varchar(50);comment:http请求方法"`
	Pid      int     `json:"pId" gorm:"comment:apiGroups 父级id"`
	Title    string  `json:"title" gorm:"type:varchar(50);uniqueIndex;comment:名称"`
	Roles    []*Role `json:"roles" gorm:"many2many:role_apis;"`
	Type     string  `json:"type" gorm:"type:varchar(5);comment:类型 0=父级 1=子级"`
	Key      uint    `json:"key"  gorm:"-"`
	Value    uint    `json:"value"  gorm:"-"`
	Children []*Api  `json:"children" gorm:"-"`
}
