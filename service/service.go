package service

import (
	"time"

	"github.com/ashis0013/feedback-ingestion/cron"
	"github.com/ashis0013/feedback-ingestion/models"
	"github.com/ashis0013/feedback-ingestion/repository"
	. "github.com/ashis0013/gollections"
)

type FeedbackIngestionService struct {
	repo              repository.Repository
	ingestion_modules []cron.IngestionModule
}

func NewFeedbackIngestionService(repository repository.Repository, modules []cron.IngestionModule) *FeedbackIngestionService {
	return &FeedbackIngestionService{
		repo:              repository,
		ingestion_modules: modules,
	}
}

func (s *FeedbackIngestionService) AppendIngestionModule(module cron.IngestionModule) {
	s.ingestion_modules = append(s.ingestion_modules, module)
	s.fireCron(module)
}

func (s *FeedbackIngestionService) OnboardTenant(tenant *models.AddTenantRequest) error {
	return s.repo.AddTenant(tenant.TenantName, tenant.Tags)
}

func (s *FeedbackIngestionService) GetFeedback(filter *models.QueryFilter) (*models.GetFeedbacksResponse, error) {
	return s.repo.GetRecords(filter)
}

func (s *FeedbackIngestionService) AddFedback(feedbacks *models.AddFeedbackRequest) error {
	records := Map(feedbacks.Feedbacks, func(req *models.FeedbackRequest) *models.Feedback {
		return req.ToFeedback()
	})
	return s.repo.AddRecord(records)
}

func (s *FeedbackIngestionService) fireCron(module cron.IngestionModule) {
	go func(m cron.IngestionModule, repo repository.Repository) {
		feedbacks := m.PullData()
		repo.AddRecord(feedbacks)
		time.Sleep(m.GetSleepDuration())
	}(module, s.repo)
}
