package models

// UserCharacterの構造体
type UserCharacter struct {
	UserCharacterID string `json:"userCharacterID"`
	CharacterID     string `json:"characterID"`
	Name            string `json:"name"`
}

// CharacterListの構造体
type CharacterList struct {
	Characters []UserCharacter `json:"characters"`
}

// users_charactersテーブルにインサートするための構造体
type UserCharacterInsert struct {
	ID          string `json:"id"`
	CharacterID string `json:"characterID"`
}
