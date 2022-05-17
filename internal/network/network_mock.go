// Code generated by MockGen. DO NOT EDIT.
// Source: network.go

// Package network is a generated GoMock package.
package network

import (
	context "context"
	io "io"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockStreamHandler is a mock of StreamHandler interface.
type MockStreamHandler struct {
	ctrl     *gomock.Controller
	recorder *MockStreamHandlerMockRecorder
}

// MockStreamHandlerMockRecorder is the mock recorder for MockStreamHandler.
type MockStreamHandlerMockRecorder struct {
	mock *MockStreamHandler
}

// NewMockStreamHandler creates a new mock instance.
func NewMockStreamHandler(ctrl *gomock.Controller) *MockStreamHandler {
	mock := &MockStreamHandler{ctrl: ctrl}
	mock.recorder = &MockStreamHandlerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockStreamHandler) EXPECT() *MockStreamHandlerMockRecorder {
	return m.recorder
}

// HandleStream mocks base method.
func (m *MockStreamHandler) HandleStream(ctx context.Context, address string, stream io.ReadWriteCloser) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "HandleStream", ctx, address, stream)
}

// HandleStream indicates an expected call of HandleStream.
func (mr *MockStreamHandlerMockRecorder) HandleStream(ctx, address, stream interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "HandleStream", reflect.TypeOf((*MockStreamHandler)(nil).HandleStream), ctx, address, stream)
}