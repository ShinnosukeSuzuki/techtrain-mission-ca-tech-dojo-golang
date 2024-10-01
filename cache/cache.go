package cache

import (
	"encoding/csv"
	"log"
	"strconv"
	"sync"
	"time"

	"github.com/ShinnosukeSuzuki/techtrain-mission-ca-tech-dojo-golang/models"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/go-co-op/gocron/v2"
)

type CharacterProbabilityCache struct {
	Characters              []models.Character
	CumulativeProbabilities []float64
	TotalProbability        float64
	CharacterNameMap        map[string]string
	mutex                   sync.RWMutex
	s3Client                *s3.S3
	bucketName              string
	filePath                string
	scheduler               gocron.Scheduler
}

func NewCharacterProbabilityCache(region, bucketName, filePath string) (*CharacterProbabilityCache, error) {
	sess := session.Must(session.NewSession(&aws.Config{
		Region: aws.String(region),
	}))
	s3Client := s3.New(sess)

	cache := &CharacterProbabilityCache{
		Characters:              []models.Character{},
		CumulativeProbabilities: []float64{},
		TotalProbability:        0,
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

func (c *CharacterProbabilityCache) StartCron() error {
	// cron実行をJSTで行うための設定
	jst, err := time.LoadLocation("Asia/Tokyo")
	if err != nil {
		return err
	}

	s, err := gocron.NewScheduler(gocron.WithLocation(jst))
	if err != nil {
		return err
	}

	_, err = s.NewJob(
		// 毎日0時に実行
		gocron.CronJob("0 0 * * *", false),
		gocron.NewTask(
			func() {
				if err := c.Update(); err != nil {
					log.Printf("Failed to update cache: %v", err)
				}
			},
		),
	)

	if err != nil {
		return err
	}

	c.scheduler = s
	s.Start()
	return nil
}

func (c *CharacterProbabilityCache) StopCron() error {
	if c.scheduler != nil {
		return c.scheduler.Shutdown()
	}
	return nil
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
	records, err := reader.ReadAll()
	if err != nil {
		return err
	}

	newCharacters := make([]models.Character, 0, len(records))
	newCumulativeProbabilities := make([]float64, 0, len(records))
	var newTotalProbability float64
	newCharacterNameMap := make(map[string]string)

	for _, r := range records {
		// csvはID, Name, Probabilityの3つのカラムを持つことを前提とする
		if len(r) != 3 {
			continue
		}
		p, err := strconv.ParseFloat(r[2], 64)
		if err != nil {
			continue
		}
		newCharacters = append(newCharacters, models.Character{
			ID:          r[0],
			Name:        r[1],
			Probability: p,
		})
		newTotalProbability += p
		newCumulativeProbabilities = append(newCumulativeProbabilities, newTotalProbability)
		newCharacterNameMap[r[0]] = r[1]
	}

	c.mutex.Lock()
	c.Characters = newCharacters
	c.CumulativeProbabilities = newCumulativeProbabilities
	c.TotalProbability = newTotalProbability
	c.CharacterNameMap = newCharacterNameMap
	c.mutex.Unlock()

	log.Println("Cache updated successfully")
	return nil
}

func (c *CharacterProbabilityCache) GetData() ([]models.Character, []float64, float64, map[string]string) {
	// 読み込み用ロックを使うことで待ち時間を短縮
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	return c.Characters, c.CumulativeProbabilities, c.TotalProbability, c.CharacterNameMap
}
