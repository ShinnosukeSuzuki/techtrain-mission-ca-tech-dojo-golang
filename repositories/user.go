package repositories

import (
	"database/sql"

	"github.com/ShinnosukeSuzuki/techtrain-mission-ca-tech-dojo-golang/models"
)

// nameとtokenから新規ユーザーを作成する
func CreateUser(db *sql.DB, name string, token string) (models.User, error) {
	const sqlInsertUser = `
		INSERT INTO users (name, token)
		VALUES (?, ?);
	`

	var newUser models.User
	newUser.Name = name
	newUser.Token = token
	result, err := db.Exec(sqlInsertUser, newUser.Name, newUser.Token)
	if err != nil {
		return models.User{}, err
	}
	id, _ := result.LastInsertId()
	newUser.ID = int(id)

	return newUser, nil
}
