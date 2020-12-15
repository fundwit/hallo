// Code generated by MockGen. DO NOT EDIT.
// Source: hallo/domain (interfaces: IdentityBindingRepository)

// Package domain is a generated GoMock package.
package domain

import (
	gomock "github.com/golang/mock/gomock"
	reflect "reflect"
)

// MockIdentityBindingRepository is a mock of IdentityBindingRepository interface
type MockIdentityBindingRepository struct {
	ctrl     *gomock.Controller
	recorder *MockIdentityBindingRepositoryMockRecorder
}

// MockIdentityBindingRepositoryMockRecorder is the mock recorder for MockIdentityBindingRepository
type MockIdentityBindingRepositoryMockRecorder struct {
	mock *MockIdentityBindingRepository
}

// NewMockIdentityBindingRepository creates a new mock instance
func NewMockIdentityBindingRepository(ctrl *gomock.Controller) *MockIdentityBindingRepository {
	mock := &MockIdentityBindingRepository{ctrl: ctrl}
	mock.recorder = &MockIdentityBindingRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockIdentityBindingRepository) EXPECT() *MockIdentityBindingRepositoryMockRecorder {
	return m.recorder
}

// Save mocks base method
func (m *MockIdentityBindingRepository) Save(arg0 uint64, arg1, arg2 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Save", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// Save indicates an expected call of Save
func (mr *MockIdentityBindingRepositoryMockRecorder) Save(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Save", reflect.TypeOf((*MockIdentityBindingRepository)(nil).Save), arg0, arg1, arg2)
}
