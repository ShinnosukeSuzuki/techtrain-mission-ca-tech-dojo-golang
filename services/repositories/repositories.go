package repositories

import "github.com/ShinnosukeSuzuki/techtrain-mission-ca-tech-dojo-golang/models"

// User関連を引き受けるリポジトリインターフェース
type UserRepository interface {
	Create(name string, token string) (models.User, error)
	GetByToken(token string) (models.User, error)
	GetById(userId string) (models.User, error)
	UpdateName(userId, name string) error
}

// Character関連を引き受けるリポジトリインターフェース
type CharacterRepository interface {
	GetAllList() ([]models.Character, error)
}

// UserCharacter関連を引き受けるリポジトリインターフェース
type UserCharacterRepository interface {
	GetList(userId string) ([]models.UserCharacter, error)
	InsertBulk(userId string, characters []models.UserCharacterInsert) error
}
