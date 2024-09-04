package services

import (
	"github.com/ShinnosukeSuzuki/techtrain-mission-ca-tech-dojo-golang/models"
	"github.com/ShinnosukeSuzuki/techtrain-mission-ca-tech-dojo-golang/repositories"
)

// サービス構造体を定義
type UserCharacterService struct {
	// userRepositoryを埋め込む
	repository repositories.UserCharacterRepository
}

// サービスのコンストラクタ
func NewUserCharacterService(r repositories.UserCharacterRepository) *UserCharacterService {
	return &UserCharacterService{repository: r}
}

// ハンドラー UserCharacterGetHandler 用のサービスメソッド
func (s *UserCharacterService) UserCharacterGetService(token string) (models.CharacterList, error) {
	userCharacters, err := s.repository.GetUserCharacterList(token)
	if err != nil {
		return models.CharacterList{}, err
	}

	characterList := models.CharacterList{Characters: userCharacters}

	return characterList, nil
}
