package models

import (
	"context"
	"errors"
	"fmt"

	"github.com/fleezesd/gin-devops/src/common"
	"github.com/uptrace/opentelemetry-go-extra/otelzap"
	"go.opentelemetry.io/otel"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model         // gorm 预定义结构体
	UserId     int     `json:"userId" gorm:"comment:用户id"`
	Username   string  `json:"username" gorm:"type:varchar(100);uniqueIndex;comment:用户登录名"` // 记得指定长度 要不migrate报错
	Password   string  `json:"password" gorm:"comment:用户登录密码"`
	RealName   string  `json:"realName" gorm:"comment:用户昵称"`
	Desc       string  `json:"desc" gorm:"comment:用户描述"`
	HomePath   string  `json:"homePath" gorm:"comment:登陆后的默认首页"`
	Enable     int     `json:"enable" gorm:"default:1;comment:用户是否被冻结 1正常 2冻结"`
	Roles      []*Role `json:"roles" gorm:"many2many:user_roles;"`
}

var tracer = otel.Tracer("gin-devops-models")

func CheckUserPassword(logger *otelzap.Logger, ctx context.Context, req *UserLoginRequest) (*User, error) {
	dbUser := User{
		Username: req.Username,
	}
	err := Db.WithContext(ctx).Where("username = ?", dbUser.Username).Preload("Roles").First(&dbUser).Error
	if err != nil {
		logger.Ctx(ctx).Error("数据库错误", zap.Error(err))
		return nil, err
	}
	// 跟db中加密的密码对比
	if err = common.BcryptCheck(req.Password, dbUser.Password); err != nil {
		logger.Ctx(ctx).Error("密码错误", zap.Error(err))
		return nil, err
	}
	return &dbUser, nil
}

func GetUserByUserName(userName string) (*User, error) {
	dbUser := User{
		Username: userName,
	}
	err := Db.Where("username = ?", dbUser.Username).Preload("Roles.Menus").First(&dbUser).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.New("用户名不存在")
		}
		return nil, fmt.Errorf("数据库错误: %w", err)
	}
	return &dbUser, nil
}
