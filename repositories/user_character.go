package repositories

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/ShinnosukeSuzuki/techtrain-mission-ca-tech-dojo-golang/dto"
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
func (r *UserCharacterRepository) GetList(userID string) ([]dto.UserCharacter, error) {
	const query = `
		SELECT BIN_TO_UUID(id) as id, BIN_TO_UUID(character_id) as character_id
		FROM users_characters
		WHERE user_id = UUID_TO_BIN(?);
	`

	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var userCharacters []dto.UserCharacter
	for rows.Next() {
		var userCharacter dto.UserCharacter
		if err := rows.Scan(&userCharacter.ID, &userCharacter.CharacterID); err != nil {
			return nil, err
		}
		userCharacters = append(userCharacters, userCharacter)
	}

	return userCharacters, nil
}

// ガチャ結果をusers_charactersテーブルにバルクインサートする
func (r *UserCharacterRepository) InsertBulk(userID string, gachaResults []models.Character) error {
	if len(gachaResults) == 0 {
		return nil
	}

	// クエリのプレースホルダーを生成
	valueStrings := make([]string, 0, len(gachaResults))
	valueArgs := make([]any, 0, len(gachaResults)*3)
	for _, g := range gachaResults {
		// UUIDの生成
		newID, err := uuid.GenerateUUID()
		if err != nil {
			return fmt.Errorf("failed to generate UUID: %w", err)
		}
		valueStrings = append(valueStrings, "(UUID_TO_BIN(?), UUID_TO_BIN(?), UUID_TO_BIN(?))")
		valueArgs = append(valueArgs, newID)
		valueArgs = append(valueArgs, userID)
		valueArgs = append(valueArgs, g.ID)
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
