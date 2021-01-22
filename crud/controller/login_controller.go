package controller

import (
	"log"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"omori.jp/model"
	"omori.jp/pagination"
)

func ShowLogin(c *gin.Context) {
	countUp(c)
	RenderHTML(c, http.StatusOK, "login.tmpl", gin.H{})
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
		RenderHTML(c, http.StatusOK, "product_index.tmpl", gin.H{
			"products":   products,
			"count":      count,
			"page":       page,
			"pageSize":   pageSize,
			"pagination": pagination.Pagination(count, page, pageSize),
		})
	} else {
		RenderHTML(c, http.StatusOK, "login.tmpl", gin.H{
			"userCode": userCode,
			"errMsg":   "ユーザーコードとパスワードの組み合わせが一致しません。",
		})
	}

}

func Logout(c *gin.Context) {
	session := sessions.Default(c)
	session.Set("UserID", "")
	session.Clear()
	session.Save()
	RenderHTML(c, http.StatusOK, "login.tmpl", gin.H{
		"errMsg": "ログアウトしました",
	})

}

func saveLoginInfo(c *gin.Context, user model.User) {
	/**
	store := cookie.NewStore([]byte("secret"))
	store.Options(sessions.Options{
		"",
		"",
		0,
		true,
		true,
		http.SameSiteStrictMode,
	})

	s := &sessions.session{"mysession", c.Request, store, nil, false, c.Writer}
	c.Set(sessions.DefaultKey, s)
	**/

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
	log.Println("Count", count)
}
