package models

import (
	"errors"
	"fmt"

	"gorm.io/gorm"
)

type Role struct {
	gorm.Model
	OrderNo   int     `json:"orderNo" gorm:"comment:排序"`
	RoleName  string  `json:"roleName" gorm:"type:varchar(100);uniqueIndex;comment:角色中文名称"`
	RoleValue string  `json:"roleValue"  gorm:"type:varchar(100);uniqueIndex;comment:角色值"`
	Remark    string  `json:"remark" gorm:"comment:用户描述"`
	HomePath  string  `json:"homePath" gorm:"comment:登陆后的默认首页"`
	Status    string  `json:"status" gorm:"default:1;comment:角色是否被冻结 1正常 2冻结"`
	Users     []*User `json:"users" gorm:"many2many:user_roles;"`
	Menus     []*Menu `json:"menus" gorm:"many2many:role_menus;"`
	MenuIds   []int   `json:"menuIds" gorm:"-"`
	Apis      []*Api  `json:"apis" gorm:"many2many:role_apis;"`
}

func (r *Role) UpdateMenus(menus []*Menu) error {
	err1 := Db.Where("id = ?", r.ID).Updates(r).Error
	err2 := Db.Model(r).Association("Menus").Replace(menus)
	if err1 == nil && err2 == nil {
		return nil
	} else {
		return fmt.Errorf("更新本体:%w 更新关联:%w", err1, err2)
	}
}

func (r *Role) UpdateApis(apis []*Api) error {
	err1 := Db.Where("id = ?", r.ID).Updates(r).Error
	err2 := Db.Model(r).Association("Apis").Replace(apis)
	if err1 == nil && err2 == nil {
		return nil
	} else {
		return fmt.Errorf("更新本体:%w 更新关联:%w", err1, err2)
	}
}

func GetRoleByRoleValue(roleValue string) (*Role, error) {
	var dbRole Role
	err := Db.Where("role_value = ?", roleValue).Preload("Menus").First(&dbRole).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("角色不存在")
		}
		return nil, fmt.Errorf("数据库错误:%w", err)
	}
	return &dbRole, nil
}

func GetRollAll() (roles []*Role, err error) {
	err = Db.Preload("Menus").Preload("Users").Preload("Apis").Find(&roles).Error
	return
}
