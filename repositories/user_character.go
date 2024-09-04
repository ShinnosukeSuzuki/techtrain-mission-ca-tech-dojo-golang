package repositories

import (
	"database/sql"

	"github.com/ShinnosukeSuzuki/techtrain-mission-ca-tech-dojo-golang/models"
)

// リポジトリ構造体を定義
type UserCharacterRepository struct {
	db *sql.DB
}

// リポジトリのコンストラクタ
func NewUserCharacterRepository(db *sql.DB) UserCharacterRepository {
	return UserCharacterRepository{db: db}
}

// userのtokenに一致するキャラクターを取得する
func (r *UserCharacterRepository) GetUserCharacterList(token string) ([]models.UserCharacter, error) {
	const sqlSelectCharacterByUserID = `
		SELECT users_characters.id, users_characters.character_id, characters.name
		FROM users_characters
		JOIN users ON users_characters.user_id = users.id
		JOIN characters ON users_characters.character_id = characters.id
		WHERE users.token = ?;
	`

	rows, err := r.db.Query(sqlSelectCharacterByUserID, token)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var userCharacters []models.UserCharacter
	for rows.Next() {
		var userCharacter models.UserCharacter
		if err := rows.Scan(&userCharacter.UserCharacterID, &userCharacter.CharacterID, &userCharacter.Name); err != nil {
			return nil, err
		}
		userCharacters = append(userCharacters, userCharacter)
	}

	return userCharacters, nil
}
