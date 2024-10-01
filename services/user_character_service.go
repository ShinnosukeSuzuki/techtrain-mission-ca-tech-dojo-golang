package services

import (
	"github.com/ShinnosukeSuzuki/techtrain-mission-ca-tech-dojo-golang/cache"
	"github.com/ShinnosukeSuzuki/techtrain-mission-ca-tech-dojo-golang/models"
	"github.com/ShinnosukeSuzuki/techtrain-mission-ca-tech-dojo-golang/repositories"
)

// サービス構造体を定義
type UserCharacterService struct {
	ucRep          repositories.UserCharacterRepository
	characterCache *cache.CharacterProbabilityCache
}

// サービスのコンストラクタ
func NewUserCharacterService(r repositories.UserCharacterRepository, characterCache *cache.CharacterProbabilityCache) *UserCharacterService {
	return &UserCharacterService{ucRep: r, characterCache: characterCache}
}

// ハンドラー GetListHandler 用のサービスメソッド
func (s *UserCharacterService) List(userID string) (models.CharacterList, error) {
	userCharacters, err := s.ucRep.GetList(userID)
	if err != nil {
		return models.CharacterList{}, err
	}

	// キャッシュからキャラクター情報を取得
	_, _, _, characterNameMap := s.characterCache.GetData()

	// キャラクターIDをキャラクター名に変換
	var userCharacterDetails []models.UserCharacterDetail
	for _, uc := range userCharacters {
		userCharacterDetails = append(userCharacterDetails, models.UserCharacterDetail{
			UserCharacterID: uc.UserCharacterID,
			CharacterID:     uc.CharacterID,
			Name:            characterNameMap[uc.CharacterID],
		})
	}

	characterList := models.CharacterList{Characters: userCharacterDetails}

	return characterList, nil
}
