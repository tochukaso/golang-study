package controller

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"omori.jp/model"
)

func InitProduct() {
	model.InitProduct()
}

func ShowProducts(c *gin.Context) {
	name := c.Query("name")
	orgCode := c.Query("orgCode")
	fmt.Println("name", name)
	fmt.Println("orgCode", orgCode)

	products := model.ReadProduct(orgCode, name)
	fmt.Println(products)

	c.HTML(http.StatusOK, "index.tmpl", gin.H{
		"name":     name,
		"orgCode":  orgCode,
		"products": products,
	})

}

func GetProduct(c *gin.Context) {
	id := c.Param("id")
	fmt.Println("id", id)

	product := model.GetProductFromId(id)
	fmt.Println(product)

	c.HTML(http.StatusOK, "detail.tmpl", gin.H{
		"P": product,
	})

}

func PutProduct(c *gin.Context) {
	product := model.Product{}
	err := c.Bind(&product)
	if err != nil {
		c.String(http.StatusBadRequest, "Bad request")
		return
	}

	isFirst := product.ID == 0
	var msg string
	if isFirst {
		model.CreateProduct(&product)
		msg = "登録しました"
	} else {
		model.UpdateProduct(product)
		msg = "保存しました"
	}

	c.HTML(http.StatusOK, "detail.tmpl", gin.H{
		"P":   product,
		"msg": msg,
	})
}

func DeleteProduct(c *gin.Context) {
	id := c.PostForm("id")
	fmt.Println("id", id)

	model.DeleteProduct(id)

	products := model.ReadProduct("", "")
	c.HTML(http.StatusOK, "index.tmpl", gin.H{
		"msg":      "削除しました",
		"products": products,
	})
}
