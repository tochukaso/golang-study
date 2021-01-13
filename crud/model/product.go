package model

import (
	"fmt"

	"gorm.io/gorm"
	"omori.jp/db"
)

type Product struct {
	gorm.Model
	Name    string
	OrgCode string
	JanCode string
	Detail  string
}

func InitProduct() {
	db := db.GetDB()
	db.AutoMigrate(&Product{})
}

func CreateProduct(product *Product) {
	db := db.GetDB()
	db.Create(product)
}

func GetProductFromId(id string) Product {
	db := db.GetDB()
	var product Product
	db.First(&product, id)
	return product
}

func ReadProduct(orgCode, name string) []Product {
	db := db.GetDB()
	var products []Product
	db.Find(&products, "org_code LIKE ? and name LIKE ?",
		fmt.Sprintf("%%%s%%", orgCode),
		fmt.Sprintf("%%%s%%", name),
	)
	return products
}

func UpdateProduct(product Product) {
	db := db.GetDB()
	db.Save(product)
}

func DeleteProduct(id string) {
	db := db.GetDB()
	db.Delete(&Product{}, id)
}
