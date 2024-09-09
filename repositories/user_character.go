package repositories

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/ShinnosukeSuzuki/techtrain-mission-ca-tech-dojo-golang/models"
	"github.com/ShinnosukeSuzuki/techtrain-mission-ca-tech-dojo-golang/pkg/uuid"
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
	const query = `
		SELECT uc.id, uc.character_id, c.name
		FROM users_characters as uc
		JOIN characters as c ON uc.character_id = c.id
		WHERE uc.user_id = ?;
	`

	rows, err := r.db.Query(query, userId)
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

// ガチャ結果をusers_charactersテーブルにバルクインサートする
func (r *UserCharacterRepository) InsertBulk(userId string, gachaResults []models.GachaResult) error {
	if len(gachaResults) == 0 {
		return nil
	}

	// クエリのプレースホルダーを生成
	valueStrings := make([]string, 0, len(gachaResults))
	valueArgs := make([]any, 0, len(gachaResults)*3)
	for _, g := range gachaResults {
		valueStrings = append(valueStrings, "(?, ?, ?)")
		valueArgs = append(valueArgs, uuid.GenerateUUID())
		valueArgs = append(valueArgs, userId)
		valueArgs = append(valueArgs, g.CharacterID)
	}

	// クエリ文字列を生成
	query := fmt.Sprintf("INSERT INTO users_characters (id, user_id, character_id) VALUES %s",
		strings.Join(valueStrings, ", "))

	_, err := r.db.Exec(query, valueArgs...)
	if err != nil {
		return err
	}
	return nil
}
