package services

import (
	"github.com/ShinnosukeSuzuki/techtrain-mission-ca-tech-dojo-golang/models"
	"github.com/ShinnosukeSuzuki/techtrain-mission-ca-tech-dojo-golang/repositories"
)

// サービス構造体を定義
type UserCharacterService struct {
	// userRepositoryを埋め込む
	ucRep repositories.UserCharacterRepository
}

// サービスのコンストラクタ
func NewUserCharacterService(r repositories.UserCharacterRepository) *UserCharacterService {
	return &UserCharacterService{ucRep: r}
}

// ハンドラー GetListHandler 用のサービスメソッド
func (s *UserCharacterService) List(userId string) (models.CharacterList, error) {
	userCharacters, err := s.ucRep.GetList(userId)
	if err != nil {
		return models.CharacterList{}, err
	}

	characterList := models.CharacterList{Characters: userCharacters}

	return characterList, nil
}
