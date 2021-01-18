package model

// Entity が兼ね備えるべきCRUDのメソッドを定義する。
type Entity interface {
	Create()
	Read()
	Update()
	Delete()
}
