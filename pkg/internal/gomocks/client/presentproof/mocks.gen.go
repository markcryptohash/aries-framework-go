// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/markcryptohash/aries-framework-go/pkg/client/presentproof (interfaces: Provider,ProtocolService)

// Package mocks is a generated GoMock package.
package mocks

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	service "github.com/markcryptohash/aries-framework-go/pkg/didcomm/common/service"
	presentproof "github.com/markcryptohash/aries-framework-go/pkg/didcomm/protocol/presentproof"
)

// MockProvider is a mock of Provider interface.
type MockProvider struct {
	ctrl     *gomock.Controller
	recorder *MockProviderMockRecorder
}

// MockProviderMockRecorder is the mock recorder for MockProvider.
type MockProviderMockRecorder struct {
	mock *MockProvider
}

// NewMockProvider creates a new mock instance.
func NewMockProvider(ctrl *gomock.Controller) *MockProvider {
	mock := &MockProvider{ctrl: ctrl}
	mock.recorder = &MockProviderMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockProvider) EXPECT() *MockProviderMockRecorder {
	return m.recorder
}

// Service mocks base method.
func (m *MockProvider) Service(arg0 string) (interface{}, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Service", arg0)
	ret0, _ := ret[0].(interface{})
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Service indicates an expected call of Service.
func (mr *MockProviderMockRecorder) Service(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Service", reflect.TypeOf((*MockProvider)(nil).Service), arg0)
}

// MockProtocolService is a mock of ProtocolService interface.
type MockProtocolService struct {
	ctrl     *gomock.Controller
	recorder *MockProtocolServiceMockRecorder
}

// MockProtocolServiceMockRecorder is the mock recorder for MockProtocolService.
type MockProtocolServiceMockRecorder struct {
	mock *MockProtocolService
}

// NewMockProtocolService creates a new mock instance.
func NewMockProtocolService(ctrl *gomock.Controller) *MockProtocolService {
	mock := &MockProtocolService{ctrl: ctrl}
	mock.recorder = &MockProtocolServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockProtocolService) EXPECT() *MockProtocolServiceMockRecorder {
	return m.recorder
}

// ActionContinue mocks base method.
func (m *MockProtocolService) ActionContinue(arg0 string, arg1 ...presentproof.Opt) error {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0}
	for _, a := range arg1 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "ActionContinue", varargs...)
	ret0, _ := ret[0].(error)
	return ret0
}

// ActionContinue indicates an expected call of ActionContinue.
func (mr *MockProtocolServiceMockRecorder) ActionContinue(arg0 interface{}, arg1 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0}, arg1...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ActionContinue", reflect.TypeOf((*MockProtocolService)(nil).ActionContinue), varargs...)
}

// ActionStop mocks base method.
func (m *MockProtocolService) ActionStop(arg0 string, arg1 error, arg2 ...presentproof.Opt) error {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0, arg1}
	for _, a := range arg2 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "ActionStop", varargs...)
	ret0, _ := ret[0].(error)
	return ret0
}

// ActionStop indicates an expected call of ActionStop.
func (mr *MockProtocolServiceMockRecorder) ActionStop(arg0, arg1 interface{}, arg2 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0, arg1}, arg2...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ActionStop", reflect.TypeOf((*MockProtocolService)(nil).ActionStop), varargs...)
}

// Actions mocks base method.
func (m *MockProtocolService) Actions() ([]presentproof.Action, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Actions")
	ret0, _ := ret[0].([]presentproof.Action)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Actions indicates an expected call of Actions.
func (mr *MockProtocolServiceMockRecorder) Actions() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Actions", reflect.TypeOf((*MockProtocolService)(nil).Actions))
}

// HandleInbound mocks base method.
func (m *MockProtocolService) HandleInbound(arg0 service.DIDCommMsg, arg1 service.DIDCommContext) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "HandleInbound", arg0, arg1)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// HandleInbound indicates an expected call of HandleInbound.
func (mr *MockProtocolServiceMockRecorder) HandleInbound(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "HandleInbound", reflect.TypeOf((*MockProtocolService)(nil).HandleInbound), arg0, arg1)
}

// HandleOutbound mocks base method.
func (m *MockProtocolService) HandleOutbound(arg0 service.DIDCommMsg, arg1, arg2 string) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "HandleOutbound", arg0, arg1, arg2)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// HandleOutbound indicates an expected call of HandleOutbound.
func (mr *MockProtocolServiceMockRecorder) HandleOutbound(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "HandleOutbound", reflect.TypeOf((*MockProtocolService)(nil).HandleOutbound), arg0, arg1, arg2)
}

// RegisterActionEvent mocks base method.
func (m *MockProtocolService) RegisterActionEvent(arg0 chan<- service.DIDCommAction) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RegisterActionEvent", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// RegisterActionEvent indicates an expected call of RegisterActionEvent.
func (mr *MockProtocolServiceMockRecorder) RegisterActionEvent(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RegisterActionEvent", reflect.TypeOf((*MockProtocolService)(nil).RegisterActionEvent), arg0)
}

// RegisterMsgEvent mocks base method.
func (m *MockProtocolService) RegisterMsgEvent(arg0 chan<- service.StateMsg) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RegisterMsgEvent", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// RegisterMsgEvent indicates an expected call of RegisterMsgEvent.
func (mr *MockProtocolServiceMockRecorder) RegisterMsgEvent(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RegisterMsgEvent", reflect.TypeOf((*MockProtocolService)(nil).RegisterMsgEvent), arg0)
}

// UnregisterActionEvent mocks base method.
func (m *MockProtocolService) UnregisterActionEvent(arg0 chan<- service.DIDCommAction) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UnregisterActionEvent", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// UnregisterActionEvent indicates an expected call of UnregisterActionEvent.
func (mr *MockProtocolServiceMockRecorder) UnregisterActionEvent(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UnregisterActionEvent", reflect.TypeOf((*MockProtocolService)(nil).UnregisterActionEvent), arg0)
}

// UnregisterMsgEvent mocks base method.
func (m *MockProtocolService) UnregisterMsgEvent(arg0 chan<- service.StateMsg) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UnregisterMsgEvent", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// UnregisterMsgEvent indicates an expected call of UnregisterMsgEvent.
func (mr *MockProtocolServiceMockRecorder) UnregisterMsgEvent(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UnregisterMsgEvent", reflect.TypeOf((*MockProtocolService)(nil).UnregisterMsgEvent), arg0)
}
