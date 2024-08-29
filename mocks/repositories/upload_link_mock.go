// Code generated by MockGen. DO NOT EDIT.
// Source: src/repositories/upload_link.go

// Package mocks is a generated GoMock package.
package mocks

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	models "github.com/tam-code/image-upload/src/models"
)

// MockUploadLinkRepository is a mock of UploadLinkRepository interface.
type MockUploadLinkRepository struct {
	ctrl     *gomock.Controller
	recorder *MockUploadLinkRepositoryMockRecorder
}

// MockUploadLinkRepositoryMockRecorder is the mock recorder for MockUploadLinkRepository.
type MockUploadLinkRepositoryMockRecorder struct {
	mock *MockUploadLinkRepository
}

// NewMockUploadLinkRepository creates a new mock instance.
func NewMockUploadLinkRepository(ctrl *gomock.Controller) *MockUploadLinkRepository {
	mock := &MockUploadLinkRepository{ctrl: ctrl}
	mock.recorder = &MockUploadLinkRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockUploadLinkRepository) EXPECT() *MockUploadLinkRepositoryMockRecorder {
	return m.recorder
}

// CreateUploadLink mocks base method.
func (m *MockUploadLinkRepository) CreateUploadLink(arg0 models.UploadLink) (*models.UploadLink, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateUploadLink", arg0)
	ret0, _ := ret[0].(*models.UploadLink)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateUploadLink indicates an expected call of CreateUploadLink.
func (mr *MockUploadLinkRepositoryMockRecorder) CreateUploadLink(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateUploadLink", reflect.TypeOf((*MockUploadLinkRepository)(nil).CreateUploadLink), arg0)
}

// GetUploadLinkByID mocks base method.
func (m *MockUploadLinkRepository) GetUploadLinkByID(arg0 string) (*models.UploadLink, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUploadLinkByID", arg0)
	ret0, _ := ret[0].(*models.UploadLink)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUploadLinkByID indicates an expected call of GetUploadLinkByID.
func (mr *MockUploadLinkRepositoryMockRecorder) GetUploadLinkByID(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUploadLinkByID", reflect.TypeOf((*MockUploadLinkRepository)(nil).GetUploadLinkByID), arg0)
}
