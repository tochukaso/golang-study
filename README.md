# Go lang study

## gin フレームワークの利用について

### リクエストメソッド、URLと処理のマッピング(ルーティング)

main.goで以下の様に記載しています。

これは、リクエストメソッドとURLの組み合わせで、どの関数が処理を担当するかを記載しています。

controller側で実際に行う処理を記載します。

``` go
	engine.GET("/", controller.ShowProducts)
	engine.GET("/product/", controller.ShowProducts)
	engine.GET("/product/detail/:id", controller.GetProduct)

	engine.GET("/product/new", func(c *gin.Context) {
		c.HTML(http.StatusOK, "detail.tmpl", gin.H{})
	})
	engine.POST("/product/", controller.PutProduct)
	engine.POST("/product/delete", controller.DeleteProduct)
```

### リクエストパラメーターのバインド

リクエストで渡されたパラメーターやボディを構造体にバインドする事ができます。

以下の記載では、model.Productの構造体にリクエストをバインドしています。

``` go
	var product model.Product
	err := c.ShouldBind(&product)
```

### レスポンスのハンドリング

以下では、レスポンスのステータスコード、
レスポンスに使用するテンプレート、レスポンスで使用するパラメーターを設定しています。

goのtemplate側では、パラメーターを使用してhtmlファイルを組み立ててレスポンスとして返却します。

``` go
	c.HTML(http.StatusOK, "index.tmpl", gin.H{
		"name":       name,
		"orgCode":    orgCode,
		"page":       page,
		"count":      count,
		"pageSize":   pageSize,
		"products":   products,
		"pagination": pagination.Pagination(count, page, pageSize),
	})
```
