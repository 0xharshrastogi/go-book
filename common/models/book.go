package models

type BookInfo struct {
	Name  string `json:"name"`
	Pages uint   `json:"pages"`
	Code  string `json:"code"`
}
