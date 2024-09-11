package services

import (
	"math/rand/v2"

	"github.com/ShinnosukeSuzuki/techtrain-mission-ca-tech-dojo-golang/models"
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

	// ガチャを行いキャラクターを選択
	gachaResults := selectRandomCharacters(characters, times)

	// ガチャの結果をDBにバルクインサート
	if err := s.ucRep.InsertBulk(userId, gachaResults); err != nil {
		return nil, err
	}

	return gachaResults, nil
}

// ガチャロジックを実装する
// キャラクターの確率に応じてランダムにキャラクターを選択する
func selectRandomCharacters(characters []models.Character, times int) []models.GachaResult {
	// キャラクターの確率の合計を計算
	var totalProbability float64
	for _, c := range characters {
		totalProbability += c.Probability
	}

	gachaResults := make([]models.GachaResult, times)

	// 累積確率を計算
	cumulativeProbabilities := make([]float64, len(characters))
	cumulativeProbabilities[0] = characters[0].Probability
	for i := 1; i < len(characters); i++ {
		cumulativeProbabilities[i] = cumulativeProbabilities[i-1] + characters[i].Probability
	}

	for i := 0; i < times; i++ {
		// 0~totalProbabilityの範囲で乱数を生成
		randomValue := rand.Float64() * totalProbability

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

		selectedChar := characters[left]

		gachaResults[i] = models.GachaResult{
			CharacterID: selectedChar.ID,
			Name:        selectedChar.Name,
		}
	}

	return gachaResults
}
