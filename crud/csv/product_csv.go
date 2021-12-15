package csv

import (
	"time"

	"github.com/tochukaso/golang-study/model"
)

type ProductCSV struct {
	ID            uint
	ProductName   string    `csv:"商品名"`
	OrgCode       string    `csv:"商品コード"`
	JanCode       string    `csv:"Janコード"`
	ProductDetail string    `csv:"商品説明"`
	ProductPrice  int       `csv:"商品価格"`
	Rating        int       `csv:"レーティング"`
	Review        int       `csv:"レビュー"`
	CreatedAt     time.Time `csv:"登録日時"`
	UpdatedAt     time.Time `csv:"更新日時"`
}

func ConvertProductCSV(p model.Product) ProductCSV {
	return ProductCSV{
		p.ID,
		p.ProductName,
		p.OrgCode,
		p.JanCode,
		p.ProductDetail,
		p.ProductPrice,
		p.Rating,
		p.Review,
		p.CreatedAt,
		p.UpdatedAt,
	}

}
