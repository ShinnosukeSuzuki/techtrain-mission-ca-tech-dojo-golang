package models

// UserCharacterの構造体
type UserCharacter struct {
	UserCharacterID string
	CharacterID     string
	Name            string
}

// CharacterListの構造体
type CharacterList struct {
	Characters []UserCharacter
}
