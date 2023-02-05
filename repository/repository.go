package repository

import (
	"github.com/ashis0013/feedback-ingestion/models"
)

type Repository interface {
	AddRecord(record []*models.Feedback) error
	GetRecords(filter *models.QueryFilter) (*models.GetFeedbacksResponse, error)
	AddTenant(tenantName string, tags []string) error
	FetchTags() ([]string, error)
	AddSource(sourceName string, metadata string) error
}
