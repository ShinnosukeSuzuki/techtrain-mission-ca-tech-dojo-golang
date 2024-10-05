package cache

import (
	"encoding/csv"
	"fmt"
	"log"
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
}

type parsedRecord struct {
	index       int
	character   models.Character
	probability float64
}

func NewCharacterProbabilityCache(region, bucketName, filePath string) (*CharacterProbabilityCache, error) {
	sess := session.Must(session.NewSession(&aws.Config{
		Region: aws.String(region),
	}))
	s3Client := s3.New(sess)

	cache := &CharacterProbabilityCache{
		Characters:              []models.Character{},
		CumulativeProbabilities: []float64{},
		CharacterNameMap:        map[string]string{},
		s3Client:                s3Client,
		bucketName:              bucketName,
		filePath:                filePath,
	}

	// 初回のデータ読み込み
	if err := cache.Update(); err != nil {
		return nil, err
	}

	return cache, nil
}

func (c *CharacterProbabilityCache) Update() error {
	resp, err := c.s3Client.GetObject(&s3.GetObjectInput{
		Bucket: aws.String(c.bucketName),
		Key:    aws.String(c.filePath),
	})
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	reader := csv.NewReader(resp.Body)

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
