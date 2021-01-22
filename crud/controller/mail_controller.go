package controller

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	csrf "github.com/utrack/gin-csrf"
	"omori.jp/message"
	"omori.jp/model"
)

func InitMailTemplate() {
	model.InitMailTemplate()
}

func ShowMailTemplates(c *gin.Context) {
	templates := model.ListMailType()
	var list []map[string]interface{}
	for _, t := range templates {
		log.Println("t", t)
		template := model.GetMailTemplateFromCode(int(t))
		mailType := model.GetMailType(t)
		v := make(map[string]interface{})
		v["code"] = t
		v["visibleCode"] = mailType[0]
		v["name"] = mailType[1]
		v["subject"] = template.Subject
		v["updatedAt"] = template.UpdatedAt
		list = append(list, v)
	}

	RenderHTML(c, http.StatusOK, "mail_index.tmpl", gin.H{
		"list": list,
	})

}

func GetMailTemplate(c *gin.Context) {
	rawCode := c.Param("code")
	log.Println("code", rawCode)

	code, err := strconv.Atoi(rawCode)
	if err != nil {
		log.Println("code(%v)が不正です。", rawCode)
		c.Redirect(301, "/mail/")
		return
	}

	mail := model.GetMailTemplateFromCode(code)
	log.Println(mail)

	if mail.ID == 0 {
		mail.MailCode = code
		mail.Subject = model.GetMailType(model.MailType(code))[1]
	}

	RenderHTML(c, http.StatusOK, "mail_detail.tmpl", gin.H{
		"P":         mail,
		"variables": model.ListTemplateValiable(model.MailType(mail.MailCode)),
	})
}

func PutMailTemplate(c *gin.Context) {
	var mailTemplate model.MailTemplate
	err := c.ShouldBind(&mailTemplate)
	validate := validator.New()
	validate.RegisterStructValidation(checkDuplicateMailCode, model.MailTemplate{})
	errors := validate.Struct(mailTemplate)
	if err != nil || errors != nil {
		errs := errors.(validator.ValidationErrors)
		sliceErrs := []string{}
		for _, e := range errs {
			sliceErrs = append(sliceErrs, message.ConvertMessage(e))
		}
		RenderHTML(c, http.StatusOK, "mail_detail.tmpl", gin.H{
			"P":      mailTemplate,
			"errMsg": sliceErrs,
		})
		return
	}

	isFirst := mailTemplate.ID == 0
	var msg string
	if isFirst {
		err := mailTemplate.Create()
		if err != nil {
			RenderHTML(c, http.StatusOK, "mail_detail.tmpl", gin.H{
				"P":      mailTemplate,
				"errMsg": "メールテンプレートの登録に失敗しました",
			})
			return
		}
		msg = "登録しました"
	} else {
		old := mailTemplate.Read().(model.MailTemplate)
		mailTemplate.CreatedAt = old.CreatedAt
		mailTemplate.Update()
		msg = "保存しました"
	}

	RenderHTML(c, http.StatusOK, "mail_detail.tmpl", gin.H{
		"_csrf": csrf.GetToken(c),
		"P":     mailTemplate,
		"msg":   msg,
	})
}

func checkDuplicateMailCode(sl validator.StructLevel) {
	template := sl.Current().Interface().(model.MailTemplate)

	if template.MailCode == 0 || template.ID != 0 {
		return
	}
	dbMailTemplate := model.GetMailTemplateFromCode(template.MailCode)
	log.Println("dbMailTemplate", dbMailTemplate)

	if dbMailTemplate.GetID() != 0 {
		sl.ReportError(template.MailCode, "MailCode", "MailCode", "duplicateCode", "")
	}
}
