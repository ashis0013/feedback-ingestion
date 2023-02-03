package main

import (
	// "time"

	"github.com/ashis0013/feedback-ingestion/models"
	"github.com/ashis0013/feedback-ingestion/repository"
)

func main() {
	pg := repository.NewPostgresRepository()
	defer pg.Terminate()
	pg.Init()
	// now := time.Now()
	// kal := time.Now().Add(time.Hour * (-24))

	pg.GetRecords(&models.QueryFilter{
		TenantId:  "",
		SourceId:  "",
		StartTime: nil,
		EndTime:   nil,
	})
}
