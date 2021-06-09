package controller

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/tochukaso/golang-study/env"
	"github.com/tochukaso/golang-study/model"
	"golang.org/x/crypto/bcrypt"
)

func ShowLogin(c *gin.Context) {
	RenderHTML(c, http.StatusOK, "login.tmpl", gin.H{})
}

func AttemptLogin(c *gin.Context) {
	userCode := c.PostForm("userCode")
	password := c.PostForm("password")
	isGuest := c.PostForm("isGuest")

	if isGuest == "true" {
		saveGuestLoginInfo(c)
		renderDefaultProductIndexView(c, "")
	} else {

		dbUser := model.GetUserFromCode(userCode)

		if bcrypt.CompareHashAndPassword([]byte(dbUser.Password), []byte(password)) == nil {
			saveLoginInfo(c, dbUser)
			renderDefaultProductIndexView(c, "")
		} else {
			RenderHTML(c, http.StatusOK, "login.tmpl", gin.H{
				"userCode": userCode,
				"errMsg":   "ユーザーコードとパスワードの組み合わせが一致しません。",
			})
		}

	}
}

func Logout(c *gin.Context) {
	session := sessions.Default(c)
	setSessionOption(session)
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
	setSessionOption(session)
	session.Set("UserID", int(user.ID))
	session.Save()
	SessionCheck(c)
}

func saveGuestLoginInfo(c *gin.Context) {
	session := sessions.Default(c)
	setSessionOption(session)
	session.Set("UserID", model.GuestLoginID)
	session.Save()
	SessionCheck(c)
}

func setSessionOption(session sessions.Session) {
	session.Options(MakeSessionOption())
}

func MakeSessionOption() sessions.Options {
	return sessions.Options{
		Path:     "/",
		MaxAge:   0,
		HttpOnly: true,
		SameSite: env.GetSameSitePolicy(),
		Secure:   env.GetCookieSSL(),
	}
}

func SessionCheck(c *gin.Context) {
	session := sessions.Default(c)
	log.Println("session", session)
	userID := session.Get("UserID")
	log.Println("UserID", userID)
	if userID == nil {
		c.Redirect(http.StatusTemporaryRedirect, "/login/")
		c.Abort()
	} else if userID == model.GuestLoginID {
		c.Set("UserID", model.GuestLoginID)
		c.Set("UserCode", model.GuestLoginID)
		c.Set("UserName", "ゲスト")
		c.Next()
	} else {
		user := model.GetUserFromId(strconv.Itoa(userID.(int)))
		c.Set("UserID", userID)
		c.Set("UserCode", user.UserCode)
		c.Set("UserName", user.UserName)
		c.Next()
	}
}
