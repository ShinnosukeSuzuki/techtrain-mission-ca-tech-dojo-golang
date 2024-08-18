package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

func NewDB() (*sql.DB, error) {
	// .envファイルから環境変数を取得
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	var (
		dbUser     = os.Getenv("USERNAME")
		dbPassword = os.Getenv("USERPASS")
		dbDatabase = os.Getenv("DATABASE")
		dbConn     = fmt.Sprintf("%s:%s@tcp(127.0.0.1:3306)/%s?parseTime=true", dbUser, dbPassword, dbDatabase)
	)

	db, err := sql.Open("mysql", dbConn)
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}

	return db, nil
}
