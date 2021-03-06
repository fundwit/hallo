// Code generated by MockGen. DO NOT EDIT.
// Source: hallo/domain (interfaces: AccountManager)

// Package domain is a generated GoMock package.
package domain

import (
	gomock "github.com/golang/mock/gomock"
	entity "hallo/domain/entity"
	reflect "reflect"
)

// MockAccountManager is a mock of AccountManager interface
type MockAccountManager struct {
	ctrl     *gomock.Controller
	recorder *MockAccountManagerMockRecorder
}

// MockAccountManagerMockRecorder is the mock recorder for MockAccountManager
type MockAccountManagerMockRecorder struct {
	mock *MockAccountManager
}

// NewMockAccountManager creates a new mock instance
func NewMockAccountManager(ctrl *gomock.Controller) *MockAccountManager {
	mock := &MockAccountManager{ctrl: ctrl}
	mock.recorder = &MockAccountManagerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockAccountManager) EXPECT() *MockAccountManagerMockRecorder {
	return m.recorder
}

// AuthenticateInternalIdentity mocks base method
func (m *MockAccountManager) AuthenticateInternalIdentity(arg0, arg1 string) (*entity.Account, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AuthenticateInternalIdentity", arg0, arg1)
	ret0, _ := ret[0].(*entity.Account)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// AuthenticateInternalIdentity indicates an expected call of AuthenticateInternalIdentity
func (mr *MockAccountManagerMockRecorder) AuthenticateInternalIdentity(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AuthenticateInternalIdentity", reflect.TypeOf((*MockAccountManager)(nil).AuthenticateInternalIdentity), arg0, arg1)
}

// CreateAccount mocks base method
func (m *MockAccountManager) CreateAccount(arg0 entity.EmailAccountCreateRequest) (*entity.Account, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateAccount", arg0)
	ret0, _ := ret[0].(*entity.Account)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateAccount indicates an expected call of CreateAccount
func (mr *MockAccountManagerMockRecorder) CreateAccount(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateAccount", reflect.TypeOf((*MockAccountManager)(nil).CreateAccount), arg0)
}
