package services

import "database/sql"

// サービス構造体を定義
type MyAppService struct {
	// フィールドにsql.DB型を持つ
	db *sql.DB
}

// サービスのコンストラクタ
func NewMyAppService(db *sql.DB) *MyAppService {
	return &MyAppService{db: db}
}
