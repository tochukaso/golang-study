package model

import (
	"time"

	"gorm.io/gorm"
)

type Product struct {
	ID            uint `gorm:"primaryKey" json:"id"`
	CreatedAt     time.Time
	UpdatedAt     time.Time
	DeletedAt     gorm.DeletedAt `gorm:"index"`
	ProductName   string         `form:"ProductName" validate:"required" gorm:"not null" json:"title"`
	OrgCode       string         `form:"OrgCode" validate:"required,ascii" gorm:"unique;not null" `
	JanCode       string         `form:"JanCode" validate:"ascii"`
	ProductDetail string         `form:"ProductDetail" json:"description"`
	ProductPrice  int            `form:"ProductPrice" json:"price"`
	Rating        int            `form:"Rating" json:"ratings"`
	Review        int            `form:"Review" json:"reviews"`
	ProductImage  string         `form:"ProductImage" json:"image"`
}
