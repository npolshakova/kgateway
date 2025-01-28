// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/kgateway-dev/kgateway/pkg/listers (interfaces: NamespaceLister)

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockNamespaceLister is a mock of NamespaceLister interface.
type MockNamespaceLister struct {
	ctrl     *gomock.Controller
	recorder *MockNamespaceListerMockRecorder
}

// MockNamespaceListerMockRecorder is the mock recorder for MockNamespaceLister.
type MockNamespaceListerMockRecorder struct {
	mock *MockNamespaceLister
}

// NewMockNamespaceLister creates a new mock instance.
func NewMockNamespaceLister(ctrl *gomock.Controller) *MockNamespaceLister {
	mock := &MockNamespaceLister{ctrl: ctrl}
	mock.recorder = &MockNamespaceListerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockNamespaceLister) EXPECT() *MockNamespaceListerMockRecorder {
	return m.recorder
}

// List mocks base method.
func (m *MockNamespaceLister) List(arg0 context.Context) ([]string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "List", arg0)
	ret0, _ := ret[0].([]string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// List indicates an expected call of List.
func (mr *MockNamespaceListerMockRecorder) List(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "List", reflect.TypeOf((*MockNamespaceLister)(nil).List), arg0)
}
