package model

import (
	"gorm.io/gorm"
	"omori.jp/db"
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
