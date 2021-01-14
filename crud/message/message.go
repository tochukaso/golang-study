package message

import "github.com/go-playground/validator/v10"

func ConvertMessage(e validator.FieldError) string {

	var field string
	switch e.Field() {
	case "Name":
		field = "商品名"
	case "OrgCode":
		field = "商品コード"
	case "JanCode":
		field = "Janコード"
	default:
		field = e.Field()
	}
	var eMsg string
	switch e.Tag() {
	case "required":
		eMsg = field + "は必須です"
	case "ascii":
		eMsg = field + "は半角英数字で入力してください"
	default:
		eMsg = field + "は不正です"
	}
	return eMsg
}
