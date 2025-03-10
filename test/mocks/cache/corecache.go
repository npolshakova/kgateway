//go:build ignore

// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/solo-io/solo-kit/pkg/api/v1/clients/kube/cache (interfaces: KubeCoreCache)

// Package mock_cache is a generated GoMock package.
package mock_cache

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	cache "github.com/solo-io/solo-kit/pkg/api/v1/clients/kube/cache"
	v1 "k8s.io/client-go/listers/core/v1"
)

// MockKubeCoreCache is a mock of KubeCoreCache interface.
type MockKubeCoreCache struct {
	ctrl     *gomock.Controller
	recorder *MockKubeCoreCacheMockRecorder
}

// MockKubeCoreCacheMockRecorder is the mock recorder for MockKubeCoreCache.
type MockKubeCoreCacheMockRecorder struct {
	mock *MockKubeCoreCache
}

// NewMockKubeCoreCache creates a new mock instance.
func NewMockKubeCoreCache(ctrl *gomock.Controller) *MockKubeCoreCache {
	mock := &MockKubeCoreCache{ctrl: ctrl}
	mock.recorder = &MockKubeCoreCacheMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockKubeCoreCache) EXPECT() *MockKubeCoreCacheMockRecorder {
	return m.recorder
}

// ConfigMapLister mocks base method.
func (m *MockKubeCoreCache) ConfigMapLister() v1.ConfigMapLister {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ConfigMapLister")
	ret0, _ := ret[0].(v1.ConfigMapLister)
	return ret0
}

// ConfigMapLister indicates an expected call of ConfigMapLister.
func (mr *MockKubeCoreCacheMockRecorder) ConfigMapLister() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ConfigMapLister", reflect.TypeOf((*MockKubeCoreCache)(nil).ConfigMapLister))
}

// NamespaceLister mocks base method.
func (m *MockKubeCoreCache) NamespaceLister() v1.NamespaceLister {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "NamespaceLister")
	ret0, _ := ret[0].(v1.NamespaceLister)
	return ret0
}

// NamespaceLister indicates an expected call of NamespaceLister.
func (mr *MockKubeCoreCacheMockRecorder) NamespaceLister() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "NamespaceLister", reflect.TypeOf((*MockKubeCoreCache)(nil).NamespaceLister))
}

// NamespacedConfigMapLister mocks base method.
func (m *MockKubeCoreCache) NamespacedConfigMapLister(arg0 string) cache.ConfigMapLister {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "NamespacedConfigMapLister", arg0)
	ret0, _ := ret[0].(cache.ConfigMapLister)
	return ret0
}

// NamespacedConfigMapLister indicates an expected call of NamespacedConfigMapLister.
func (mr *MockKubeCoreCacheMockRecorder) NamespacedConfigMapLister(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "NamespacedConfigMapLister", reflect.TypeOf((*MockKubeCoreCache)(nil).NamespacedConfigMapLister), arg0)
}

// NamespacedPodLister mocks base method.
func (m *MockKubeCoreCache) NamespacedPodLister(arg0 string) cache.PodLister {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "NamespacedPodLister", arg0)
	ret0, _ := ret[0].(cache.PodLister)
	return ret0
}

// NamespacedPodLister indicates an expected call of NamespacedPodLister.
func (mr *MockKubeCoreCacheMockRecorder) NamespacedPodLister(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "NamespacedPodLister", reflect.TypeOf((*MockKubeCoreCache)(nil).NamespacedPodLister), arg0)
}

// NamespacedSecretLister mocks base method.
func (m *MockKubeCoreCache) NamespacedSecretLister(arg0 string) cache.SecretLister {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "NamespacedSecretLister", arg0)
	ret0, _ := ret[0].(cache.SecretLister)
	return ret0
}

// NamespacedSecretLister indicates an expected call of NamespacedSecretLister.
func (mr *MockKubeCoreCacheMockRecorder) NamespacedSecretLister(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "NamespacedSecretLister", reflect.TypeOf((*MockKubeCoreCache)(nil).NamespacedSecretLister), arg0)
}

// NamespacedServiceLister mocks base method.
func (m *MockKubeCoreCache) NamespacedServiceLister(arg0 string) cache.ServiceLister {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "NamespacedServiceLister", arg0)
	ret0, _ := ret[0].(cache.ServiceLister)
	return ret0
}

// NamespacedServiceLister indicates an expected call of NamespacedServiceLister.
func (mr *MockKubeCoreCacheMockRecorder) NamespacedServiceLister(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "NamespacedServiceLister", reflect.TypeOf((*MockKubeCoreCache)(nil).NamespacedServiceLister), arg0)
}

// PodLister mocks base method.
func (m *MockKubeCoreCache) PodLister() v1.PodLister {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "PodLister")
	ret0, _ := ret[0].(v1.PodLister)
	return ret0
}

// PodLister indicates an expected call of PodLister.
func (mr *MockKubeCoreCacheMockRecorder) PodLister() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PodLister", reflect.TypeOf((*MockKubeCoreCache)(nil).PodLister))
}

// SecretLister mocks base method.
func (m *MockKubeCoreCache) SecretLister() v1.SecretLister {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SecretLister")
	ret0, _ := ret[0].(v1.SecretLister)
	return ret0
}

// SecretLister indicates an expected call of SecretLister.
func (mr *MockKubeCoreCacheMockRecorder) SecretLister() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SecretLister", reflect.TypeOf((*MockKubeCoreCache)(nil).SecretLister))
}

// ServiceLister mocks base method.
func (m *MockKubeCoreCache) ServiceLister() v1.ServiceLister {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ServiceLister")
	ret0, _ := ret[0].(v1.ServiceLister)
	return ret0
}

// ServiceLister indicates an expected call of ServiceLister.
func (mr *MockKubeCoreCacheMockRecorder) ServiceLister() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ServiceLister", reflect.TypeOf((*MockKubeCoreCache)(nil).ServiceLister))
}

// Subscribe mocks base method.
func (m *MockKubeCoreCache) Subscribe() <-chan struct{} {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Subscribe")
	ret0, _ := ret[0].(<-chan struct{})
	return ret0
}

// Subscribe indicates an expected call of Subscribe.
func (mr *MockKubeCoreCacheMockRecorder) Subscribe() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Subscribe", reflect.TypeOf((*MockKubeCoreCache)(nil).Subscribe))
}

// Unsubscribe mocks base method.
func (m *MockKubeCoreCache) Unsubscribe(arg0 <-chan struct{}) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Unsubscribe", arg0)
}

// Unsubscribe indicates an expected call of Unsubscribe.
func (mr *MockKubeCoreCacheMockRecorder) Unsubscribe(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Unsubscribe", reflect.TypeOf((*MockKubeCoreCache)(nil).Unsubscribe), arg0)
}
