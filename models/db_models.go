package models

import (
	"fmt"
	"time"
)

type Feedback struct {
	RecordId        string    `db:"record_id"`
	SourceId        string    `db:"source_id"`
	TenantId        string    `db:"tenant_id"`
	PersonName      string    `db:"person_name"`
	PersonEmail     string    `db:"person_email"`
	FeedbackType    string    `db:"feedback_type"`
	FeedbackContent string    `db:"feedback_content"`
	FeedbackLang    string    `db:"record_lang"`
	CreatedOn       time.Time `db:"created_on"`
	AdditionalData  string    `db:"additional_data"`
}

type SourceRecord struct {
	SourceId       string `db:"source_id"`
	SourceName     string `db:"source_name"`
	SourceMetadata string `db:"source_meta"`
}

type TenantRecord struct {
	TenantId   string `db:"tenant_id"`
	TenantName string `db:"tenant_name"`
	Tags       string `db:"product_tags"`
}

func (f *Feedback) ToFeedbackResponse(sourceInfo map[string][]string) *FeedbackResponse {
	resp := &FeedbackResponse{
		TenantId:        f.TenantId,
		FeedbackType:    f.FeedbackType,
		FeedbackLang:    f.FeedbackLang,
		FeedbackContent: f.FeedbackContent,
		PersonName:      f.PersonName,
		PersonEmail:     f.PersonEmail,
		FetchedOn:       fmt.Sprintf("%v", f.CreatedOn.Unix()),
		AdditionalData:  f.AdditionalData,
	}
	if len(sourceInfo[f.SourceId]) != 0 {
		resp.SourceName = sourceInfo[f.SourceId][0]
		resp.SourceMetaData = sourceInfo[f.SourceId][1]
	}
	return resp
}
