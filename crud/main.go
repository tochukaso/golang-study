package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"omori.jp/product"
)

func main() {
	product.Init()
	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	r.LoadHTMLGlob("templates/*")
	r.GET("/index", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.tmpl", gin.H{})
	})

	r.GET("/product", func(c *gin.Context) {
		name := c.Query("name")
		orgCode := c.Query("orgCode")
		fmt.Println("name", name)
		fmt.Println("orgCode", orgCode)

		products := product.Read(orgCode, name)
		fmt.Println(products)

		c.HTML(http.StatusOK, "index.tmpl", gin.H{
			"name":     name,
			"orgCode":  orgCode,
			"products": products,
		})
	})

	r.GET("/product/new", func(c *gin.Context) {
		c.HTML(http.StatusOK, "detail.tmpl", gin.H{})
	})

	r.POST("/product", func(c *gin.Context) {
		name := c.PostForm("name")
		orgCode := c.PostForm("orgCode")
		fmt.Println("name", name)
		fmt.Println("orgCode", orgCode)

		product.Create(&product.Product{Name: name, OrgCode: orgCode})

		c.HTML(http.StatusOK, "detail.tmpl", gin.H{
			"name":    name,
			"orgCode": orgCode,
		})
	})

	loadTemplates(r)
	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}

func loadTemplates(r *gin.Engine) {
	r.Static("/static", "static")
}
