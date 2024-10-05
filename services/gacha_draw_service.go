package services

import (
	"math/rand/v2"

	"github.com/ShinnosukeSuzuki/techtrain-mission-ca-tech-dojo-golang/cache"
	"github.com/ShinnosukeSuzuki/techtrain-mission-ca-tech-dojo-golang/models"
	"github.com/ShinnosukeSuzuki/techtrain-mission-ca-tech-dojo-golang/repositories"
)

// サービス構造体を定義
type GachaDrawService struct {
	ucRep          repositories.UserCharacterRepository
	characterCache *cache.CharacterProbabilityCache
}

// サービスのコンストラクタ
func NewGachaDrawService(ucRep repositories.UserCharacterRepository, characterCache *cache.CharacterProbabilityCache) *GachaDrawService {
	return &GachaDrawService{ucRep: ucRep, characterCache: characterCache}
}

// ハンドラー GachaDrawHandler 用のサービスメソッド
func (s *GachaDrawService) Draw(times int, userID string) ([]models.Character, error) {
	// キャラクター全件取得をキャッシュから取得
	characters, cumulativeProbabilities := s.characterCache.GetDataForGacha()

	// ガチャを行いキャラクターを選択
	gachaResults := selectRandomCharacters(times, characters, cumulativeProbabilities)

	// ガチャの結果をDBにバルクインサート
	if err := s.ucRep.InsertBulk(userID, gachaResults); err != nil {
		return nil, err
	}

	return gachaResults, nil
}

// ガチャロジックを実装する
// キャラクターの確率に応じてランダムにキャラクターを選択する
func selectRandomCharacters(times int, characters []models.Character, cumulativeProbabilities []float64) []models.Character {

	gachaResults := make([]models.Character, times)

	for i := 0; i < times; i++ {
		// 0~totalProbabilityの範囲で乱数を生成
		randomValue := rand.Float64() * cumulativeProbabilities[len(cumulativeProbabilities)-1]

		// 二分探索で乱数に対応するキャラクターを選択
		left, right := 0, len(cumulativeProbabilities)-1
		for left < right {
			mid := (left + right) / 2
			if cumulativeProbabilities[mid] < randomValue {
				left = mid + 1
			} else {
				right = mid
			}
		}

		gachaResults[i] = characters[left]
	}

	return gachaResults
}
