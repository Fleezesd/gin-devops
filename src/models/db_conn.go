package models

import (
	"github.com/fleezesd/gin-devops/src/common"
	"github.com/fleezesd/gin-devops/src/config"
	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var (
	Db *gorm.DB
)

func InitDB(dsn string) error {
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return err
	}
	Db = db

	return nil
}

// MigrateTable 自动建表的逻辑
func MigrateTable() error {
	return Db.AutoMigrate(
		&User{},
	)
}

func MockUserRegister(sc *config.ServerConfig) {
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
			},
		},
	}
	user1.Password = common.BcryptHash(user1.Password)
	if err := Db.Create(&user1).Error; err != nil {
		sc.Logger.Error("模拟用户注册失败",
			zap.Error(err),
		)
		return
	}
	sc.Logger.Info("模拟用户注册成功!")
}
