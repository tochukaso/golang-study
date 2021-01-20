package main

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"omori.jp/controller"
	"omori.jp/middleware"
	"omori.jp/model"
)

func main() {
	engine := gin.Default()
	store := cookie.NewStore([]byte("secret"))
	engine.Use(sessions.Sessions("mysession", store))
	engine.Use(middleware.RecordUaAndTime)
	engine.LoadHTMLGlob("templates/*.tmpl")

	addControllers(engine)

	loadTemplates(engine)
	engine.Run(":8081")
}

func loadTemplates(engine *gin.Engine) {
	engine.Static("/static", "static")
}

func addControllers(engine *gin.Engine) {
	addLoginController(engine)
	addProductController(engine)
	addProductUploadController(engine)
	addUserController(engine)
}

func addLoginController(engine *gin.Engine) {
	engine.GET("/", controller.ShowLogin)
	engine.GET("/login/", controller.ShowLogin)
	engine.POST("/login/", controller.AttemptLogin)
	engine.GET("/logout/", controller.Logout)
}

func addProductController(engine *gin.Engine) {
	controller.InitProduct()
	group := engine.Group("/product")
	group.Use(sessionCheck)
	group.GET("/", controller.ShowProducts)
	group.GET("/detail/:id", controller.GetProduct)
	group.GET("/download", controller.DownloadProduct)

	group.GET("/new", func(c *gin.Context) {
		controller.RenderHTML(c, http.StatusOK, "product_detail.tmpl", gin.H{})
	})
	group.POST("/", controller.PutProduct)
	group.POST("/delete", controller.DeleteProduct)
}

func addProductUploadController(engine *gin.Engine) {
	controller.InitProduct()
	group := engine.Group("/product")
	group.Use(sessionCheck)
	group.GET("/upload", func(c *gin.Context) {
		controller.RenderHTML(c, http.StatusOK, "product_upload.tmpl", gin.H{})
	})
}

func addUserController(engine *gin.Engine) {
	controller.InitUser()
	group := engine.Group("/user")
	group.Use(sessionCheck)
	group.GET("/", controller.ShowUsers)
	group.GET("/detail/:id", controller.GetUser)

	group.GET("/new", func(c *gin.Context) {
		controller.RenderHTML(c, http.StatusOK, "user_detail.tmpl", gin.H{})
	})
	group.POST("/", controller.PutUser)
	group.POST("/delete", controller.DeleteUser)
}

func sessionCheck(c *gin.Context) {
	session := sessions.Default(c)
	fmt.Println("session", session)
	userID := session.Get("UserID")
	fmt.Println("UserID", userID)
	if userID == nil {
		c.Redirect(http.StatusTemporaryRedirect, "/login/")
		c.Abort()
	} else {
		user := model.GetUserFromId(strconv.Itoa(userID.(int)))
		c.Set("UserID", userID)
		c.Set("UserCode", user.UserCode)
		c.Set("UserName", user.UserName)
		c.Next()
	}
}
