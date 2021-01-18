package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"omori.jp/model"
	"omori.jp/pagination"
)

func ShowLogin(c *gin.Context) {
	c.HTML(http.StatusOK, "login.tmpl", gin.H{})

}

func AttemptLogin(c *gin.Context) {
	userCode := c.PostForm("userCode")
	password := c.PostForm("password")
	dbUser := model.GetUserFromCode(userCode)

	if bcrypt.CompareHashAndPassword([]byte(dbUser.Password), []byte(password)) == nil {
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
