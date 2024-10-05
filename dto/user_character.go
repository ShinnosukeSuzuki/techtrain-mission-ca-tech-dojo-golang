package dto

// UserCharacterの構造体
type UserCharacter struct {
	ID          string `db:"id"`
	CharacterID string `db:"character_id"`
}
