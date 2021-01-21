package csv

import "time"

type ProductCSV struct {
	ID            uint
	ProductName   string    `csv:"商品名"`
	OrgCode       string    `csv:"商品コード"`
	JanCode       string    `csv:"Janコード"`
	ProductDetail string    `csv:"商品説明"`
	CreatedAt     time.Time `csv:"登録日時"`
	UpdatedAt     time.Time `csv:"更新日時"`
}
