package cron

import (
	"log"
	"time"

	"github.com/ashis0013/feedback-ingestion/models"
	"github.com/ashis0013/feedback-ingestion/repository"
	. "github.com/ashis0013/gollections"
)

type IngestionModule interface {
	GetSleepDuration() time.Duration
	PullData() []*models.Feedback
	GetSourceName() string
	GetMetadata() string
	SetSourceId(id string)
}

func InitModules(repo repository.Repository) []IngestionModule {
	sources, err := repo.GetSources()
	if err != nil {
		log.Fatalln(err)
	}
	crons := []IngestionModule{
		NewDiscourseIngestor(repo, &CronHttpClient{}, sources["Discourse"], time.Hour * 24),
	}

	ForEach(crons, func(cron IngestionModule) {
		if sources == nil || ContainsKey(sources, cron.GetSourceName()) {
			return
		}
		id, err := repo.AddSource(cron.GetSourceName(), cron.GetMetadata())
		if err != nil {
			log.Fatalln(err)
			return
		}
		cron.SetSourceId(id)
	})
	return crons
}
