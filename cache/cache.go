package cache

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"runtime"
	"strconv"
	"sync"

	"github.com/ShinnosukeSuzuki/techtrain-mission-ca-tech-dojo-golang/models"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

type CharacterProbabilityCache struct {
	Characters              []models.Character
	CumulativeProbabilities []float64
	CharacterNameMap        map[string]string
	mutex                   sync.RWMutex
	s3Client                *s3.S3
	bucketName              string
	filePath                string
	useS3                   bool
}

type parsedRecord struct {
	index       int
	character   models.Character
	probability float64
}

func NewCharacterProbabilityCache() (*CharacterProbabilityCache, error) {
	env := os.Getenv("ENV")

	// ローカル環境かS3を使用するかを判定
	useS3 := env == "Prod" || env == "Dev"
	var bucketName, filePath, region string
	var s3Client *s3.S3

	if useS3 {
		// 環境変数からS3の設定を取得
		region = os.Getenv("REGION")
		bucketName = os.Getenv("BUCKET_NAME")
		filePath = os.Getenv("FILE_PATH")

		if region == "" || bucketName == "" || filePath == "" {
			return nil, fmt.Errorf("missing S3 configuration: REGION, BUCKET_NAME, and FILE_PATH must be set")
		}

		// S3クライアントを初期化
		sess := session.Must(session.NewSession(&aws.Config{
			Region: aws.String(region),
		}))
		s3Client = s3.New(sess)
	} else {
		// ローカル環境では固定ファイルパスを使用
		filePath = "infra/game-api-infrastructure/S3/characters.csv"
	}

	cache := &CharacterProbabilityCache{
		Characters:              []models.Character{},
		CumulativeProbabilities: []float64{},
		CharacterNameMap:        map[string]string{},
		s3Client:                s3Client,
		bucketName:              bucketName,
		filePath:                filePath,
		useS3:                   useS3,
	}

	// 初回のデータ読み込み
	if err := cache.Update(); err != nil {
		return nil, err
	}

	return cache, nil
}

func (c *CharacterProbabilityCache) Update() error {
	var reader *csv.Reader
	if c.useS3 {
		// S3からデータを取得
		resp, err := c.s3Client.GetObject(&s3.GetObjectInput{
			Bucket: aws.String(c.bucketName),
			Key:    aws.String(c.filePath),
		})
		if err != nil {
			return fmt.Errorf("failed to fetch from S3: %w", err)
		}
		defer resp.Body.Close()
		reader = csv.NewReader(resp.Body)
	} else {
		// ローカルファイルを読み込む
		file, err := os.Open(c.filePath)
		if err != nil {
			return fmt.Errorf("failed to open local file: %w", err)
		}
		defer file.Close()
		reader = csv.NewReader(file)
	}

	// ヘッダー行を読み飛ばす
	if _, err := reader.Read(); err != nil {
		return fmt.Errorf("failed to read header: %w", err)
	}

	records, err := reader.ReadAll()
	if err != nil {
		return err
	}

	// ワーカー数を設定（CPU数の2倍を使用）
	numWorkers := runtime.NumCPU() * 2

	// チャネルの作成
	jobs := make(chan struct {
		index  int
		record []string
	}, len(records))
	results := make(chan parsedRecord, len(records))
	errors := make(chan error, numWorkers)

	// ワーカーの起動
	var wg sync.WaitGroup
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go worker(jobs, results, errors, &wg)
	}

	// ジョブの送信
	for i, record := range records {
		jobs <- struct {
			index  int
			record []string
		}{i, record}
	}
	close(jobs)

	// 全てのワーカーの完了を待つ
	go func() {
		wg.Wait()
		close(results)
		close(errors)
	}()

	// エラーチェック
	for err := range errors {
		if err != nil {
			return fmt.Errorf("error in worker: %w", err)
		}
	}

	// 結果の集計と順序付け
	parsedRecords := make([]parsedRecord, len(records))
	for result := range results {
		parsedRecords[result.index] = result
	}

	// 累積確率の計算と最終データの作成
	newCharacters := make([]models.Character, 0, len(records))
	newCumulativeProbabilities := make([]float64, 0, len(records))
	var newTotalProbability float64
	newCharacterNameMap := make(map[string]string)

	for _, r := range parsedRecords {
		newCharacters = append(newCharacters, r.character)
		newTotalProbability += r.probability
		newCumulativeProbabilities = append(newCumulativeProbabilities, newTotalProbability)
		newCharacterNameMap[r.character.ID] = r.character.Name
	}

	c.mutex.Lock()
	c.Characters = newCharacters
	c.CumulativeProbabilities = newCumulativeProbabilities
	c.CharacterNameMap = newCharacterNameMap
	c.mutex.Unlock()

	log.Println("Cache updated successfully")
	return nil
}

func (c *CharacterProbabilityCache) GetDataForGacha() ([]models.Character, []float64) {
	// 読み込み用ロックを使うことで待ち時間を短縮
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	return c.Characters, c.CumulativeProbabilities
}

func (c *CharacterProbabilityCache) GetNameMap() map[string]string {
	// 読み込み用ロックを使うことで待ち時間を短縮
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	return c.CharacterNameMap
}

func worker(jobs <-chan struct {
	index  int
	record []string
}, results chan<- parsedRecord, errors chan<- error, wg *sync.WaitGroup) {
	defer wg.Done()
	for job := range jobs {
		if len(job.record) != 3 {
			errors <- fmt.Errorf("invalid record length at index %d: %v", job.index, job.record)
			continue
		}
		probability, err := strconv.ParseFloat(job.record[2], 64)
		if err != nil {
			errors <- fmt.Errorf("invalid probability at index %d: %w", job.index, err)
			continue
		}
		results <- parsedRecord{
			index: job.index,
			character: models.Character{
				ID:   job.record[0],
				Name: job.record[1],
			},
			probability: probability,
		}
	}
}
