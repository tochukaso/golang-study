package db

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func GetDB() *gorm.DB {
	dsn := "gostudy:gostudy@/gostudy?parseTime=true"
	db, _ := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	return db
}
