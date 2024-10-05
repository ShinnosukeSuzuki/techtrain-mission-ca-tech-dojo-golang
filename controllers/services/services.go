package services

import (
	"github.com/ShinnosukeSuzuki/techtrain-mission-ca-tech-dojo-golang/models"
)

// User関連を引き受けるサービス
type UserServicer interface {
	Create(name string) (models.User, error)
	Get(userID string) (models.User, error)
	UpdateName(userID, name string) error
}

// GacheDraw関連を引き受けるサービス
type GachaDrawServicer interface {
	Draw(times int, userID string) ([]models.Character, error)
}

// UserCharacter関連を引き受けるサービス
type UserCharacterServicer interface {
	List(userID string) (models.CharacterList, error)
}
