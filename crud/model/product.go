package model

import (
	"fmt"
	"log"
	"strconv"
	"time"

	"gorm.io/gorm"
	"omori.jp/db"
	"omori.jp/env"
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

type ProductJSON struct {
	Product
	ProductImagePath string `json:"imagePath"`
}

func InitProduct() {
	db := db.GetDB()
	db.AutoMigrate(&Product{})
}

func (e Product) Create() error {
	result := GetDB().Create(&e)
	return result.Error
}

func (e Product) Read() Entity {
	var res = e
	GetDB().Find(&res)
	return res
}

func (e Product) Update() error {
	result := GetDB().Save(&e)
	return result.Error
}

func (e Product) Delete() error {
	result := GetDB().Delete(&e)
	return result.Error
}

func (e Product) GetID() uint {
	return e.ID
}

func (e Product) GetFromCode() Entity {
	var user Product
	GetDB().Find(&user, "org_code = ?", e.OrgCode)
	return user
}

func (e Product) GetAbsoluteImagePath() string {
	if e.ProductImage == "" {
		return ""
	}
	return env.GetEnv().BaseURL + e.GetImagePath()
}

func (e Product) GetImagePath() string {
	if e.ProductImage == "" {
		return ""
	}
	return "static/assets/product/" + strconv.Itoa(int(e.ID)) + "/" + e.ProductImage
}

func GetProductFromId(id string) Product {
	db := db.GetDB()
	var product Product
	db.First(&product, id)
	return product
}

func GetProductFromCode(code string) Product {
	var product Product
	GetDB().Find(&product, "org_code = ?", code)
	return product
}

func ReadProduct(orgCode, productName string) ([]Product, int) {
	return readProduct(db.GetDB(), orgCode, productName)
}

func ReadProductWithPaging(page, pageSize int, orgCode, productName string) ([]Product, int) {
	db := db.GetDB()
	offset := (page - 1) * pageSize
	db = db.Offset(offset).Limit(pageSize)
	return readProduct(db, orgCode, productName)
}

func readProduct(gdb *gorm.DB, orgCode, productName string) ([]Product, int) {
	var products []Product
	var count int64
	where := []interface{}{"org_code LIKE ? and product_name LIKE ?"}
	args := []interface{}{fmt.Sprintf("%%%s%%", orgCode),
		fmt.Sprintf("%%%s%%", productName)}

	gdb.Find(&products, append(where, args...)...)
	log.Println("products", products)
	GetDB().Model(&Product{}).Where(where[0], args...).Count(&count)

	log.Println(count)
	return products, int(count)
}
