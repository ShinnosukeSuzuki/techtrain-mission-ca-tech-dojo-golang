package cron

import (
	"log"
	"time"

	"github.com/go-co-op/gocron/v2"
)

type CronJob struct {
	scheduler gocron.Scheduler
}

func NewCronJob() (*CronJob, error) {
	jst, err := time.LoadLocation("Asia/Tokyo")
	if err != nil {
		return nil, err
	}

	s, err := gocron.NewScheduler(gocron.WithLocation(jst))
	if err != nil {
		return nil, err
	}

	return &CronJob{scheduler: s}, nil
}

func (c *CronJob) AddJob(cronExpression string, job func() error) error {
	_, err := c.scheduler.NewJob(
		gocron.CronJob(cronExpression, false),
		gocron.NewTask(
			func() {
				if err := job(); err != nil {
					log.Printf("Failed to execute job: %v", err)
				}
			},
		),
	)

	return err
}

func (c *CronJob) Start() {
	c.scheduler.Start()
}

func (c *CronJob) Stop() error {
	return c.scheduler.Shutdown()
}
