# Go lang study

## Gin Webフレームワークの利用について

GinはGolangで書かれているWebアプリケーション用のフレームワークです。

_martini-like_ APIを備えており、非常に高速なパフォーマンスを誇ります。

Githubのスター数も多く、現時点で44.8Kほどあります。
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
		"productName":       productName,
		"orgCode":    orgCode,
		"page":       page,
		"count":      count,
		"pageSize":   pageSize,
		"products":   products,
		"pagination": pagination.Pagination(count, page, pageSize),
	})
```

## GORM フレームワークの利用について

ORMマッパーとして使用します。

データベースへの登録や読込などをSQLを記載せずに、実行できます。

構造体を直接保存、読込するような感覚で利用することが出来ます。

また、構造体に従って、DDLを実行することも出来ます。

構造体の定義とDDLの実行は以下の様に行います。

``` go

type Product struct {
	gorm.Model
	ProductName    string `form:"ProductName" binding:"required" validate:"required"`
	OrgCode string `form:"OrgCode" validate:"required,ascii"`
	JanCode string `form:"JanCode" validate:"ascii"`
	ProductDetail  string
}

func InitProduct() {
	db := db.GetDB()
	db.AutoMigrate(&Product{})
}

```

`db(gorm.DB).AutoMigrate`でDDLを実行することが出来ます。

DBへのデータの更新は以下のように構造体を指定して保存するだけです。

``` go
	db.Save(product)
```

ページング処理についてもOffsetとLimitを利用して簡単に実装することが出来ます。

``` go 
	db := db.GetDB()
	offset := (page - 1) * pageSize
	db = db.Offset(offset).Limit(pageSize)
	return readProduct(db, orgCode, productName)
```


