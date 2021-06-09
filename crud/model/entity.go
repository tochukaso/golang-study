package model

import (
	"github.com/tochukaso/golang-study/db"
	"gorm.io/gorm"
)

// Entity が兼ね備えるべきCRUDのメソッドを定義する。
type Entity interface {
	GetID() uint
	GetFromCode() Entity
	Create() error
	Read() Entity
	Update() error
	Delete() error
}

func GetDB() *gorm.DB {
	return db.GetDB()
}
