package models

import (
	"fmt"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

/*
{
    path: '/system',
    name: 'System',
    component: 'LAYOUT',
    redirect: '/system/account',
    icon: 'ion:settings-outline',
    title: '系统管理',
    id: "1",
    dbId: "1",
    type: "0",
    show: "1",
    parentMenu: "0",
    orderNo: 0,
  },
*/

type Menu struct {
	gorm.Model
	Name       string    `json:"name" gorm:"type:varchar(100);uniqueIndex;comment:名称"`
	Title      string    `json:"title" gorm:"comment:名称"`
	Pid        int       `json:"pId" gorm:"comment:父级的id"`
	ParentMenu string    `json:"parentMenu" gorm:"varchar(5);comment:父级的id"`
	Meta       *MenuMeta `json:"meta" gorm:"-"`
	Icon       string    `json:"icon" gorm:"comment:图标"`
	Type       string    `json:"type" gorm:"type:varchar(5);comment:类型 0=目录 1=子菜单"`
	Show       string    `json:"show" gorm:"type:varchar(5);comment:类型 0=禁用 1=启用"`
	OrderNo    int       `json:"orderNo" gorm:"comment:排序"`
	Component  string    `json:"component" gorm:"type:varchar(50);comment:前端组件 菜单就是LAYOUT"`
	Redirect   string    `json:"redirect" gorm:"type:varchar(50);comment:显示路径"`
	Path       string    `json:"path" gorm:"type:varchar(50);comment:路由路径"`
	Remark     string    `json:"remark" gorm:"comment:用户描述"`
	HomePath   string    `json:"homePath" gorm:"comment:登陆后的默认首页"`
	Status     string    `json:"status" gorm:"default:1;comment:是否启用 0禁用 1启用"` //用户是否被冻结 1正常 2冻结
	Roles      []*Role   `json:"roles" gorm:"many2many:role_menus;"`
	Children   []*Menu   `json:"children" gorm:"-"`
	Key        uint      `json:"value"  gorm:"-"`
	Value      uint      `json:"key"  gorm:"-"`
}

type MenuMeta struct {
	Title           string `json:"title" gorm:"-"`
	Icon            string `json:"icon" gorm:"-"`
	ShowMenu        bool   `json:"showMenu" gorm:"-"`
	HideMenu        bool   `json:"hideMenu" gorm:"-"`
	IgnoreKeepAlive bool   `json:"ignoreKeepAlive" gorm:"-"`
}

func GetMenuById(id int) (*Menu, error) {
	var dbMenu Menu
	err := Db.Where("id = ?", id).Preload("Roles").First(&dbMenu).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("菜单不存在")
		}
		return nil, fmt.Errorf("数据库错误:%w", err)
	}
	return &dbMenu, nil
}

func GetMenuAll() (menus []*Menu, err error) {
	err = Db.Find(&menus).Error
	return
}

func (m *Menu) UpdateOne() error {
	return Db.Updates(m).Error
}

func (m *Menu) CreateOne() error {
	return Db.Create(m).Error
}

func (m *Menu) DeleteOne() error {
	// Select(clause.Associations) 用于自动预加载 Menu 结构体的所有关联
	// Unscoped() 用于执行硬删除，这将永久性地从数据库中删除记录
	// Delete(m) 删除与 Menu 实例 m 的主键匹配的记录
	// Error 用于返回在删除操作过程中可能发生的任何错误。如果没有发生错误，它将返回 nil
	return Db.Select(clause.Associations).Unscoped().Delete(m).Error
}
