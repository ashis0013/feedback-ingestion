// Code generated by mockery v2.18.0. DO NOT EDIT.

package mocks

import (
	models "github.com/ashis0013/feedback-ingestion/models"
	mock "github.com/stretchr/testify/mock"
)

// Repository is an autogenerated mock type for the Repository type
type Repository struct {
	mock.Mock
}

// AddRecord provides a mock function with given fields: record
func (_m *Repository) AddRecord(record []*models.Feedback) error {
	ret := _m.Called(record)

	var r0 error
	if rf, ok := ret.Get(0).(func([]*models.Feedback) error); ok {
		r0 = rf(record)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// AddSource provides a mock function with given fields: sourceName, metadata
func (_m *Repository) AddSource(sourceName string, metadata string) (string, error) {
	ret := _m.Called(sourceName, metadata)

	var r0 string
	if rf, ok := ret.Get(0).(func(string, string) string); ok {
		r0 = rf(sourceName, metadata)
	} else {
		r0 = ret.Get(0).(string)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string, string) error); ok {
		r1 = rf(sourceName, metadata)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// AddTenant provides a mock function with given fields: tenantName, tags
func (_m *Repository) AddTenant(tenantName string, tags []string) error {
	ret := _m.Called(tenantName, tags)

	var r0 error
	if rf, ok := ret.Get(0).(func(string, []string) error); ok {
		r0 = rf(tenantName, tags)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// FetchTags provides a mock function with given fields:
func (_m *Repository) FetchTags() (map[string]string, error) {
	ret := _m.Called()

	var r0 map[string]string
	if rf, ok := ret.Get(0).(func() map[string]string); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(map[string]string)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetRecords provides a mock function with given fields: filter
func (_m *Repository) GetRecords(filter *models.QueryFilter) (*models.GetFeedbacksResponse, error) {
	ret := _m.Called(filter)

	var r0 *models.GetFeedbacksResponse
	if rf, ok := ret.Get(0).(func(*models.QueryFilter) *models.GetFeedbacksResponse); ok {
		r0 = rf(filter)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.GetFeedbacksResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(*models.QueryFilter) error); ok {
		r1 = rf(filter)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetSources provides a mock function with given fields:
func (_m *Repository) GetSources() (map[string]string, error) {
	ret := _m.Called()

	var r0 map[string]string
	if rf, ok := ret.Get(0).(func() map[string]string); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(map[string]string)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

type mockConstructorTestingTNewRepository interface {
	mock.TestingT
	Cleanup(func())
}

// NewRepository creates a new instance of Repository. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewRepository(t mockConstructorTestingTNewRepository) *Repository {
	mock := &Repository{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}