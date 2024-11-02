package repositories

import (
	"database/sql"

	"github.com/ShinnosukeSuzuki/techtrain-mission-ca-tech-dojo-golang/dto"
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
func (r *UserRepository) Create(id, name, token string) (dto.User, error) {
	const query = `
		INSERT INTO users (id, name, token)
		VALUES (UUID_TO_BIN(?), ?, UUID_TO_BIN(?));
	`

	_, err := r.db.Exec(query, id, name, token)
	if err != nil {
		return dto.User{}, err
	}

	var newUser dto.User
	newUser.ID = id
	newUser.Name = name
	newUser.Token = token

	return newUser, nil
}

// tokenからユーザーを取得する
func (r *UserRepository) GetByToken(token string) (dto.User, error) {
	const query = `
		SELECT BIN_TO_UUID(id) as id, name, BIN_TO_UUID(token) as token
		FROM users
		WHERE token = UUID_TO_BIN(?);
	`

	var user dto.User
	err := r.db.QueryRow(query, token).Scan(&user.ID, &user.Name, &user.Token)
	if err != nil {
		return dto.User{}, err
	}

	return user, nil
}

// idからユーザーを取得する
func (r *UserRepository) GetById(userID string) (dto.User, error) {
	const query = `
		SELECT BIN_TO_UUID(id) as id, name, BIN_TO_UUID(token) as token
		FROM users
		WHERE id = UUID_TO_BIN(?);
	`

	var user dto.User
	err := r.db.QueryRow(query, userID).Scan(&user.ID, &user.Name, &user.Token)
	if err != nil {
		return dto.User{}, err
	}

	return user, nil
}

// userIdが一致するユーザーのnameを更新する
func (r *UserRepository) UpdateName(userID, name string) error {
	const query = `
		UPDATE users
		SET name = ?
		WHERE id = UUID_TO_BIN(?);
	`

	_, err := r.db.Exec(query, name, userID)
	if err != nil {
		return err
	}

	return nil
}
