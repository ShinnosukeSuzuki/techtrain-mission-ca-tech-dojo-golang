package services

import "github.com/ShinnosukeSuzuki/techtrain-mission-ca-tech-dojo-golang/models"

// User関連を引き受けるサービス
type UserServicer interface {
	Create(name string) (models.User, error)
	Get(userId string) (models.User, error)
	UpdateName(userId, name string) error
}

// UserCharacter関連を引き受けるサービス
type UserCharacterServicer interface {
	UserCharacterGetService(token string) (models.CharacterList, error)
}
