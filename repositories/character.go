package repositories

import (
	"database/sql"

	"github.com/ShinnosukeSuzuki/techtrain-mission-ca-tech-dojo-golang/models"
)

// リポジトリ構造体を定義
type CharacterRepository struct {
	db *sql.DB
}

// リポジトリのコンストラクタ
func NewCharacterRepository(db *sql.DB) CharacterRepository {
	return CharacterRepository{db: db}
}

// キャラクター全件取得
func (r *CharacterRepository) GetAllList() ([]models.Character, error) {
	const query = `SELECT id, name, probability FROM characters;`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var characters []models.Character
	for rows.Next() {
		var character models.Character
		if err := rows.Scan(&character.ID, &character.Name, &character.Probability); err != nil {
			return nil, err
		}
		characters = append(characters, character)
	}

	return characters, nil
}
