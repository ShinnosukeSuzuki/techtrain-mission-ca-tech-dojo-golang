package db

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

func NewDB() (*sql.DB, error) {
	// DOCKER_ENV 環境変数が設定されていない場合は .env ファイルを読み込む(ローカル環境では .env ファイルを読み込むが、Docker環境ではコンテナに設定された環境変数を読み込む)
	if os.Getenv("DOCKER_ENV") == "" {
		err := godotenv.Load()
		if err != nil {
			return nil, fmt.Errorf("error loading .env file: %v", err)
		}
	}

	var (
		dbUser     = os.Getenv("USERNAME")
		dbPassword = os.Getenv("USERPASS")
		dbDatabase = os.Getenv("DATABASE")
		dbHost     = os.Getenv("DBHOST")
		dbPort     = os.Getenv("DBPORT")
		dbConn     = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", dbUser, dbPassword, dbHost, dbPort, dbDatabase)
	)

	db, err := sql.Open("mysql", dbConn)

	if err != nil {
		return nil, err
	}

	return db, nil
}
