package service_test

import (
	"testing"

	mockCron "github.com/ashis0013/feedback-ingestion/mocks/cron"
	mockRepo "github.com/ashis0013/feedback-ingestion/mocks/repository"
	"github.com/onsi/ginkgo"
	"github.com/onsi/gomega"
)

var (
	BeforeEach     = ginkgo.BeforeEach
	BeEquivalentTo = gomega.BeEquivalentTo
	BeNil          = gomega.BeNil
	BeFalse        = gomega.BeFalse
	BeTrue         = gomega.BeTrue
	Context        = ginkgo.Context
	Describe       = ginkgo.Describe
	Equal          = gomega.Equal
	Expect         = gomega.Expect
	It             = ginkgo.It
	When           = ginkgo.When
)

type MockAssertVariables struct {
	mockRepo  *mockRepo.Repository
	mockCron1 *mockCron.IngestionModule
	mockCron2 *mockCron.IngestionModule
}

func NewMockAssertVariables() *MockAssertVariables {
	return &MockAssertVariables{
		mockRepo:  new(mockRepo.Repository),
		mockCron1: new(mockCron.IngestionModule),
		mockCron2: new(mockCron.IngestionModule),
	}
}

func TestGollections(t *testing.T) {
	gomega.RegisterFailHandler(ginkgo.Fail)
	ginkgo.RunSpecs(t, "Service Test Suite")
}
