package model

import (
	"gorm.io/gorm"
)

type MailTemplate struct {
	gorm.Model
	MailCode   int    `form:"MailCode" validate:"required,numeric" gorm:"unique;not null"`
	Subject    string `form:"Subject" validate:"required" gorm:"not null"`
	From       string `form:"From" validate:"required" gorm:"not null"`
	CC         string `form:"CC"`
	MailDetail string `form:"MailDetail"`
}

func InitMailTemplate() {
	GetDB().AutoMigrate(&MailTemplate{})
}

func (e MailTemplate) Create() error {
	result := GetDB().Create(&e)
	return result.Error
}

func (e MailTemplate) Read() Entity {
	var res = e
	GetDB().Find(&res)
	return res
}

func (e MailTemplate) Update() error {
	result := GetDB().Save(&e)
	return result.Error
}

func (e MailTemplate) Delete() error {
	result := GetDB().Delete(&e)
	return result.Error
}

func (e MailTemplate) GetID() uint {
	return e.ID
}

func (e MailTemplate) GetFromCode() Entity {
	var user MailTemplate
	GetDB().Find(&user, "mail_code = ?", e.MailCode)
	return user
}

func GetMailTemplateFromCode(code int) MailTemplate {
	var user MailTemplate
	GetDB().Find(&user, "mail_code = ?", code)
	return user
}
