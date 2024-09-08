package services

import (
	"math/rand/v2"

	"github.com/ShinnosukeSuzuki/techtrain-mission-ca-tech-dojo-golang/models"
	"github.com/ShinnosukeSuzuki/techtrain-mission-ca-tech-dojo-golang/pkg/uuid"
	"github.com/ShinnosukeSuzuki/techtrain-mission-ca-tech-dojo-golang/repositories"
)

// サービス構造体を定義
type GachaDrawService struct {
	// UserCharacterRepositoryを埋め込む
	ucRep repositories.UserCharacterRepository
	// CharacterRepositoryを埋め込む
	cRep repositories.CharacterRepository
}

// サービスのコンストラクタ
func NewGachaDrawService(ucRep repositories.UserCharacterRepository, cRep repositories.CharacterRepository) *GachaDrawService {
	return &GachaDrawService{ucRep: ucRep, cRep: cRep}
}

// ハンドラー GachaDrawHandler 用のサービスメソッド
func (s *GachaDrawService) Draw(times int, userId string) ([]models.GachaResult, error) {
	// キャラクター全件取得
	characters, err := s.cRep.GetAllList()
	if err != nil {
		return nil, err
	}

	// ガチャの結果をDBにバルクインサートするための構造体を作成
	var userCharacterInserts []models.UserCharacterInsert
	// ガチャ結果を保存するための構造体を作成
	var gachaResults []models.GachaResult
	for i := 0; i < times; i++ {
		// ガチャのIDをuuidで生成
		id := uuid.GenerateUUID()

		character := selectRandomCharacter(characters)
		// ガチャの結果をDBにバルクインサートするための構造体に追加
		userCharacterInserts = append(userCharacterInserts, models.UserCharacterInsert{
			ID:          id,
			CharacterID: character.ID,
		})
		// ガチャ結果を保存するための構造体に追加
		gachaResults = append(gachaResults, models.GachaResult{
			CharacterID: character.ID,
			Name:        character.Name,
		})
	}

	// ガチャの結果をDBにバルクインサート
	if err := s.ucRep.InsertBulk(userId, userCharacterInserts); err != nil {
		return nil, err
	}

	return gachaResults, nil
}

// ガチャロジックを実装する
// キャラクターの確率に応じてランダムにキャラクターを選択する
func selectRandomCharacter(characters []models.Character) *models.Character {
	totalProbability := 0.0
	for _, char := range characters {
		totalProbability += char.Probability
	}

	randomValue := rand.Float64() * totalProbability
	cumulativeProbability := 0.0

	for _, char := range characters {
		cumulativeProbability += char.Probability
		if cumulativeProbability > randomValue {
			return &char
		}
	}

	return nil
}
