package services

import "github.com/ShinnosukeSuzuki/techtrain-mission-ca-tech-dojo-golang/models"

// User関連を引き受けるサービス
type UserServicer interface {
	UserCreateService(name string) (models.User, error)
	UserGetService(token string) (models.User, error)
	UserUpdateService(token string, name string) error
}

// UserCharacter関連を引き受けるサービス
type UserCharacterServicer interface {
	UserCharacterGetService(token string) (models.CharacterList, error)
}
