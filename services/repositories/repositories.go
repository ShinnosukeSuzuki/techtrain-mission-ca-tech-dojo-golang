package repositories

	"github.com/ShinnosukeSuzuki/techtrain-mission-ca-tech-dojo-golang/dto"

// User関連を引き受けるリポジトリインターフェース
type UserRepository interface {
	Create(name string, token string) (dto.User, error)
	GetByToken(token string) (dto.User, error)
	GetById(userID string) (dto.User, error)
	UpdateName(userID, name string) error
}

// Character関連を引き受けるリポジトリインターフェース
type CharacterRepository interface {
	GetAllList() ([]models.Character, error)
}

// UserCharacter関連を引き受けるリポジトリインターフェース
type UserCharacterRepository interface {
	GetList(userID string) ([]models.UserCharacter, error)
	InsertBulk(userID string, characters []models.GachaResult) error
}
