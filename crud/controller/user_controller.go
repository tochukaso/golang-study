package controller

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/tochukaso/golang-study/mail"
	"github.com/tochukaso/golang-study/message"
	"github.com/tochukaso/golang-study/model"
	"github.com/tochukaso/golang-study/pagination"
	"golang.org/x/crypto/bcrypt"
)

func InitUser() {
	model.InitUser()
}

func ShowUsers(c *gin.Context) {
	userName := c.Query("userName")
	userCode := c.Query("userCode")
	log.Println("userName", userName)
	log.Println("userCode", userCode)
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))
	log.Println("page", page)
	log.Println("pageSize", pageSize)

	users, count := model.ReadUserWithPaging(page, pageSize, userCode, userName)
	log.Println(users)
	log.Println(count)

	RenderHTML(c, http.StatusOK, "user_index.tmpl", gin.H{
		"userName":   userName,
		"userCode":   userCode,
		"page":       page,
		"count":      count,
		"pageSize":   pageSize,
		"users":      users,
		"pagination": pagination.Pagination(count, page, pageSize),
	})

}

func GetUser(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	log.Println("id", id)
	user := createIDUser(id)
	user = user.Read().(model.User)
	user.Password = ""

	RenderHTML(c, http.StatusOK, "user_detail.tmpl", gin.H{
		"P": user,
	})
}

func PutUser(c *gin.Context) {
	var user model.User
	err := c.ShouldBind(&user)
	validate := validator.New()
	validate.RegisterStructValidation(checkDuplicateUserCode, model.User{})
	errors := validate.Struct(user)
	if err != nil || errors != nil {
		errs := errors.(validator.ValidationErrors)
		sliceErrs := []string{}
		for _, e := range errs {
			sliceErrs = append(sliceErrs, message.ConvertMessage(e))
		}
		RenderHTML(c, http.StatusOK, "user_detail.tmpl", gin.H{
			"P":      user,
			"errMsg": sliceErrs,
		})
		return
	}

	hash, _ := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.MinCost)
	user.Password = string(hash)

	isFirst := user.ID == 0
	var msg string
	if isFirst {
		err := user.Create()
		if err != nil {
			user.Password = ""
			RenderHTML(c, http.StatusOK, "user_detail.tmpl", gin.H{
				"P":      user,
				"errMsg": "ユーザーの登録に失敗しました",
			})
			return
		}
		msg = "登録しました"
		mail.SendUserRegisterMail(user)
	} else {
		dbUser := user.Read().(model.User)
		user.CreatedAt = dbUser.CreatedAt
		user.Update()
		msg = "保存しました"
	}
	user.Password = ""

	RenderHTML(c, http.StatusOK, "user_detail.tmpl", gin.H{
		"P":   user,
		"msg": msg,
	})
}

func DeleteUser(c *gin.Context) {
	id, _ := strconv.Atoi(c.PostForm("id"))
	user := createIDUser(id)

	user.Delete()

	RenderHTML(c, http.StatusOK, "user_index.tmpl", gin.H{
		"msg": "削除しました",
	})
}

func createIDUser(id int) model.User {
	var user model.User
	user.ID = uint(id)
	return user
}

func checkDuplicateUserCode(sl validator.StructLevel) {
	user := sl.Current().Interface().(model.User)

	if user.UserCode == "" || user.ID != 0 {
		return
	}
	dbUser := model.GetUserFromCode(user.UserCode)
	log.Println("dbUser", dbUser)

	if dbUser.GetID() != 0 {
		sl.ReportError(user.UserCode, "UserCode", "UserCode", "duplicateCode", "")
	}
}
