package service_test

import (
	"time"

	"github.com/ashis0013/feedback-ingestion/cron"
	"github.com/ashis0013/feedback-ingestion/models"
	"github.com/ashis0013/feedback-ingestion/service"
	. "github.com/ashis0013/gollections"
	"github.com/stretchr/testify/mock"
)

var _ = Describe("Tests for feedback ingestion service", func() {
	var svc *service.FeedbackIngestionService
	var mav *MockAssertVariables

	BeforeEach(func() {
		mav = NewMockAssertVariables()
		svc = service.NewFeedbackIngestionService(mav.mockRepo, []cron.IngestionModule{mav.mockCron1})
		mav.mockCron1.On("PullData").Return([]*models.Feedback{})
		mav.mockCron1.On("GetSleepDuration").Return(time.Nanosecond)
		mav.mockRepo.On("AddRecord", []*models.Feedback{}).Return(nil)
	})

	Context("AppendIngestionModule()", func() {
		It("should fire the processe when called with an ingestion module", func() {
			mav.mockCron2.On("PullData").Return([]*models.Feedback{})
			mav.mockCron2.On("GetSleepDuration").Return(time.Nanosecond)

			svc.AppendIngestionModule(mav.mockCron2)
			time.Sleep(time.Second)

			Expect(len(mav.mockCron1.Calls) > 0).Should(BeTrue())
		})
	})

	Context("OnboardTenant()", func() {
		It("should called correct repository function with proper arguments", func() {
			tenantName := new(string)
			tags := new([]string)
			mav.mockRepo.On("AddTenant", mock.Anything, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
				*tenantName = args.Get(0).(string)
				*tags = args.Get(1).([]string)
			})

			svc.OnboardTenant(buildAddTenantReq())

			Expect(*tenantName).Should(Equal("tenat1"))
			Expect(*tags).Should(Equal([]string{"tag2"}))
		})
	})

	Context("GetFeedback()", func() {
		It("should called correct repository function with proper arguments", func() {
			filter := &models.QueryFilter{}
			mav.mockRepo.On("GetRecords", filter).Return(buildGetFeedbackResponse(), nil)

			resp, err := svc.GetFeedback(filter)

			Expect(resp).Should(Equal(buildGetFeedbackResponse()))
			Expect(err).Should(BeNil())
		})
	})

	Context("AddFedback", func() {
		It("should called correct repository function with proper arguments", func() {
			param := new([]*models.Feedback)
			mav.mockRepo.On("AddRecord", mock.Anything).Return(nil).Run(func(args mock.Arguments) {
				*param = args.Get(0).([]*models.Feedback)
			})

			svc.AddFedback(buildAddFeedbackRequest())

			Expect(*param).Should(Equal(Map(buildAddFeedbackRequest().Feedbacks, func(it *models.FeedbackRequest) *models.Feedback {
				return it.ToFeedback()
			})))
		})
	})
})

func buildAddTenantReq() *models.AddTenantRequest {
	return &models.AddTenantRequest{
		TenantName: "tenat1",
		Tags:       []string{"tag2"},
	}
}

func buildGetFeedbackResponse() *models.GetFeedbacksResponse {
	return &models.GetFeedbacksResponse{
		Feedbacks: []*models.FeedbackResponse{
			{
				FeedbackContent: "content",
			},
		},
	}
}

func buildAddFeedbackRequest() *models.AddFeedbackRequest {
	return &models.AddFeedbackRequest{
		Feedbacks: []*models.FeedbackRequest{
			{
				SourceId:        "s1",
				FeedbackContent: "this is a post",
			},
			{
				SourceId:        "s1",
				FeedbackContent: "another psot",
			},
		},
	}
}
