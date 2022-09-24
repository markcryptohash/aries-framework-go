// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/markcryptohash/aries-framework-go/pkg/didcomm/protocol/issuecredential (interfaces: Provider)

// Package mocks is a generated GoMock package.
package mocks

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	service "github.com/markcryptohash/aries-framework-go/pkg/didcomm/common/service"
	storage "github.com/markcryptohash/aries-framework-go/spi/storage"
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

// Messenger mocks base method.
func (m *MockProvider) Messenger() service.Messenger {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Messenger")
	ret0, _ := ret[0].(service.Messenger)
	return ret0
}

// Messenger indicates an expected call of Messenger.
func (mr *MockProviderMockRecorder) Messenger() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Messenger", reflect.TypeOf((*MockProvider)(nil).Messenger))
}

// StorageProvider mocks base method.
func (m *MockProvider) StorageProvider() storage.Provider {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "StorageProvider")
	ret0, _ := ret[0].(storage.Provider)
	return ret0
}

// StorageProvider indicates an expected call of StorageProvider.
func (mr *MockProviderMockRecorder) StorageProvider() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "StorageProvider", reflect.TypeOf((*MockProvider)(nil).StorageProvider))
}
