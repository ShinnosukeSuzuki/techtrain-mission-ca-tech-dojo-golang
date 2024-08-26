package repositories

import (
	"database/sql"

	"github.com/ShinnosukeSuzuki/techtrain-mission-ca-tech-dojo-golang/models"
)

// リポジトリ構造体を定義
type UserRepository struct {
	db *sql.DB
}

// リポジトリのコンストラクタ
func NewUserRepository(db *sql.DB) UserRepository {
	return UserRepository{db: db}
}

// nameとtokenから新規ユーザーを作成する
func (r *UserRepository) CreateUser(name string, token string) (models.User, error) {
	const sqlInsertUser = `
		INSERT INTO users (name, token)
		VALUES (?, ?);
	`

	result, err := r.db.Exec(sqlInsertUser, name, token)
	if err != nil {
		return models.User{}, err
	}
	id, _ := result.LastInsertId()

	var newUser models.User
	newUser.ID = int(id)
	newUser.Name = name
	newUser.Token = token

	return newUser, nil
}

// tokenからユーザーを取得する
func (r *UserRepository) GetUserByToken(token string) (models.User, error) {
	const sqlSelectUserByToken = `
		SELECT id, name, token
		FROM users
		WHERE token = ?;
	`

	var user models.User
	err := r.db.QueryRow(sqlSelectUserByToken, token).Scan(&user.ID, &user.Name, &user.Token)
	if err != nil {
		return models.User{}, err
	}

	return user, nil
}

// tokenが一致するユーザーのnameを更新する
func (r *UserRepository) UpdateUserNameByToken(token string, name string) error {
	const sqlUpdateUserNameByToken = `
		UPDATE users
		SET name = ?
		WHERE token = ?;
	`

	_, err := r.db.Exec(sqlUpdateUserNameByToken, name, token)
	if err != nil {
		return err
	}

	return nil
}
