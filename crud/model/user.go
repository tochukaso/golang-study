package model

import (
	"fmt"

	"gorm.io/gorm"
	"omori.jp/db"
)

type User struct {
	gorm.Model
	UserCode string `form:"UserCode" validate:"required,ascii"`
	UserName string `form:"UserName" binding:"required" validate:"required"`
	Password string `form:"Password" binding:"required" validate:"required,ascii,gte=8"`
}

func InitUser() {
	db := db.GetDB()
	db.AutoMigrate(&User{})
}

func (e *User) Create() {
	db := db.GetDB()
	db.Create(e)
}

func (e *User) Read() {
	db := db.GetDB()
	db.Find(&e)
}

func (e *User) Update() {
	db := db.GetDB()
	db.Save(e)
}

func (e *User) Delete() {
	db := db.GetDB()
	db.Delete(e)
}

func GetUserFromId(id string) User {
	db := db.GetDB()
	var user User
	db.First(&user, id)
	return user
}

func GetUserFromCode(code string) User {
	db := db.GetDB()
	var user User
	db.Find(&user, "user_code = ?", code)
	return user
}

func ReadUserWithPaging(page, pageSize int, code, name string) ([]User, int) {
	db := db.GetDB()
	offset := (page - 1) * pageSize
	db = db.Offset(offset).Limit(pageSize)
	return readUser(db, code, name)
}

func readUser(gdb *gorm.DB, code, name string) ([]User, int) {
	var users []User
	var count int64
	where := []interface{}{"user_code LIKE ? and user_name LIKE ?"}
	args := []interface{}{fmt.Sprintf("%%%s%%", code),
		fmt.Sprintf("%%%s%%", name)}

	gdb.Find(&users, append(where, args...)...)
	fmt.Println("users", users)
	db.GetDB().Model(&User{}).Where(where[0], args...).Count(&count)

	fmt.Println(count)
	return users, int(count)
}
