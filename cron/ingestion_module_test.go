package cron_test

import (
	"github.com/ashis0013/feedback-ingestion/cron"
	"github.com/stretchr/testify/mock"
)

var _ = Describe("test for module functions", func() {
	var mav *MockAssertVariables

	BeforeEach(func() {
		mav = NewMockAssertVariables()
		mav.mockRepo.On("FetchTags").Return(map[string]string{}, nil)
	})

	Context("InitModules()", func() {
		It("should assign correct source ids for existing modules", func() {
			mav.mockRepo.On("GetSources").Return(map[string]string{"Discourse": "s1"}, nil)

			modules := cron.InitModules(mav.mockRepo)
			discourseModule := modules[0].(*cron.DiscourseIngestor)

			Expect(discourseModule.SourceId).Should(Equal("s1"))
		})

		It("should register new modules", func() {
			mav.mockRepo.On("GetSources").Return(map[string]string{}, nil)
			mav.mockRepo.On("AddSource", mock.Anything, mock.Anything).Return("s2", nil)

			modules := cron.InitModules(mav.mockRepo)
			discourseModule := modules[0].(*cron.DiscourseIngestor)

			Expect(len(mav.mockRepo.Calls) > 0).Should(BeTrue())
			Expect(discourseModule.SourceId).Should(Equal("s2"))
		})
	})
})
