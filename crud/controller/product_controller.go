package controller

import (
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/gocarina/gocsv"
	"omori.jp/csv"
	"omori.jp/env"
	"omori.jp/mail"
	"omori.jp/message"
	"omori.jp/model"
	"omori.jp/pagination"
)

type productSearchQuery struct {
	productName string
	orgCode     string
	page        int
	pageSize    int
}

const (
	productDefaultPage     int = 1
	productDefaultPageSize int = 10
)

func InitProduct() {
	model.InitProduct()
}

func ShowProducts(c *gin.Context) {
	renderProductIndexView(c, extractProductSearchQuery(c), "")
}

func ShowProductsJSON(c *gin.Context) {
	products, count := searchProduct(extractProductSearchQuery(c))
	ResponseJSON(c, http.StatusOK, map[string]interface{}{
		"count":    count,
		"products": convertProductsJSON(products),
	})
}

func GetProduct(c *gin.Context) {
	id := c.Param("id")
	log.Println("id", id)

	product := model.GetProductFromId(id)
	log.Println(product)

	RenderHTML(c, http.StatusOK, "product_detail.tmpl", gin.H{
		"P": product,
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
			log.Println("画像アップロードエラー", err)
		}
	}

	RenderHTML(c, http.StatusOK, "product_detail.tmpl", gin.H{
		"P":   product,
		"msg": msg,
	})

}

func DeleteProduct(c *gin.Context) {
	id, _ := strconv.Atoi(c.PostForm("id"))
	product := createIDProduct(id)

	product.Delete()
	renderDefaultProductIndexView(c, "削除しました")
}

func DownloadProduct(c *gin.Context) {

	productName := c.Query("productName")
	orgCode := c.Query("orgCode")
	products, _ := model.ReadProduct(orgCode, productName)

	csvStr, err := gocsv.MarshalString(convertProduct(products))
	log.Println("csvStr", csvStr)
	if err == nil {
		renderProductIndexView(c, extractProductSearchQuery(c), "ダウンロードに失敗しました")
		c.Abort()
		return
	}

	header := c.Writer.Header()
	header["Content-type"] = []string{"text/csv"}
	header["Content-Disposition"] = []string{"attachment; filename= products.csv"}
	header["Content-Length"] = []string{""}

	log.Println("products", products)
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
	log.Println("dbProduct", dbProduct)

	if dbProduct.GetID() != 0 {
		sl.ReportError(product.OrgCode, "OrgCode", "OrgCode", "duplicateCode", "")
	}
}

func convertProduct(products []model.Product) []csv.ProductCSV {
	var result []csv.ProductCSV
	for i, p := range products {
		result = append(result, csv.ConvertProductCSV(p))
		log.Println("r", result[i])
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

func searchProduct(query productSearchQuery) ([]model.Product, int) {
	return model.ReadProductWithPaging(
		query.page,
		query.pageSize,
		query.orgCode,
		query.productName,
	)
}

func renderProductIndexView(c *gin.Context, query productSearchQuery, msg string) {
	products, count := searchProduct(query)
	RenderHTML(c, http.StatusOK, "product_index.tmpl", gin.H{
		"msg":         msg,
		"productName": query.productName,
		"orgCode":     query.orgCode,
		"page":        query.page,
		"count":       count,
		"pageSize":    query.pageSize,
		"products":    products,
		"pagination":  pagination.Pagination(count, query.page, query.pageSize),
	})
}

func renderDefaultProductIndexView(c *gin.Context, msg string) {
	renderProductIndexView(c, createDefaultProductSearchQuery(), msg)
}

func createDefaultProductSearchQuery() productSearchQuery {
	return productSearchQuery{"", "", productDefaultPage, productDefaultPageSize}
}

func extractProductSearchQuery(c *gin.Context) productSearchQuery {
	page, _ := strconv.Atoi(c.DefaultQuery("page", strconv.Itoa(productDefaultPage)))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", strconv.Itoa(productDefaultPageSize)))
	return productSearchQuery{
		c.Query("productName"),
		c.Query("orgCode"),
		page,
		pageSize,
	}
}

func convertProductsJSON(products []model.Product) []model.ProductJSON {
	var result []model.ProductJSON

	for _, p := range products {
		result = append(result, convertProductJSON(p))
	}

	return result
}

func convertProductJSON(p model.Product) model.ProductJSON {
	return model.ProductJSON{
		p,
		p.GetAbsoluteImagePath(),
	}
}
