package controller

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"omori.jp/message"
	"omori.jp/model"
	"omori.jp/pagination"
)

func InitProduct() {
	model.InitProduct()
}

func ShowProducts(c *gin.Context) {
	name := c.Query("name")
	orgCode := c.Query("orgCode")
	fmt.Println("name", name)
	fmt.Println("orgCode", orgCode)
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))
	fmt.Println("page", page)
	fmt.Println("pageSize", pageSize)

	products, count := model.ReadProductWithPaging(page, pageSize, orgCode, name)
	fmt.Println(products)
	fmt.Println(count)

	c.HTML(http.StatusOK, "product_index.tmpl", gin.H{
		"name":       name,
		"orgCode":    orgCode,
		"page":       page,
		"count":      count,
		"pageSize":   pageSize,
		"products":   products,
		"pagination": pagination.Pagination(count, page, pageSize),
	})

}

func GetProduct(c *gin.Context) {
	id := c.Param("id")
	fmt.Println("id", id)

	product := model.GetProductFromId(id)
	fmt.Println(product)

	c.HTML(http.StatusOK, "product_detail.tmpl", gin.H{
		"P": product,
	})
}

func PutProduct(c *gin.Context) {
	var product model.Product
	err := c.ShouldBind(&product)
	validate := validator.New()
	validate.RegisterValidation("duplicateCode", checkDuplicateOrgCode)
	errors := validate.Struct(product)
	if err != nil || errors != nil {
		errs := errors.(validator.ValidationErrors)
		sliceErrs := []string{}
		for _, e := range errs {
			sliceErrs = append(sliceErrs, message.ConvertMessage(e))
		}
		c.HTML(http.StatusOK, "product_detail.tmpl", gin.H{
			"P":      product,
			"errMsg": sliceErrs,
		})
		return
	}

	isFirst := product.ID == 0
	var msg string
	if isFirst {
		err := product.Create()
		if err != nil {
			c.HTML(http.StatusOK, "product_detail.tmpl", gin.H{
				"P":      product,
				"errMsg": "商品の登録に失敗しました",
			})
			return
		}
		msg = "登録しました"
	} else {
		product.Update()
		msg = "保存しました"
	}

	c.HTML(http.StatusOK, "product_detail.tmpl", gin.H{
		"P":   product,
		"msg": msg,
	})
}

func DeleteProduct(c *gin.Context) {
	id, _ := strconv.Atoi(c.PostForm("id"))
	product := createIDProduct(id)

	product.Delete()
	fmt.Println("id", id)

	products, count := model.ReadProduct("", "")
	c.HTML(http.StatusOK, "product_index.tmpl", gin.H{
		"msg":      "削除しました",
		"products": products,
		"count":    count,
	})
}

func createIDProduct(id int) model.Product {
	var product model.Product
	product.ID = uint(id)
	return product
}

func checkDuplicateOrgCode(fl validator.FieldLevel) bool {
	var data model.Product
	data.OrgCode = fl.Field().String()
	user := data.GetFromCode()
	fmt.Println("checkDuplicate", user)

	if user.GetID() != 0 {
		return false
	}
	return true
}
