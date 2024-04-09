package models

import (
	"fmt"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

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

func GetApiById(id int) (*Api, error) {
	var dbObj Api
	err := Db.Where("id = ?", id).Preload("Roles").First(&dbObj).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("api不存在")
		}
		return nil, fmt.Errorf("数据库错误:%w", err)
	}
	return &dbObj, nil
}

func GetApiAll() (objs []*Api, err error) {
	err = Db.Find(&objs).Error
	return
}

func (obj *Api) DeleteOne() error {
	return Db.Select(clause.Associations).Unscoped().Delete(obj).Error
}

func (obj *Api) CreateOne() error {
	return Db.Create(obj).Error

}

func (obj *Api) UpdateOne() error {
	return Db.Updates(obj).Error

}
