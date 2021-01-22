package controller

import (
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/gocarina/gocsv"
	csrf "github.com/utrack/gin-csrf"
	"omori.jp/csv"
	"omori.jp/env"
	"omori.jp/mail"
	"omori.jp/message"
	"omori.jp/model"
	"omori.jp/pagination"
)

func InitProduct() {
	model.InitProduct()
}

func ShowProducts(c *gin.Context) {
	productName := c.Query("productName")
	orgCode := c.Query("orgCode")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))

	products, count := model.ReadProductWithPaging(page, pageSize, orgCode, productName)
	RenderHTML(c, http.StatusOK, "product_index.tmpl", gin.H{
		"productName": productName,
		"orgCode":     orgCode,
		"page":        page,
		"count":       count,
		"pageSize":    pageSize,
		"products":    products,
		"pagination":  pagination.Pagination(count, page, pageSize),
	})

}

func GetProduct(c *gin.Context) {
	id := c.Param("id")
	fmt.Println("id", id)

	product := model.GetProductFromId(id)
	fmt.Println(product)

	RenderHTML(c, http.StatusOK, "product_detail.tmpl", gin.H{
		"_csrf": csrf.GetToken(c),
		"P":     product,
	})
}

func PutProduct(c *gin.Context) {
	var product model.Product
	err := c.ShouldBind(&product)
	errors := validateProduct(product)
	if err != nil || errors != nil {
		RenderHTML(c, http.StatusOK, "product_detail.tmpl", gin.H{
			"P":      product,
			"errMsg": errors,
		})
		return
	}

	file, fileErr := c.FormFile("productImage")
	if fileErr == nil {
		product.ProductImage = file.Filename
	}

	msg, _ := saveProduct(product)
	if product.ID == 0 {
		product.ID = model.GetProductFromCode(product.OrgCode).GetID()
	}
	if fileErr == nil {
		dirPath := env.GetEnv().ProductImagePath + "/" + strconv.Itoa(int(product.ID))
		os.Mkdir(dirPath, 0755)
		err = c.SaveUploadedFile(file, dirPath+"/"+file.Filename)
		if err != nil {
			fmt.Println("画像アップロードエラー", err)
		}
	}

	RenderHTML(c, http.StatusOK, "product_detail.tmpl", gin.H{
		"_csrf": csrf.GetToken(c),
		"P":     product,
		"msg":   msg,
	})

}

func DeleteProduct(c *gin.Context) {
	id, _ := strconv.Atoi(c.PostForm("id"))
	product := createIDProduct(id)

	product.Delete()
	fmt.Println("id", id)

	products, count := model.ReadProduct("", "")
	RenderHTML(c, http.StatusOK, "product_index.tmpl", gin.H{
		"msg":      "削除しました",
		"products": products,
		"count":    count,
	})
}

func DownloadProduct(c *gin.Context) {

	productName := c.Query("productName")
	orgCode := c.Query("orgCode")
	products, count := model.ReadProduct(orgCode, productName)

	header := c.Writer.Header()
	header["Content-type"] = []string{"text/csv"}
	header["Content-Disposition"] = []string{"attachment; filename= products.csv"}
	header["Content-Length"] = []string{""}

	fmt.Println("products", products)

	csvStr, err := gocsv.MarshalString(convertProduct(products))
	fmt.Println("csvStr", csvStr)
	if err != nil {
		c.Status(http.StatusBadRequest)
		RenderHTML(c, http.StatusOK, "product_index.tmpl", gin.H{
			"productName": productName,
			"msg":         "ダウンロードに失敗しました",
			"orgCode":     orgCode,
			"products":    products,
			"pagination":  pagination.Pagination(count, 1, 10),
		})
		return
	}
	writer := c.Writer
	writer.Write([]byte(csvStr))
	writer.Flush()

	c.Status(http.StatusOK)
}

func createIDProduct(id int) model.Product {
	var product model.Product
	product.ID = uint(id)
	return product
}

func checkDuplicateOrgCode(sl validator.StructLevel) {
	product := sl.Current().Interface().(model.Product)

	if product.OrgCode == "" || product.ID != 0 {
		return
	}
	dbProduct := model.GetProductFromCode(product.OrgCode)
	fmt.Println("dbProduct", dbProduct)

	if dbProduct.GetID() != 0 {
		sl.ReportError(product.OrgCode, "OrgCode", "OrgCode", "duplicateCode", "")
	}
}

func convertProduct(products []model.Product) []csv.ProductCSV {
	var result []csv.ProductCSV
	for i, p := range products {
		result = append(result, csv.ProductCSV{p.ID, p.ProductName, p.OrgCode, p.JanCode, p.ProductDetail, p.CreatedAt, p.UpdatedAt})
		fmt.Println("r", result[i])
	}
	return result
}

func validateProduct(product model.Product) []string {
	validate := validator.New()
	validate.RegisterStructValidation(checkDuplicateOrgCode, model.Product{})
	errors := validate.Struct(product)
	if errors == nil {
		return nil
	}
	errs := errors.(validator.ValidationErrors)
	sliceErrs := []string{}
	for _, e := range errs {
		sliceErrs = append(sliceErrs, message.ConvertMessage(e))
	}
	return sliceErrs
}

func saveProduct(product model.Product) (string, error) {
	isFirst := product.ID == 0
	var msg string
	var err error
	if isFirst {
		err = product.Create()
		if err != nil {
			msg = "商品の登録に失敗しました"
		} else {
			msg = "登録しました"
		}
		mail.SendProductRegisterMail(product)
	} else {
		dbProduct := product.Read().(model.Product)
		product.CreatedAt = dbProduct.CreatedAt
		if product.ProductImage == "" {
			product.ProductImage = dbProduct.ProductImage
		}
		product.Update()
		msg = "保存しました"
	}

	return msg, err
}
