// Code generated by MockGen. DO NOT EDIT.
// Source: internal/services/password.go
//
// Generated by this command:
//
//	mockgen -source internal/services/password.go -destination internal/services/_mock/password.go -package mock_services
//

// Package mock_services is a generated GoMock package.
package mock_services

import (
	reflect "reflect"

	gomock "go.uber.org/mock/gomock"
)

// MockPasswordHasher is a mock of PasswordHasher interface.
type MockPasswordHasher struct {
	ctrl     *gomock.Controller
	recorder *MockPasswordHasherMockRecorder
	isgomock struct{}
}

// MockPasswordHasherMockRecorder is the mock recorder for MockPasswordHasher.
type MockPasswordHasherMockRecorder struct {
	mock *MockPasswordHasher
}

// NewMockPasswordHasher creates a new mock instance.
func NewMockPasswordHasher(ctrl *gomock.Controller) *MockPasswordHasher {
	mock := &MockPasswordHasher{ctrl: ctrl}
	mock.recorder = &MockPasswordHasherMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockPasswordHasher) EXPECT() *MockPasswordHasherMockRecorder {
	return m.recorder
}

// CompareHashAndPassword mocks base method.
func (m *MockPasswordHasher) CompareHashAndPassword(hashedPassword, password []byte) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CompareHashAndPassword", hashedPassword, password)
	ret0, _ := ret[0].(error)
	return ret0
}

// CompareHashAndPassword indicates an expected call of CompareHashAndPassword.
func (mr *MockPasswordHasherMockRecorder) CompareHashAndPassword(hashedPassword, password any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CompareHashAndPassword", reflect.TypeOf((*MockPasswordHasher)(nil).CompareHashAndPassword), hashedPassword, password)
}

// GenerateFromPassword mocks base method.
func (m *MockPasswordHasher) GenerateFromPassword(password []byte, cost int) ([]byte, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GenerateFromPassword", password, cost)
	ret0, _ := ret[0].([]byte)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GenerateFromPassword indicates an expected call of GenerateFromPassword.
func (mr *MockPasswordHasherMockRecorder) GenerateFromPassword(password, cost any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GenerateFromPassword", reflect.TypeOf((*MockPasswordHasher)(nil).GenerateFromPassword), password, cost)
}
