package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"omori.jp/controller"
	"omori.jp/middleware"
)

func main() {
	engine := gin.Default()
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
	addUserController(engine)
}

func addLoginController(engine *gin.Engine) {
	engine.GET("/", controller.ShowLogin)
	engine.GET("/login/", controller.ShowLogin)
	engine.POST("/login/", controller.AttemptLogin)
}

func addProductController(engine *gin.Engine) {
	controller.InitProduct()
	engine.GET("/product/", controller.ShowProducts)
	engine.GET("/product/detail/:id", controller.GetProduct)

	engine.GET("/product/new", func(c *gin.Context) {
		c.HTML(http.StatusOK, "product_detail.tmpl", gin.H{})
	})
	engine.POST("/product/", controller.PutProduct)
	engine.POST("/product/delete", controller.DeleteProduct)
}

func addUserController(engine *gin.Engine) {
	controller.InitUser()
	engine.GET("/user/", controller.ShowUsers)
	engine.GET("/user/detail/:id", controller.GetUser)

	engine.GET("/user/new", func(c *gin.Context) {
		c.HTML(http.StatusOK, "user_detail.tmpl", gin.H{})
	})
	engine.POST("/user/", controller.PutUser)
	engine.POST("/user/delete", controller.DeleteUser)
}
