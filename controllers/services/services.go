package services

import "github.com/ShinnosukeSuzuki/techtrain-mission-ca-tech-dojo-golang/models"

// User関連を引き受けるサービス
type UserServicer interface {
	UserCreateService(name string) (models.User, error)
	UserGetService(userId string) (models.User, error)
	UserUpdateService(userId, name string) error
}
