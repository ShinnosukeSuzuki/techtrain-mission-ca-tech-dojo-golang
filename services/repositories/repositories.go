package repositories

import "github.com/ShinnosukeSuzuki/techtrain-mission-ca-tech-dojo-golang/models"

// User関連を引き受けるリポジトリインターフェース
type UserRepository interface {
	CreateUser(name string, token string) (models.User, error)
	GetUserByToken(token string) (models.User, error)
	GetUserById(userId string) (models.User, error)
	UpdateUserName(userId, name string) error
}
