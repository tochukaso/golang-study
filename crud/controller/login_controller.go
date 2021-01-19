package controller

import (
	"fmt"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"omori.jp/model"
	"omori.jp/pagination"
)

func ShowLogin(c *gin.Context) {
	countUp(c)
	c.HTML(http.StatusOK, "login.tmpl", gin.H{})

}

func AttemptLogin(c *gin.Context) {
	countUp(c)
	userCode := c.PostForm("userCode")
	password := c.PostForm("password")
	dbUser := model.GetUserFromCode(userCode)

	if bcrypt.CompareHashAndPassword([]byte(dbUser.Password), []byte(password)) == nil {
		saveLoginInfo(c, dbUser)
		products, count := model.ReadProduct("", "")
		page := 1
		pageSize := 10
		c.HTML(http.StatusOK, "product_index.tmpl", gin.H{
			"products":   products,
			"count":      count,
			"page":       page,
			"pageSize":   pageSize,
			"pagination": pagination.Pagination(count, page, pageSize),
		})
	} else {
		c.HTML(http.StatusOK, "login.tmpl", gin.H{
			"errMsg": "ユーザーコードとパスワードの組み合わせが一致しません。",
		})
	}

}

func Logout(c *gin.Context) {
	session := sessions.Default(c)
	session.Set("UserID", "")
	session.Clear()
	session.Save()
	c.HTML(http.StatusOK, "login.tmpl", gin.H{
		"errMsg": "ログアウトしました",
	})

}

func saveLoginInfo(c *gin.Context, user model.User) {
	session := sessions.Default(c)
	session.Set("UserID", int(user.ID))
	session.Save()
}

func countUp(c *gin.Context) {
	session := sessions.Default(c)
	var count int
	v := session.Get("count")
	if v == nil {
		count = 0
	} else {
		count = v.(int)
		count++
	}
	session.Set("count", count)
	session.Save()
	fmt.Println("Count", count)
}
