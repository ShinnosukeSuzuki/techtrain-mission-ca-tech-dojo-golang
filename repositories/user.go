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
func (r *UserRepository) CreateUser(id, name, token string) (models.User, error) {
	const sqlInsertUser = `
		INSERT INTO users (id, name, token)
		VALUES (?, ?, ?);
	`

	_, err := r.db.Exec(sqlInsertUser, id, name, token)
	if err != nil {
		return models.User{}, err
	}

	var newUser models.User
	newUser.ID = id
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

// idからユーザーを取得する
func (r *UserRepository) GetUserById(userId string) (models.User, error) {
	const sqlSelectUserById = `
		SELECT id, name, token
		FROM users
		WHERE id = ?;
	`

	var user models.User
	err := r.db.QueryRow(sqlSelectUserById, userId).Scan(&user.ID, &user.Name, &user.Token)
	if err != nil {
		return models.User{}, err
	}

	return user, nil
}

// userIdが一致するユーザーのnameを更新する
func (r *UserRepository) UpdateUserName(userId, name string) error {
	const sqlUpdateUserNameByToken = `
		UPDATE users
		SET name = ?
		WHERE id = ?;
	`

	_, err := r.db.Exec(sqlUpdateUserNameByToken, name, userId)
	if err != nil {
		return err
	}

	return nil
}
