// Code generated by MockGen. DO NOT EDIT.
// Source: v1_phone_sessions_post.go

// Package mocks_phone_sessions_post is a generated GoMock package.
package mocks_phone_sessions_post

import (
	context "context"
	gomock "github.com/golang/mock/gomock"
	reflect "reflect"
	session "example.com/back/auth/internal/session"
	"net/http"
)

// MockphoneSessionPOSTValidator is a mock of phoneSessionPOSTValidator interface
type MockphoneSessionPOSTValidator struct {
	ctrl     *gomock.Controller
	recorder *MockphoneSessionPOSTValidatorMockRecorder
}

// MockphoneSessionPOSTValidatorMockRecorder is the mock recorder for MockphoneSessionPOSTValidator
type MockphoneSessionPOSTValidatorMockRecorder struct {
	mock *MockphoneSessionPOSTValidator
}

// NewMockphoneSessionPOSTValidator creates a new mock instance
func NewMockphoneSessionPOSTValidator(ctrl *gomock.Controller) *MockphoneSessionPOSTValidator {
	mock := &MockphoneSessionPOSTValidator{ctrl: ctrl}
	mock.recorder = &MockphoneSessionPOSTValidatorMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockphoneSessionPOSTValidator) EXPECT() *MockphoneSessionPOSTValidatorMockRecorder {
	return m.recorder
}

// Struct mocks base method
func (m *MockphoneSessionPOSTValidator) Struct(arg0 interface{}) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Struct", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// Struct indicates an expected call of Struct
func (mr *MockphoneSessionPOSTValidatorMockRecorder) Struct(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Struct", reflect.TypeOf((*MockphoneSessionPOSTValidator)(nil).Struct), arg0)
}

// MockphoneSessionPOSTGetter is a mock of phoneSessionPOSTGetter interface
type MockphoneSessionPOSTGetter struct {
	ctrl     *gomock.Controller
	recorder *MockphoneSessionPOSTGetterMockRecorder
}

// MockphoneSessionPOSTGetterMockRecorder is the mock recorder for MockphoneSessionPOSTGetter
type MockphoneSessionPOSTGetterMockRecorder struct {
	mock *MockphoneSessionPOSTGetter
}

// NewMockphoneSessionPOSTGetter creates a new mock instance
func NewMockphoneSessionPOSTGetter(ctrl *gomock.Controller) *MockphoneSessionPOSTGetter {
	mock := &MockphoneSessionPOSTGetter{ctrl: ctrl}
	mock.recorder = &MockphoneSessionPOSTGetterMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockphoneSessionPOSTGetter) EXPECT() *MockphoneSessionPOSTGetterMockRecorder {
	return m.recorder
}

