package models

// ユーザーの構造体
type User struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Token string `json:"token"`
}
