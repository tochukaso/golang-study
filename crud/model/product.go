package model

import (
	"fmt"

	"gorm.io/gorm"
	"omori.jp/db"
)

type Product struct {
	gorm.Model
	Name    string `form:"Name" binding:"required" validate:"required"`
	OrgCode string `form:"OrgCode" validate:"required,ascii"`
	JanCode string `form:"JanCode" validate:"ascii"`
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

func ReadProduct(orgCode, name string) ([]Product, int64) {
	return readProduct(db.GetDB(), orgCode, name)
}

func ReadProductWithPaging(page, pageSize int, orgCode, name string) ([]Product, int64) {
	db := db.GetDB()
	offset := (page - 1) * pageSize
	db = db.Offset(offset).Limit(pageSize)
	return readProduct(db, orgCode, name)
}

func UpdateProduct(product Product) {
	db := db.GetDB()
	db.Save(product)
}

func DeleteProduct(id string) {
	db := db.GetDB()
	db.Delete(&Product{}, id)
}

func readProduct(gdb *gorm.DB, orgCode, name string) ([]Product, int64) {
	var products []Product
	var count int64
	where := []interface{}{"org_code LIKE ? and name LIKE ?"}
	args := []interface{}{fmt.Sprintf("%%%s%%", orgCode),
		fmt.Sprintf("%%%s%%", name)}

	gdb.Find(&products, append(where, args...)...)
	fmt.Println("products", products)
	db.GetDB().Model(&Product{}).Where(where[0], args...).Count(&count)

	fmt.Println(count)
	return products, count
}
