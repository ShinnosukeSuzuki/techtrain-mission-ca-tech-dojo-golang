package models

// UserCharacterの構造体
type UserCharacter struct {
	UserCharacterID string
	CharacterID     string
}

// UserCharacterDetailの構造体
type UserCharacterDetail struct {
	UserCharacterID string
	CharacterID     string
	Name            string
}

// CharacterListの構造体
type CharacterList struct {
	Characters []UserCharacterDetail
}
