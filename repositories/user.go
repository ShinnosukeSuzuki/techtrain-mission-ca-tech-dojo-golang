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
func (r *UserRepository) Create(id, name, token string) (models.User, error) {
	const query = `
		INSERT INTO users (id, name, token)
		VALUES (?, ?, ?);
	`

	_, err := r.db.Exec(query, id, name, token)
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
func (r *UserRepository) GetByToken(token string) (models.User, error) {
	const query = `
		SELECT id, name, token
		FROM users
		WHERE token = ?;
	`

	var user models.User
	err := r.db.QueryRow(query, token).Scan(&user.ID, &user.Name, &user.Token)
	if err != nil {
		return models.User{}, err
	}

	return user, nil
}

// idからユーザーを取得する
func (r *UserRepository) GetById(userID string) (models.User, error) {
	const query = `
		SELECT id, name, token
		FROM users
		WHERE id = ?;
	`

	var user models.User
	err := r.db.QueryRow(query, userID).Scan(&user.ID, &user.Name, &user.Token)
	if err != nil {
		return models.User{}, err
	}

	return user, nil
}

// userIdが一致するユーザーのnameを更新する
func (r *UserRepository) UpdateName(userID, name string) error {
	const query = `
		UPDATE users
		SET name = ?
		WHERE id = ?;
	`

	_, err := r.db.Exec(query, name, userID)
	if err != nil {
		return err
	}

	return nil
}
