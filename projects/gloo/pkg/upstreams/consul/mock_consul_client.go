// Code generated by MockGen. DO NOT EDIT.
// Source: consul_client.go

// Package consul is a generated GoMock package.
package consul

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	api "github.com/hashicorp/consul/api"
)

// MockConsulClient is a mock of ConsulClient interface
type MockConsulClient struct {
	ctrl     *gomock.Controller
	recorder *MockConsulClientMockRecorder
}

// MockConsulClientMockRecorder is the mock recorder for MockConsulClient
type MockConsulClientMockRecorder struct {
	mock *MockConsulClient
}

// NewMockConsulClient creates a new mock instance
func NewMockConsulClient(ctrl *gomock.Controller) *MockConsulClient {
	mock := &MockConsulClient{ctrl: ctrl}
	mock.recorder = &MockConsulClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockConsulClient) EXPECT() *MockConsulClientMockRecorder {
	return m.recorder
}

// CanConnect mocks base method
func (m *MockConsulClient) CanConnect() bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CanConnect")
	ret0, _ := ret[0].(bool)
	return ret0
}

// CanConnect indicates an expected call of CanConnect
func (mr *MockConsulClientMockRecorder) CanConnect() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CanConnect", reflect.TypeOf((*MockConsulClient)(nil).CanConnect))
}

// DataCenters mocks base method
func (m *MockConsulClient) DataCenters() ([]string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DataCenters")
	ret0, _ := ret[0].([]string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// DataCenters indicates an expected call of DataCenters
func (mr *MockConsulClientMockRecorder) DataCenters() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DataCenters", reflect.TypeOf((*MockConsulClient)(nil).DataCenters))
}

// Services mocks base method
func (m *MockConsulClient) Services(q *api.QueryOptions) (map[string][]string, *api.QueryMeta, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Services", q)
	ret0, _ := ret[0].(map[string][]string)
	ret1, _ := ret[1].(*api.QueryMeta)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// Services indicates an expected call of Services
func (mr *MockConsulClientMockRecorder) Services(q interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Services", reflect.TypeOf((*MockConsulClient)(nil).Services), q)
}
