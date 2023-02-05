package models

type AddFeedbackRequest struct {
	Feedbacks []*FeedbackRequest `json:"feedbacks"`
}

type FeedbackRequest struct {
	SourceId        string `json:"source_id"`
	TenantId        string `json:"tenant_id"`
	PersonName      string `json:"person_name"`
	PersonEmail     string `json:"person_email"`
	FeedbackType    string `json:"feedback_type"`
	FeedbackContent string `json:"feedback_content"`
	FeedbackLang    string `json:"record_lang"`
	AdditionalData  string `json:"additional_data"`
}

type GetFeedbacksResponse struct {
	Feedbacks []*FeedbackResponse `json:"feebnacks"`
}

type FeedbackResponse struct {
	TenantId        string `json:"tenant_id"`
	SourceName      string `json:"source_name"`
	SourceMetaData  string `json:"source_meta"`
	FeedbackType    string `json:"feedback_type"`
	FeedbackLang    string `json:"feedback_lang"`
	FeedbackContent string `json:"feedback_content"`
	PersonName      string `json:"person_name"`
	PersonEmail     string `json:"person_email"`
	FetchedOn       string `json:"fetched_on"`
	AdditionalData  string `json:"additional_data"`
}

type AddTenantRequest struct {
	TenantName string   `json:"tenant_name"`
	Tags       []string `json:"tags"`
}

func (req *FeedbackRequest) ToFeedback() *Feedback {
	return &Feedback{
		TenantId:        req.TenantId,
		SourceId:        req.SourceId,
		PersonName:      req.PersonName,
		PersonEmail:     req.PersonEmail,
		FeedbackType:    req.FeedbackType,
		FeedbackContent: req.FeedbackContent,
		FeedbackLang:    req.FeedbackLang,
		AdditionalData:  req.AdditionalData,
	}
}
