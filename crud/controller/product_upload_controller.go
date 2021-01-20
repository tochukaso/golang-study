package controller

import (
	"encoding/csv"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"omori.jp/model"
)

func UploadProduct(c *gin.Context) {

	filePath, err := c.FormFile("productFile")
	if err != nil {
		fmt.Println("err", err)
		RenderHTML(c, http.StatusOK, "product_upload.tmpl", gin.H{
			"errMsg": "CSVファイルが設定されていません。",
		})
		return
	}

	file, _ := filePath.Open()
	defer file.Close()
	errors, i, updateCount := processCsv(file)
	RenderHTML(c, http.StatusOK, "product_upload.tmpl", gin.H{
		"msg":      "取り込み処理を完了しました。",
		"totalRow": strconv.Itoa(i),
		"saveRow":  updateCount,
		"errMsg":   errors,
	})
}

func processCsv(file multipart.File) ([]string, int, int) {
	r := csv.NewReader(file)
	r.TrailingComma = true
	var errors []string
	var i int
	var updateCount int
	// 先頭行は読み込みを飛ばす
	r.Read()
	for {
		i++

		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			errors = append(errors, fmt.Sprintf("%v行目でエラーが発生しました。取り込みをスキップします", i))
			continue
		}

		if len(record) <= 5 {
			errors = append(errors, fmt.Sprintf("%v行目のCSVのフォーマットが不正です。取り込みをスキップします", i))
			continue
		}

		product, convErr := convProduct(record)
		if convErr != "" {
			errors = append(errors, fmt.Sprintf("%v行目で、%v。取り込みをスキップします", i, convErr))
			continue
		}

		errs := validateProduct(*product)
		if errs != nil {
			for _, e := range errs {
				errors = append(errors, fmt.Sprintf("%v行目で、%v。取り込みをスキップします", i, e))
			}
			continue
		}
		_, saveErr := saveProduct(*product)
		if saveErr != nil {
			errors = append(errors, fmt.Sprintf("%v行目で、データの更新時にエラーが発生しました。取り込みをスキップします。", i))
			continue
		}
		updateCount++
	}

	return errors, i, updateCount
}

func convProduct(record []string) (*model.Product, string) {
	var product model.Product
	id, err := convUint(record[0])
	if err != nil {
		return nil, fmt.Sprintf("IDに不正な値(%v)が設定されています", record[1])
	}
	product.ID = id
	product.Name = record[1]
	product.OrgCode = record[2]
	product.JanCode = record[3]
	product.Detail = record[4]
	return &product, ""
}

func convUint(s string) (uint, error) {
	i, err := strconv.Atoi(s)
	if err != nil {
		return 0, err
	}

	return uint(i), nil
}
