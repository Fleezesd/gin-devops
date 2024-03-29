package models

import (
	"github.com/fleezesd/gin-devops/src/common"
	"github.com/fleezesd/gin-devops/src/config"
	"github.com/uptrace/opentelemetry-go-extra/otelgorm"
	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var (
	Db *gorm.DB
)

func InitDB(dsn string) error {
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return err
	}
	if err = db.Use(otelgorm.NewPlugin()); err != nil {
		return err
	}
	Db = db
	return nil
}

// MigrateTable 自动建表的逻辑
func MigrateTable() error {
	return Db.AutoMigrate(
		&User{},
		&Role{},
		&Menu{},
	)
}

func MockUserRegister(sc *config.ServerConfig) {
	menus := []*Menu{
		{
			Name:      "System",
			Icon:      "ion:settings-outline",
			Title:     "系统管理",
			Component: "LAYOUT",
			Redirect:  "/system/account",
			Path:      "/system",
		},
	}

	user1 := User{
		Username: "vben",
		Password: "123456",
		RealName: "超管",
		Desc:     "",
		HomePath: "/system/account",
		Enable:   1,
		Roles: []*Role{
			{
				RoleName:  "管理员",
				RoleValue: "admin",
				Menus:     menus,
			},
			{
				RoleName:  "前端权限管理员",
				RoleValue: "frontAdmin",
			},
		},
	}
	user1.Password = common.BcryptHash(user1.Password)
	// 保存用户信息   关联role menu会自动migrate到各自表中
	if err := Db.Create(&user1).Error; err != nil {
		sc.Logger.Error("模拟用户注册失败", zap.Error(err))
		return
	}
	sc.Logger.Info("模拟用户注册成功!")
}