// FindActiveByUserPhone mocks base method
func (m *MockphoneSessionPOSTGetter) FindActiveByUserPhone(ctx context.Context, usersPhone string) (*session.Session, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindActiveByUserPhone", ctx, usersPhone)
	ret0, _ := ret[0].(*session.Session)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindActiveByUserPhone indicates an expected call of FindActiveByUserPhone
func (mr *MockphoneSessionPOSTGetterMockRecorder) FindActiveByUserPhone(ctx, usersPhone interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindActiveByUserPhone", reflect.TypeOf((*MockphoneSessionPOSTGetter)(nil).FindActiveByUserPhone), ctx, usersPhone)
}

// FindLastActiveByUserPhone mocks base method
func (m *MockphoneSessionPOSTGetter) FindLastActiveByUserPhone(ctx context.Context, usersPhone string) (*session.Session, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindLastActiveByUserPhone", ctx, usersPhone)
	ret0, _ := ret[0].(*session.Session)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindLastActiveByUserPhone indicates an expected call of FindLastActiveByUserPhone
func (mr *MockphoneSessionPOSTGetterMockRecorder) FindLastActiveByUserPhone(ctx, usersPhone interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindLastActiveByUserPhone", reflect.TypeOf((*MockphoneSessionPOSTGetter)(nil).FindLastActiveByUserPhone), ctx, usersPhone)
}

// MockphoneSessionPOSTSaver is a mock of phoneSessionPOSTSaver interface
type MockphoneSessionPOSTSaver struct {
	ctrl     *gomock.Controller
	recorder *MockphoneSessionPOSTSaverMockRecorder
}

// MockphoneSessionPOSTSaverMockRecorder is the mock recorder for MockphoneSessionPOSTSaver
type MockphoneSessionPOSTSaverMockRecorder struct {
	mock *MockphoneSessionPOSTSaver
}

// NewMockphoneSessionPOSTSaver creates a new mock instance
func NewMockphoneSessionPOSTSaver(ctrl *gomock.Controller) *MockphoneSessionPOSTSaver {
	mock := &MockphoneSessionPOSTSaver{ctrl: ctrl}
	mock.recorder = &MockphoneSessionPOSTSaverMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockphoneSessionPOSTSaver) EXPECT() *MockphoneSessionPOSTSaverMockRecorder {
	return m.recorder
}

// Save mocks base method
func (m *MockphoneSessionPOSTSaver) Save(ctx context.Context, usersPhone string, code int64) (*session.Session, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Save", ctx, usersPhone, code)
	ret0, _ := ret[0].(*session.Session)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Save indicates an expected call of Save
func (mr *MockphoneSessionPOSTSaverMockRecorder) Save(ctx, usersPhone, code interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Save", reflect.TypeOf((*MockphoneSessionPOSTSaver)(nil).Save), ctx, usersPhone, code)
}

// DeleteLastByPhone mocks base method
func (m *MockphoneSessionPOSTSaver) DeleteLastByPhone(ctx context.Context, usersPhone string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteLastByPhone", ctx, usersPhone)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteLastByPhone indicates an expected call of DeleteLastByPhone
func (mr *MockphoneSessionPOSTSaverMockRecorder) DeleteLastByPhone(ctx, usersPhone interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteLastByPhone", reflect.TypeOf((*MockphoneSessionPOSTSaver)(nil).DeleteLastByPhone), ctx, usersPhone)
}

// Mocknotifier is a mock of notifier interface
type Mocknotifier struct {
	ctrl     *gomock.Controller
	recorder *MocknotifierMockRecorder
}

// MocknotifierMockRecorder is the mock recorder for Mocknotifier
type MocknotifierMockRecorder struct {
	mock *Mocknotifier
}

// NewMocknotifier creates a new mock instance
func NewMocknotifier(ctrl *gomock.Controller) *Mocknotifier {
	mock := &Mocknotifier{ctrl: ctrl}
	mock.recorder = &MocknotifierMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *Mocknotifier) EXPECT() *MocknotifierMockRecorder {
	return m.recorder
}

// POSTV1UsersSendMessageSMS mocks base method
func (m *Mocknotifier) POSTV1UsersSendMessageSMS(userPhone, message string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "POSTV1UsersSendMessageSMS", userPhone, message)
	ret0, _ := ret[0].(error)
	return ret0
}

// POSTV1UsersSendMessageSMS indicates an expected call of POSTV1UsersSendMessageSMS
func (mr *MocknotifierMockRecorder) POSTV1UsersSendMessageSMS(userPhone, message interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "POSTV1UsersSendMessageSMS", reflect.TypeOf((*Mocknotifier)(nil).POSTV1UsersSendMessageSMS), userPhone, message)
}

// MockvalidStruct is a mock of validStruct interface
type MockvalidStruct struct {
	ctrl     *gomock.Controller
	recorder *MockvalidStructMockRecorder
}

// MockvalidStructMockRecorder is the mock recorder for MockvalidStruct
type MockvalidStructMockRecorder struct {
	mock *MockvalidStruct
}

// NewMockvalidStruct creates a new mock instance
func NewMockvalidStruct(ctrl *gomock.Controller) *MockvalidStruct {
	mock := &MockvalidStruct{ctrl: ctrl}
	mock.recorder = &MockvalidStructMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockvalidStruct) EXPECT() *MockvalidStructMockRecorder {
	return m.recorder
}

// Struct mocks base method
func (m *MockvalidStruct) Struct(arg0 interface{}) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Struct", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// Struct indicates an expected call of Struct
func (mr *MockvalidStructMockRecorder) Struct(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Struct", reflect.TypeOf((*MockvalidStruct)(nil).Struct), arg0)
}

// MockgenerateRandom is a mock of generateRandom interface
type MockgenerateRandom struct {
	ctrl     *gomock.Controller
	recorder *MockgenerateRandomMockRecorder
}

// MockgenerateRandomMockRecorder is the mock recorder for MockgenerateRandom
type MockgenerateRandomMockRecorder struct {
	mock *MockgenerateRandom
}

// NewMockgenerateRandom creates a new mock instance
func NewMockgenerateRandom(ctrl *gomock.Controller) *MockgenerateRandom {
	mock := &MockgenerateRandom{ctrl: ctrl}
	mock.recorder = &MockgenerateRandomMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockgenerateRandom) EXPECT() *MockgenerateRandomMockRecorder {
	return m.recorder
}

// GenerateRandomPassword mocks base method
func (m *MockgenerateRandom) GenerateRandomPassword(length, lenDigits, lenSymbol int) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GenerateRandomPassword", length, lenDigits, lenSymbol)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GenerateRandomPassword indicates an expected call of GenerateRandomPassword
func (mr *MockgenerateRandomMockRecorder) GenerateRandomPassword(length, lenDigits, lenSymbol interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GenerateRandomPassword", reflect.TypeOf((*MockgenerateRandom)(nil).GenerateRandomPassword), length, lenDigits, lenSymbol)
}

// POSTV2SendEmailByUserID mocks base method
func (m *Mocknotifier) POSTV2SendEmailByUserID(ctx context.Context, userID int64, templateName string, tags map[string]string) (*http.Response, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "POSTV2SendEmailByUserID", ctx, userID, templateName, tags)
	ret0, _ := ret[0].(*http.Response)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// POSTV2SendEmailByUserID indicates an expected call of POSTV2SendEmailByUserID
func (mr *MocknotifierMockRecorder) POSTV2SendEmailByUserID(ctx, userID, templateName, tags interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "POSTV2SendEmailByUserID", reflect.TypeOf((*Mocknotifier)(nil).POSTV2SendEmailByUserID), ctx, userID, templateName, tags)
}
