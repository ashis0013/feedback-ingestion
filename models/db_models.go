package models

import "time"

type Feedback struct {
	RecordId        string    `db:"record_id"`
	SourceId        string    `db:"source_id"`
	TenantId        string    `db:"tenant_id"`
	PersonName      string    `db:"person_name"`
	PersonEmail     string    `db:"person_email"`
	FeedbackType    string    `db:"feedback_type"`
	FeedbackContent string    `db:"feedback_content"`
	RecordLang      string    `db:"record_lang"`
	CreatedOn       time.Time `db:"created_on"`
	AdditionalData  string    `db:"additional_data"`
}
