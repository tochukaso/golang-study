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
	engine.LoadHTMLGlob("templates/*")

	addProductController(engine)

	loadTemplates(engine)
	engine.Run(":8081")
}

func loadTemplates(engine *gin.Engine) {
	engine.Static("/static", "static")
}

func addProductController(engine *gin.Engine) {
	controller.InitProduct()
	engine.GET("/", controller.ShowProducts)
	engine.GET("/product/", controller.ShowProducts)
	engine.GET("/product/detail/:id", controller.GetProduct)

	engine.GET("/product/new", func(c *gin.Context) {
		c.HTML(http.StatusOK, "detail.tmpl", gin.H{})
	})
	engine.POST("/product/", controller.PutProduct)
	engine.POST("/product/delete", controller.DeleteProduct)
}
