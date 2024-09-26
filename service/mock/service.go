// Code generated by MockGen. DO NOT EDIT.
// Source: service.go
//
// Generated by this command:
//
//	mockgen -source=service.go -destination=mock/service.go -package=mock
//

// Package mock is a generated GoMock package.
package mock

import (
	context "context"
	reflect "reflect"

	user "github.com/Mitra-Apps/be-user-service/domain/proto/user"
	entity "github.com/Mitra-Apps/be-user-service/domain/user/entity"
	uuid "github.com/google/uuid"
	gomock "go.uber.org/mock/gomock"
)

// MockServiceInterface is a mock of ServiceInterface interface.
type MockServiceInterface struct {
	ctrl     *gomock.Controller
	recorder *MockServiceInterfaceMockRecorder
}

// MockServiceInterfaceMockRecorder is the mock recorder for MockServiceInterface.
type MockServiceInterfaceMockRecorder struct {
	mock *MockServiceInterface
}

// NewMockServiceInterface creates a new mock instance.
func NewMockServiceInterface(ctrl *gomock.Controller) *MockServiceInterface {
	mock := &MockServiceInterface{ctrl: ctrl}
	mock.recorder = &MockServiceInterfaceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockServiceInterface) EXPECT() *MockServiceInterfaceMockRecorder {
	return m.recorder
}

// ChangePassword mocks base method.
func (m *MockServiceInterface) ChangePassword(ctx context.Context, req *user.ChangePasswordRequest) (*entity.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ChangePassword", ctx, req)
	ret0, _ := ret[0].(*entity.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ChangePassword indicates an expected call of ChangePassword.
func (mr *MockServiceInterfaceMockRecorder) ChangePassword(ctx, req any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ChangePassword", reflect.TypeOf((*MockServiceInterface)(nil).ChangePassword), ctx, req)
}

// CreateRole mocks base method.
func (m *MockServiceInterface) CreateRole(ctx context.Context, role *entity.Role) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateRole", ctx, role)
	ret0, _ := ret[0].(error)
	return ret0
}

// CreateRole indicates an expected call of CreateRole.
func (mr *MockServiceInterfaceMockRecorder) CreateRole(ctx, role any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateRole", reflect.TypeOf((*MockServiceInterface)(nil).CreateRole), ctx, role)
}

// GetAll mocks base method.
func (m *MockServiceInterface) GetAll(ctx context.Context) ([]*entity.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAll", ctx)
	ret0, _ := ret[0].([]*entity.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAll indicates an expected call of GetAll.
func (mr *MockServiceInterfaceMockRecorder) GetAll(ctx any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAll", reflect.TypeOf((*MockServiceInterface)(nil).GetAll), ctx)
}

// GetByID mocks base method.
func (m *MockServiceInterface) GetByID(ctx context.Context, id uuid.UUID) (*entity.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetByID", ctx, id)
	ret0, _ := ret[0].(*entity.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetByID indicates an expected call of GetByID.
func (mr *MockServiceInterfaceMockRecorder) GetByID(ctx, id any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByID", reflect.TypeOf((*MockServiceInterface)(nil).GetByID), ctx, id)
}

// GetRole mocks base method.
func (m *MockServiceInterface) GetRole(ctx context.Context) ([]entity.Role, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetRole", ctx)
	ret0, _ := ret[0].([]entity.Role)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetRole indicates an expected call of GetRole.
func (mr *MockServiceInterfaceMockRecorder) GetRole(ctx any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetRole", reflect.TypeOf((*MockServiceInterface)(nil).GetRole), ctx)
}

// Login mocks base method.
func (m *MockServiceInterface) Login(ctx context.Context, payload entity.LoginRequest) (*entity.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Login", ctx, payload)
	ret0, _ := ret[0].(*entity.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Login indicates an expected call of Login.
func (mr *MockServiceInterfaceMockRecorder) Login(ctx, payload any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Login", reflect.TypeOf((*MockServiceInterface)(nil).Login), ctx, payload)
}

// Logout mocks base method.
func (m *MockServiceInterface) Logout(ctx context.Context, id uuid.UUID) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Logout", ctx, id)
	ret0, _ := ret[0].(error)
	return ret0
}

// Logout indicates an expected call of Logout.
func (mr *MockServiceInterfaceMockRecorder) Logout(ctx, id any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Logout", reflect.TypeOf((*MockServiceInterface)(nil).Logout), ctx, id)
}

// Register mocks base method.
func (m *MockServiceInterface) Register(ctx context.Context, req *user.UserRegisterRequest) (*entity.OtpMailReq, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Register", ctx, req)
	ret0, _ := ret[0].(*entity.OtpMailReq)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Register indicates an expected call of Register.
func (mr *MockServiceInterfaceMockRecorder) Register(ctx, req any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Register", reflect.TypeOf((*MockServiceInterface)(nil).Register), ctx, req)
}

// ResendOTP mocks base method.
func (m *MockServiceInterface) ResendOTP(ctx context.Context, email string) (*entity.OtpMailReq, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ResendOTP", ctx, email)
	ret0, _ := ret[0].(*entity.OtpMailReq)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ResendOTP indicates an expected call of ResendOTP.
func (mr *MockServiceInterfaceMockRecorder) ResendOTP(ctx, email any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ResendOTP", reflect.TypeOf((*MockServiceInterface)(nil).ResendOTP), ctx, email)
}

// Save mocks base method.
func (m *MockServiceInterface) Save(ctx context.Context, user *entity.User) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Save", ctx, user)
	ret0, _ := ret[0].(error)
	return ret0
}

// Save indicates an expected call of Save.
func (mr *MockServiceInterfaceMockRecorder) Save(ctx, user any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Save", reflect.TypeOf((*MockServiceInterface)(nil).Save), ctx, user)
}

// VerifyOTP mocks base method.
func (m *MockServiceInterface) VerifyOTP(ctx context.Context, otp int, redisKey string) (*entity.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "VerifyOTP", ctx, otp, redisKey)
	ret0, _ := ret[0].(*entity.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// VerifyOTP indicates an expected call of VerifyOTP.
func (mr *MockServiceInterfaceMockRecorder) VerifyOTP(ctx, otp, redisKey any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "VerifyOTP", reflect.TypeOf((*MockServiceInterface)(nil).VerifyOTP), ctx, otp, redisKey)
}
