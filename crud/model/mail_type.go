package model

type MailType int

const (
	UserRegister MailType = iota + 1
	ProductRegister
)

func ListMailType() []MailType {
	return []MailType{
		UserRegister,
		ProductRegister,
	}
}

func GetMailType(code MailType) []string {
	var result []string
	switch code {
	case UserRegister:
		result = []string{"UserRegister", "会員登録"}
	case ProductRegister:
		result = []string{"ProductRegister", "商品登録"}
	}

	return result
}

func ListTemplateValiable(code MailType) [][]string {
	var result [][]string
	switch code {
	case UserRegister:
		result = [][]string{
			{"{{.UserCode}}", "ユーザーコード"},
			{"{{.UserName}}", "ユーザー名"},
			{"{{.Email}}", "メールアドレス"},
		}

	case ProductRegister:
		result = [][]string{
			{"{{.ProductName}}", "商品名"},
			{"{{.OrgCode}}", "商品コード"},
			{"{{.JanCode}}", "Janコード"},
			{"{{.ProductDetail}}", "商品詳細"},
		}
	}

	return result
}
