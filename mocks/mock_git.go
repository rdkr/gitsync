// Code generated by MockGen. DO NOT EDIT.
// Source: gitsync/cmd (interfaces: Git)

// Package mocks is a generated GoMock package.
package mocks

import (
	gomock "github.com/golang/mock/gomock"
	go_git_v4 "gopkg.in/src-d/go-git.v4"
	reflect "reflect"
)

// MockGit is a mock of Git interface
type MockGit struct {
	ctrl     *gomock.Controller
	recorder *MockGitMockRecorder
}

// MockGitMockRecorder is the mock recorder for MockGit
type MockGitMockRecorder struct {
	mock *MockGit
}

// NewMockGit creates a new mock instance
func NewMockGit(ctrl *gomock.Controller) *MockGit {
	mock := &MockGit{ctrl: ctrl}
	mock.recorder = &MockGitMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockGit) EXPECT() *MockGitMockRecorder {
	return m.recorder
}

// Fetch mocks base method
func (m *MockGit) Fetch(arg0 *go_git_v4.Repository) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Fetch", arg0)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Fetch indicates an expected call of Fetch
func (mr *MockGitMockRecorder) Fetch(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Fetch", reflect.TypeOf((*MockGit)(nil).Fetch), arg0)
}

// PlainClone mocks base method
func (m *MockGit) PlainClone() (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "PlainClone")
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// PlainClone indicates an expected call of PlainClone
func (mr *MockGitMockRecorder) PlainClone() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PlainClone", reflect.TypeOf((*MockGit)(nil).PlainClone))
}

// PlainOpen mocks base method
func (m *MockGit) PlainOpen() (*go_git_v4.Repository, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "PlainOpen")
	ret0, _ := ret[0].(*go_git_v4.Repository)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// PlainOpen indicates an expected call of PlainOpen
func (mr *MockGitMockRecorder) PlainOpen() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PlainOpen", reflect.TypeOf((*MockGit)(nil).PlainOpen))
}

// Pull mocks base method
func (m *MockGit) Pull(arg0 *go_git_v4.Worktree) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Pull", arg0)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Pull indicates an expected call of Pull
func (mr *MockGitMockRecorder) Pull(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Pull", reflect.TypeOf((*MockGit)(nil).Pull), arg0)
}