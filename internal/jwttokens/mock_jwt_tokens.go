// Code generated by MockGen. DO NOT EDIT.
// Source: gitlab.com/mediasoft-internship/internship/mamoru777/authservice/internal/jwttokens (interfaces: ITokens)
//
// Generated by this command:
//
//	mockgen -destination mock_jwt_tokens.go -package jwttokens . ITokens
//
// Package jwttokens is a generated GoMock package.
package jwttokens

import (
	reflect "reflect"

	jwt "github.com/dgrijalva/jwt-go"
	uuid "github.com/google/uuid"
	mylogger "gitlab.com/mediasoft-internship/internship/mamoru777/authservice/internal/mylogger"
	gomock "go.uber.org/mock/gomock"
)

// MockITokens is a mock of ITokens interface.
type MockITokens struct {
	ctrl     *gomock.Controller
	recorder *MockITokensMockRecorder
}

// MockITokensMockRecorder is the mock recorder for MockITokens.
type MockITokensMockRecorder struct {
	mock *MockITokens
}

// NewMockITokens creates a new mock instance.
func NewMockITokens(ctrl *gomock.Controller) *MockITokens {
	mock := &MockITokens{ctrl: ctrl}
	mock.recorder = &MockITokensMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockITokens) EXPECT() *MockITokensMockRecorder {
	return m.recorder
}

// CreateAccessToken mocks base method.
func (m *MockITokens) CreateAccessToken(arg0 uuid.UUID) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateAccessToken", arg0)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateAccessToken indicates an expected call of CreateAccessToken.
func (mr *MockITokensMockRecorder) CreateAccessToken(arg0 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateAccessToken", reflect.TypeOf((*MockITokens)(nil).CreateAccessToken), arg0)
}

// CreateRefreshToken mocks base method.
func (m *MockITokens) CreateRefreshToken(arg0 uuid.UUID) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateRefreshToken", arg0)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateRefreshToken indicates an expected call of CreateRefreshToken.
func (mr *MockITokensMockRecorder) CreateRefreshToken(arg0 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateRefreshToken", reflect.TypeOf((*MockITokens)(nil).CreateRefreshToken), arg0)
}

// VerifyToken mocks base method.
func (m *MockITokens) VerifyToken(arg0 string, arg1 *mylogger.Logger) (jwt.MapClaims, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "VerifyToken", arg0, arg1)
	ret0, _ := ret[0].(jwt.MapClaims)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// VerifyToken indicates an expected call of VerifyToken.
func (mr *MockITokensMockRecorder) VerifyToken(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "VerifyToken", reflect.TypeOf((*MockITokens)(nil).VerifyToken), arg0, arg1)
}