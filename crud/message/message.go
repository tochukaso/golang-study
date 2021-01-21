package message

import (
	"fmt"

	"github.com/go-playground/validator/v10"
)

func ConvertMessage(e validator.FieldError) string {

	var field string
	switch e.Field() {
	case "Name":
		field = "商品名"
	case "OrgCode":
		field = "商品コード"
	case "JanCode":
		field = "Janコード"
	case "UserName":
		field = "ユーザー名"
	case "UserCode":
		field = "ユーザーコード"
	case "Email":
		field = "メールアドレス"
	case "Password":
		field = "パスワード"
	default:
		field = e.Field()
	}
	var eMsg string
	switch e.Tag() {
	case "required":
		eMsg = field + "は必須です"
	case "ascii":
		eMsg = field + "は半角英数字で入力してください"
	case "gte":
		eMsg = fmt.Sprintf("%vは%v文字以上の長さで入力してください", field, e.Param())
	case "duplicateCode":
		eMsg = fmt.Sprintf("%vは既に登録されています", field)
	default:
		fmt.Println("tagname", e.Tag())
		eMsg = field + "は不正です"
	}
	return eMsg
}
