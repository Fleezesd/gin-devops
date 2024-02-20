package models

import (
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
