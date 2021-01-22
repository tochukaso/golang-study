package main

import (
	"io"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	csrf "github.com/utrack/gin-csrf"
	"omori.jp/controller"
	"omori.jp/env"
	"omori.jp/middleware"
	"omori.jp/model"
)

func main() {
	setLogger()
	engine := gin.Default()
	store := cookie.NewStore([]byte("secret"))
	/** comment out Cookie options
	store.Options(sessions.Options{
		"",
		"",
		0,
		false,
		true,
		http.SameSiteStrictMode,
	})
	**/
	engine.Use(sessions.Sessions("mysession", store))
	engine.Use(csrf.Middleware(csrf.Options{
		Secret: "token",
		ErrorFunc: func(c *gin.Context) {
			controller.RenderHTML(c, http.StatusForbidden, "token_error.tmpl", gin.H{})
			c.Abort()
		},
	}))

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
	addMailController(engine)
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
	group.POST("/upload", controller.UploadProduct)
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

func addMailController(engine *gin.Engine) {
	controller.InitMailTemplate()
	group := engine.Group("/mail")
	group.Use(sessionCheck)
	group.GET("/", controller.ShowMailTemplates)
	group.GET("/detail/:code", controller.GetMailTemplate)

	group.POST("/", controller.PutMailTemplate)
}

func sessionCheck(c *gin.Context) {
	session := sessions.Default(c)
	log.Println("session", session)
	userID := session.Get("UserID")
	log.Println("UserID", userID)
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

func setLogger() {
	logFilePath := env.GetEnv().LogFilePath
	f, _ := os.Create(logFilePath)

	//gin.DefaultWriter = io.MultiWriter(f)
	// Use the following code if you need to write the logs to file and console at the same time.
	gin.DefaultWriter = io.MultiWriter(f, os.Stdout)

	log.SetOutput(f)
}
