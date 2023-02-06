package cron

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/ashis0013/feedback-ingestion/models"
	"github.com/ashis0013/feedback-ingestion/repository"
	. "github.com/ashis0013/gollections"
)

type DiscourseIngestor struct {
	SourceId      string
	sleepDuration time.Duration
	batchSize     int
	maxFetch      int
	repository    repository.Repository
	httpClient    HttpClient
	tags          map[string]string
}

func NewDiscourseIngestor(repo repository.Repository, client HttpClient, id string, sleppTime time.Duration) *DiscourseIngestor {
	instance := &DiscourseIngestor{
		SourceId:      id,
		sleepDuration: sleppTime,
		batchSize:     2,
		maxFetch:      10,
		repository:    repo,
		httpClient:    client,
		tags:          make(map[string]string),
	}
	instance.repositoryRoutine()
	return instance
}

func (di *DiscourseIngestor) GetSleepDuration() time.Duration {
	return di.sleepDuration
}

func (di *DiscourseIngestor) PullData() []*models.Feedback {
	workerCount := di.maxFetch / di.batchSize
	jobs := make(chan []int, workerCount)
	results := make(chan []*models.Feedback, workerCount)
	for i := 0; i < workerCount; i++ {
		go di.worker(jobs, results)
	}
	for i := 0; i < workerCount; i++ {
		jobs <- []int{i * di.batchSize, (i + 1) * di.batchSize}
	}
	close(jobs)

	feedbacks := []*models.Feedback{}
	for i := 0; i < workerCount; i++ {
		feedbacks = append(feedbacks, <-results...)
	}
	return feedbacks
}

func (di *DiscourseIngestor) GetSourceName() string {
	return "Discourse"
}
func (di *DiscourseIngestor) GetMetadata() string {
	meta, _ := json.Marshal(&Item{})
	return string(meta)
}

func (di *DiscourseIngestor) SetSourceId(id string) {
	di.SourceId = id
}

func (di *DiscourseIngestor) repositoryRoutine() {
	go func(i *DiscourseIngestor) {
		for {
			tags, err := i.repository.FetchTags()
			if err != nil {
				log.Fatalln(err)
			}
			ForEachEntry(tags, func(key string, val string) { i.tags[key] = val })
			time.Sleep(time.Minute * 5)
		}
	}(di)
}

func (di *DiscourseIngestor) worker(jobs chan []int, res chan []*models.Feedback) {
	for pages := range jobs {
		res <- di.fetchData(pages[0], pages[1])
	}
}

func (di *DiscourseIngestor) fetchData(pageFrom int, pageTo int) []*models.Feedback {
	posts := []*Item{}
	validPosts := make(map[int][]string)
	for page := pageFrom; page < pageTo; page++ {
		body, err := di.httpClient.Get(getUrl(page))
		if err != nil {
			continue
		}
		var respStruct discourseResp
		err = json.Unmarshal(body, &respStruct)
		if err != nil {
			continue
		}
		ForEach(respStruct.Posts, func(it Item) {
			body, err = di.httpClient.Get(getPostUrl(it.TopicId, it.Id))
			var postStruct postResp
			json.Unmarshal(body, &postStruct)
			if len(postStruct.PostStream.Posts) == 0 {
				return
			}
			ForEachEntry(di.tags, func(key string, val string) {
				rawPost := postStruct.PostStream.Posts[0].Cooked
				if strings.Contains(rawPost, key) {
					validPosts[it.Id] = []string{val, rawPost}
				}
			})
			time.Sleep(time.Second)
		})
		ForEach(respStruct.Posts, func(it Item) {
			if ContainsKey(validPosts, it.Id) {
				posts = append(posts, &it)
			}
		})

		time.Sleep(time.Second)
	}
	return Map(posts, func(it *Item) *models.Feedback {
		feedback := it.toFeedback()
		feedback.TenantId = validPosts[it.Id][0]
		feedback.SourceId = di.SourceId
		feedback.FeedbackContent = validPosts[it.Id][1]
		return feedback
	})
}

func getUrl(page int) string {
	url := "https://meta.discourse.org/search.json?page=%v"
	url = fmt.Sprintf(url, page)
	after := "&q=after%3A"
	before := "before%3A"
	yesterday := time.Now().AddDate(0, 0, -1).Format("2006-01-02")
	dayBefore := time.Now().AddDate(0, 0, -2).Format("2006-01-02")
	return fmt.Sprintf("%s%s%s%s%s", url, after, dayBefore, before, yesterday)
}

func getPostUrl(topicId int, postId int) string {
	return fmt.Sprintf("https://meta.discourse.org/t/%d/posts.json?post_ids[]=%d", topicId, postId)
}

type discourseResp struct {
	Posts []Item `json:"posts"`
}
type Item struct {
	Id                 int    `json:"id"`
	Name               string `json:"name"`
	Username           string `json:"username"`
	AvatarTemplate     string `json:"avatar_template"`
	CreatedAt          string `json:"created_at"`
	LikeCount          int    `json:"like_count"`
	Blurb              string `json:"blurb"`
	PostNumber         int    `json:"post_number"`
	TopicTitleHeadline string `json:"topic_title_headline"`
	TopicId            int    `json:"topic_id"`
}

type postResp struct {
	PostStream Stream `json:"post_stream"`
	Id         int    `json:"id"`
}

type Stream struct {
	Posts []PostItem `json:"posts"`
}

type PostItem struct {
	Cooked string `json:"cooked"`
}

func (apiItem *Item) toFeedback() *models.Feedback {
	enc, _ := json.Marshal(apiItem)
	return &models.Feedback{
		PersonName:     apiItem.Username,
		FeedbackType:   "Post",
		FeedbackLang:   "English",
		AdditionalData: string(enc),
	}
}
