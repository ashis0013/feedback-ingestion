package cron

import (
	"time"

	"github.com/ashis0013/feedback-ingestion/models"
)

type IngestionModule interface {
	GetSleepDuration() time.Duration
	PullData() []*models.Feedback
}
