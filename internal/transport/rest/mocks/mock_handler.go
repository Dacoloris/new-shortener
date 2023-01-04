// Code generated by MockGen. DO NOT EDIT.
// Source: ./internal/transport/rest/handler.go

// Package mock_rest is a generated GoMock package.
package mock_rest

import (
	context "context"
	domain "new-shortner/internal/domain"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockURLs is a mock of URLs interface.
type MockURLs struct {
	ctrl     *gomock.Controller
	recorder *MockURLsMockRecorder
}

// MockURLsMockRecorder is the mock recorder for MockURLs.
type MockURLsMockRecorder struct {
	mock *MockURLs
}

// NewMockURLs creates a new mock instance.
func NewMockURLs(ctrl *gomock.Controller) *MockURLs {
	mock := &MockURLs{ctrl: ctrl}
	mock.recorder = &MockURLsMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockURLs) EXPECT() *MockURLsMockRecorder {
	return m.recorder
}

// Create mocks base method.
func (m *MockURLs) Create(ctx context.Context, url domain.URL, baseURL string) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", ctx, url, baseURL)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Create indicates an expected call of Create.
func (mr *MockURLsMockRecorder) Create(ctx, url, baseURL interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockURLs)(nil).Create), ctx, url, baseURL)
}

// GetAllURLsByUserID mocks base method.
func (m *MockURLs) GetAllURLsByUserID(ctx context.Context, UserID string) ([]domain.URL, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAllURLsByUserID", ctx, UserID)
	ret0, _ := ret[0].([]domain.URL)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAllURLsByUserID indicates an expected call of GetAllURLsByUserID.
func (mr *MockURLsMockRecorder) GetAllURLsByUserID(ctx, UserID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAllURLsByUserID", reflect.TypeOf((*MockURLs)(nil).GetAllURLsByUserID), ctx, UserID)
}

// GetOriginalByShort mocks base method.
func (m *MockURLs) GetOriginalByShort(ctx context.Context, short string) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetOriginalByShort", ctx, short)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetOriginalByShort indicates an expected call of GetOriginalByShort.
func (mr *MockURLsMockRecorder) GetOriginalByShort(ctx, short interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetOriginalByShort", reflect.TypeOf((*MockURLs)(nil).GetOriginalByShort), ctx, short)
}
