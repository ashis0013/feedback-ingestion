package cron_test

import (
	"encoding/json"
	"strings"
	"time"

	"github.com/ashis0013/feedback-ingestion/cron"
	mockCron "github.com/ashis0013/feedback-ingestion/mocks/cron"
	"github.com/ashis0013/feedback-ingestion/models"
)

type wrappedClient struct {
	mockCron.HttpClient
}

func (c *wrappedClient) Get(url string) ([]byte, error) {
	if strings.Contains(url, "page=0") {
		return []byte(pageResponse), nil
	} else if url == "https://meta.discourse.org/t/171809/posts.json?post_ids[]=882453" {
		return []byte(postResponse), nil
	}
	return []byte{}, &json.MarshalerError{}
}

var _ = Describe("Tests for discourse ingestion module", func() {
	var mav *MockAssertVariables
	var di *cron.DiscourseIngestor
	const sleepTime = time.Hour

	BeforeEach(func() {
		mav = NewMockAssertVariables()
		mav.mockRepo.On("FetchTags").Return(map[string]string{"prd1": "tenant1"}, nil)
		di = cron.NewDiscourseIngestor(mav.mockRepo, new(wrappedClient), "s1", sleepTime)
		time.Sleep(time.Second)
	})

	Context("GetSleepDuration()", func() {
		It("should return the specifid sleep duration", func() {
			Expect(di.GetSleepDuration()).Should(Equal(sleepTime))
		})
	})

	Context("GetSourceName()", func() {
		It("should return the specifid source name", func() {
			Expect(di.GetSourceName()).Should(Equal("Discourse"))
		})
	})

	Context("GetMetadata()", func() {
		It("should return the correct metadata", func() {
			meta, err := json.Marshal(&cron.Item{})
			if err != nil {
				Expect(false).Should(BeTrue()) //force fail the test
			}
			Expect(di.GetMetadata()).Should(Equal(string(meta)))
		})
	})

	Context("SetSourceId()", func() {
		It("should set the specifid source id", func() {
			di.SetSourceId("s2")

			Expect(di.SourceId).Should(Equal("s2"))
		})
	})

	Context("PullData()", func() {
		It("should start fetching data from discourse api", func() {
			feedbacks := di.PullData()

			Expect(len(feedbacks)).Should(Equal(1))
			Expect(*feedbacks[0]).Should(Equal(*buildMockFeedback()))
		})
	})
})

func buildMockFeedback() *models.Feedback {
	return &models.Feedback{
		SourceId:        "s1",
		TenantId:        "tenant1",
		FeedbackType:    "Post",
		FeedbackContent: "raw post data about prd1",
		FeedbackLang:    "English",
		AdditionalData:  `{"id":882453,"name":"","username":"","avatar_template":"","created_at":"","like_count":0,"blurb":"","post_number":0,"topic_title_headline":"","topic_id":171809}`,
	}
}

const (
	postResponse = `{
        "post_stream": {
            "posts": [
                {
                    "cooked": "raw post data about prd1"
                }
            ]
        },
        "id": 171809
    }`

	pageResponse = `{"posts":[{"id": 882453, "topic_id": 171809}]}`
)
