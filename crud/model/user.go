package model

import (
	"fmt"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	UserCode string `form:"UserCode" binding:"required" validate:"required,ascii" gorm:"unique;not null"`
	UserName string `form:"UserName" binding:"required" validate:"required" gorm:"not null"`
	Password string `form:"Password" binding:"required" validate:"required,ascii,gte=8" gorm:"not null"`
	Mail     string `form:"Mail" validate:"email" `
}

func InitUser() {
	GetDB().AutoMigrate(&User{})
}

func (e User) Create() error {
	result := GetDB().Create(&e)
	return result.Error
}

func (e User) Read() Entity {
	var res = e
	GetDB().Find(&res)
	return res
}

func (e User) Update() error {
	result := GetDB().Save(&e)
	return result.Error
}

func (e User) Delete() error {
	result := GetDB().Delete(&e)
	return result.Error
}

func (e User) GetID() uint {
	return e.ID
}

func (e User) GetFromCode() Entity {
	var user User
	GetDB().Find(&user, "user_code = ?", e.UserCode)
	return user
}

func GetUserFromId(id string) User {
	var user User
	GetDB().First(&user, id)
	return user
}

func GetUserFromCode(code string) User {
	var user User
	GetDB().Find(&user, "user_code = ?", code)
	return user
}

func ReadUserWithPaging(page, pageSize int, code, name string) ([]User, int) {
	db := GetDB()
	offset := (page - 1) * pageSize
	db = GetDB().Offset(offset).Limit(pageSize)
	return readUser(db, code, name)
}

func readUser(gdb *gorm.DB, code, name string) ([]User, int) {
	var users []User
	var count int64
	where := []interface{}{"user_code LIKE ? and user_name LIKE ?"}
	args := []interface{}{fmt.Sprintf("%%%s%%", code),
		fmt.Sprintf("%%%s%%", name)}

	GetDB().Find(&users, append(where, args...)...)
	fmt.Println("users", users)
	GetDB().Model(&User{}).Where(where[0], args...).Count(&count)

	fmt.Println(count)
	return users, int(count)
}
