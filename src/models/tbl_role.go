package models

import (
	"errors"
	"fmt"

	"github.com/fleezesd/gin-devops/src/config"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
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
	ApiIds    []int   `json:"apiIds" gorm:"-"`
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

func (r *Role) UpdateApis(apis []*Api, sc *config.ServerConfig) error {
	err1 := Db.Where("id = ?", r.ID).Updates(r).Error
	err2 := Db.Model(r).Association("Apis").Replace(apis)
	// 还应该遍历
	rules := [][]string{}
	for _, api := range apis {
		oneRule := []string{
			r.RoleValue,
			api.Path,
			api.Method,
		}
		_, err := CasbinAddOnePolicy(r.RoleValue, api.Path, api.Method)
		if err != nil {
			sc.Logger.Error("CasbinAddOnePolicy错误",
				zap.Error(err),
				zap.String("角色", r.RoleValue),
				zap.String("api.Path", api.Path),
				zap.String("api.Method", api.Method),
			)
		}

		rules = append(rules, oneRule)
	}
	_, err3 := CasbinAddPolicies(rules)

	if err1 == nil && err2 == nil && err3 == nil {
		return nil
	} else {
		return fmt.Errorf("更新本体:w 更新关联:%w 更新casbin:%w", err1, err2, err3)
	}
}

func GetRoleById(id int) (*Role, error) {
	var dbRole Role
	err := Db.Where("id = ?", id).Preload("Menus").First(&dbRole).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("角色不存在")
		}
		return nil, fmt.Errorf("数据库错误:%w", err)
	}
	return &dbRole, nil
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

func GetRoleAll() (roles []*Role, err error) {
	err = Db.Preload("Menus").Preload("Users").Preload("Apis").Find(&roles).Error
	return
}

func (r *Role) CreateOne() error {
	return Db.Create(r).Error

}

func (r *Role) DeleteOne() error {
	return Db.Select(clause.Associations).Unscoped().Delete(r).Error
}
