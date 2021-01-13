package product

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
}

func Init() {
	db := db.GetDB()
	db.AutoMigrate(&Product{})
}

func Create(product *Product) {
	db := db.GetDB()
	db.Create(product)
}

func Read(orgCode, name string) []Product {
	db := db.GetDB()
	var products []Product
	db.Find(&products, "org_code LIKE ? and name LIKE ?",
		fmt.Sprintf("%%%s%%", orgCode),
		fmt.Sprintf("%%%s%%", name),
	)
	return products
}

func Update(product Product) {
	db := db.GetDB()
	db.Save(product)
}

func Delete(product Product) {
	db := db.GetDB()
	db.Delete(&product)
}
