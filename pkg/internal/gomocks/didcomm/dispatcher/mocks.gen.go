// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/markcryptohash/aries-framework-go/pkg/didcomm/dispatcher (interfaces: Outbound)

// Package mocks is a generated GoMock package.
package mocks

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	service "github.com/markcryptohash/aries-framework-go/pkg/didcomm/common/service"
)

// MockOutbound is a mock of Outbound interface.
type MockOutbound struct {
	ctrl     *gomock.Controller
	recorder *MockOutboundMockRecorder
}

// MockOutboundMockRecorder is the mock recorder for MockOutbound.
type MockOutboundMockRecorder struct {
	mock *MockOutbound
}

// NewMockOutbound creates a new mock instance.
func NewMockOutbound(ctrl *gomock.Controller) *MockOutbound {
	mock := &MockOutbound{ctrl: ctrl}
	mock.recorder = &MockOutboundMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockOutbound) EXPECT() *MockOutboundMockRecorder {
	return m.recorder
}

// Forward mocks base method.
func (m *MockOutbound) Forward(arg0 interface{}, arg1 *service.Destination) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Forward", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// Forward indicates an expected call of Forward.
func (mr *MockOutboundMockRecorder) Forward(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Forward", reflect.TypeOf((*MockOutbound)(nil).Forward), arg0, arg1)
}

// Send mocks base method.
func (m *MockOutbound) Send(arg0 interface{}, arg1 string, arg2 *service.Destination) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Send", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// Send indicates an expected call of Send.
func (mr *MockOutboundMockRecorder) Send(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Send", reflect.TypeOf((*MockOutbound)(nil).Send), arg0, arg1, arg2)
}

// SendToDID mocks base method.
func (m *MockOutbound) SendToDID(arg0 interface{}, arg1, arg2 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SendToDID", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// SendToDID indicates an expected call of SendToDID.
func (mr *MockOutboundMockRecorder) SendToDID(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SendToDID", reflect.TypeOf((*MockOutbound)(nil).SendToDID), arg0, arg1, arg2)
}
