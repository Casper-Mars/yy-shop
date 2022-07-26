// Code generated by MockGen. DO NOT EDIT.
// Source: yy-shop/internal/biz (interfaces: EncryptService,UserRepo)

// Package biz is a generated GoMock package.
package biz

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockEncryptService is a mock of EncryptService interface.
type MockEncryptService struct {
	ctrl     *gomock.Controller
	recorder *MockEncryptServiceMockRecorder
}

// MockEncryptServiceMockRecorder is the mock recorder for MockEncryptService.
type MockEncryptServiceMockRecorder struct {
	mock *MockEncryptService
}

// NewMockEncryptService creates a new mock instance.
func NewMockEncryptService(ctrl *gomock.Controller) *MockEncryptService {
	mock := &MockEncryptService{ctrl: ctrl}
	mock.recorder = &MockEncryptServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockEncryptService) EXPECT() *MockEncryptServiceMockRecorder {
	return m.recorder
}

// Encrypt mocks base method.
func (m *MockEncryptService) Encrypt(arg0 context.Context, arg1 []byte) ([]byte, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Encrypt", arg0, arg1)
	ret0, _ := ret[0].([]byte)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Encrypt indicates an expected call of Encrypt.
func (mr *MockEncryptServiceMockRecorder) Encrypt(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Encrypt", reflect.TypeOf((*MockEncryptService)(nil).Encrypt), arg0, arg1)
}

// Token mocks base method.
func (m *MockEncryptService) Token(arg0 context.Context, arg1 *User) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Token", arg0, arg1)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Token indicates an expected call of Token.
func (mr *MockEncryptServiceMockRecorder) Token(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Token", reflect.TypeOf((*MockEncryptService)(nil).Token), arg0, arg1)
}

// MockUserRepo is a mock of UserRepo interface.
type MockUserRepo struct {
	ctrl     *gomock.Controller
	recorder *MockUserRepoMockRecorder
}

// MockUserRepoMockRecorder is the mock recorder for MockUserRepo.
type MockUserRepoMockRecorder struct {
	mock *MockUserRepo
}

// NewMockUserRepo creates a new mock instance.
func NewMockUserRepo(ctrl *gomock.Controller) *MockUserRepo {
	mock := &MockUserRepo{ctrl: ctrl}
	mock.recorder = &MockUserRepoMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockUserRepo) EXPECT() *MockUserRepoMockRecorder {
	return m.recorder
}

// FetchByUsername mocks base method.
func (m *MockUserRepo) FetchByUsername(arg0 context.Context, arg1 string) (*User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FetchByUsername", arg0, arg1)
	ret0, _ := ret[0].(*User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FetchByUsername indicates an expected call of FetchByUsername.
func (mr *MockUserRepoMockRecorder) FetchByUsername(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FetchByUsername", reflect.TypeOf((*MockUserRepo)(nil).FetchByUsername), arg0, arg1)
}

// Save mocks base method.
func (m *MockUserRepo) Save(arg0 context.Context, arg1 *User) (int64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Save", arg0, arg1)
	ret0, _ := ret[0].(int64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Save indicates an expected call of Save.
func (mr *MockUserRepoMockRecorder) Save(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Save", reflect.TypeOf((*MockUserRepo)(nil).Save), arg0, arg1)
}
