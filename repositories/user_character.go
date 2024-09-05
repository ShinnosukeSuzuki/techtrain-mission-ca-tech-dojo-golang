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

// userのidに一致するキャラクターを取得する
func (r *UserCharacterRepository) GetList(userId string) ([]models.UserCharacter, error) {
	const sqlSelectCharacterByUserID = `
		SELECT uc.id, uc.character_id, c.name
		FROM users_characters as uc
		JOIN characters as c ON uc.character_id = c.id
		WHERE uc.user_id = ?;
	`

	rows, err := r.db.Query(sqlSelectCharacterByUserID, userId)
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
