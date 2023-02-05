package models

import (
	"fmt"
	"strconv"
	"time"
)

type QueryFilter struct {
	TenantId  string `json:"tenant_id"`
	SourceId  string `json:"source_id"`
	StartTime string `json:"start"`
	EndTime   string `json:"end"`
}

func takeNonEmpty(field string, fieldName string, callback func(string)) {
	if field == "" {
		return
	}
	if callback != nil {
		callback(fieldName)
	}
}

func (f *QueryFilter) BuildSelectquery() string {
	query := "select * from public.feedback_records where"
	clauseCount := 1
	addClause := func(field string) {
		query = field
		clauseCount++
	}
	takeNonEmpty(f.TenantId, fmt.Sprintf("%s tenant_id = $%d and", query, clauseCount), addClause)
	takeNonEmpty(f.SourceId, fmt.Sprintf("%s source_id = $%d and", query, clauseCount), addClause)
	takeNonEmpty(f.StartTime, fmt.Sprintf("%s created_on >= $%d and", query, clauseCount), addClause)
	takeNonEmpty(f.EndTime, fmt.Sprintf("%s created_on <= $%d and", query, clauseCount), addClause)

	if query[len(query)-4:] == " and" {
		query = query[:len(query)-4]
	} else {
		query = query[:len(query)-6]
	}
	return query
}

func (f *QueryFilter) GetQueryArgs() []any {
	args := []any{}
	if f.TenantId != "" {
		args = append(args, f.TenantId)
	}
	if f.SourceId != "" {
		args = append(args, f.SourceId)
	}
	if f.StartTime != "" {
		args = append(args, toTimestamp(f.StartTime))
	}
	if f.EndTime != "" {
		args = append(args, toTimestamp(f.EndTime))
	}
	return args
}

func (f *QueryFilter) IsInvalid() bool {
	if _, err := strconv.ParseInt(f.StartTime, 10, 64); err != nil {
		return f.StartTime != ""
	}
	if _, err := strconv.ParseInt(f.EndTime, 10, 64); err != nil {
		return f.EndTime != ""
	}
	return false
}

func toTimestamp(timestamp string) string {
	i64, _ := strconv.ParseInt(timestamp, 10, 64)
	return time.Unix(i64, 0).Format("2006-01-02 15:04:05")
}
