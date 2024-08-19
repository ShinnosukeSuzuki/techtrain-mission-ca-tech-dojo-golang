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

// tokenからユーザーを取得する
func GetUserByToken(db *sql.DB, token string) (models.User, error) {
	const sqlSelectUserByToken = `
		SELECT *
		FROM users
		WHERE token = ?;
	`

	var user models.User
	err := db.QueryRow(sqlSelectUserByToken, token).Scan(&user.ID, &user.Name, &user.Token)
	if err != nil {
		return models.User{}, err
	}

	return user, nil
}

// tokenが一致するユーザーのnameを更新する
func UpdateUserNameByToken(db *sql.DB, token string, name string) error {
	const sqlUpdateUserNameByToken = `
		UPDATE users
		SET name = ?
		WHERE token = ?;
	`

	_, err := db.Exec(sqlUpdateUserNameByToken, name, token)
	if err != nil {
		return err
	}

	return nil
}
