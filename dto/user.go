package dto

// ユーザーの構造体
type User struct {
	ID    string `db:"id"`
	Name  string `db:"name"`
	Token string `db:"token"`
}
