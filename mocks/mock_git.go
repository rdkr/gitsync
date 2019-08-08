// Code generated by MockGen. DO NOT EDIT.
// Source: gitsync/cmd (interfaces: Cloner)

// Package mocks is a generated GoMock package.
package mocks

import (
	gomock "github.com/golang/mock/gomock"
	cmd "gitsync/cmd"
	go_git_v4 "gopkg.in/src-d/go-git.v4"
	reflect "reflect"
)

// MockCloner is a mock of Cloner interface
type MockCloner struct {
	ctrl     *gomock.Controller
	recorder *MockClonerMockRecorder
}

// MockClonerMockRecorder is the mock recorder for MockCloner
type MockClonerMockRecorder struct {
	mock *MockCloner
}

// NewMockCloner creates a new mock instance
func NewMockCloner(ctrl *gomock.Controller) *MockCloner {
	mock := &MockCloner{ctrl: ctrl}
	mock.recorder = &MockClonerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockCloner) EXPECT() *MockClonerMockRecorder {
	return m.recorder
}

// Fetch mocks base method
func (m *MockCloner) Fetch(arg0 *go_git_v4.Repository) cmd.Status {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Fetch", arg0)
	ret0, _ := ret[0].(cmd.Status)
	return ret0
}

// Fetch indicates an expected call of Fetch
func (mr *MockClonerMockRecorder) Fetch(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Fetch", reflect.TypeOf((*MockCloner)(nil).Fetch), arg0)
}

// PlainClone mocks base method
func (m *MockCloner) PlainClone() cmd.Status {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "PlainClone")
	ret0, _ := ret[0].(cmd.Status)
	return ret0
}

// PlainClone indicates an expected call of PlainClone
func (mr *MockClonerMockRecorder) PlainClone() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PlainClone", reflect.TypeOf((*MockCloner)(nil).PlainClone))
}

// PlainOpen mocks base method
func (m *MockCloner) PlainOpen() (*go_git_v4.Repository, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "PlainOpen")
	ret0, _ := ret[0].(*go_git_v4.Repository)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// PlainOpen indicates an expected call of PlainOpen
func (mr *MockClonerMockRecorder) PlainOpen() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PlainOpen", reflect.TypeOf((*MockCloner)(nil).PlainOpen))
}

// Pull mocks base method
func (m *MockCloner) Pull(arg0 *go_git_v4.Worktree) cmd.Status {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Pull", arg0)
	ret0, _ := ret[0].(cmd.Status)
	return ret0
}

// Pull indicates an expected call of Pull
func (mr *MockClonerMockRecorder) Pull(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Pull", reflect.TypeOf((*MockCloner)(nil).Pull), arg0)
}
