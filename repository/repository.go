package repository

import (
	"github.com/ashis0013/feedback-ingestion/models"
)

type Repository interface {
	Init()
	AddRecord(record *models.Feedback)
	GetRecords(filter *models.QueryFilter) []*models.Feedback
}
