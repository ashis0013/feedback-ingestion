package models

import (
	"fmt"
	"time"
)

type QueryFilter struct {
	TenantId  string
	SourceId  string
	StartTime *time.Time
	EndTime   *time.Time
}

func validateField[T comparable](field T, zero T, fieldName string, callback func(string)) {
	if field == zero {
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
	validateField(f.TenantId, "", fmt.Sprintf("%s tenant_id = $%d and", query, clauseCount), addClause)
	validateField(f.SourceId, "", fmt.Sprintf("%s source_id = $%d and", query, clauseCount), addClause)
	validateField(f.StartTime, nil, fmt.Sprintf("%s created_on >= $%d and", query, clauseCount), addClause)
	validateField(f.EndTime, nil, fmt.Sprintf("%s created_on <= $%d and", query, clauseCount), addClause)

	if query[len(query)-4:] == " and" {
		query = query[:len(query)-4]
	} else {
		query = query[:len(query)-6]
	}
	return query
}
