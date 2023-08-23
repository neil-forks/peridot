// The MIT License
//
// Copyright (c) 2020 Temporal Technologies Inc.  All rights reserved.
//
// Copyright (c) 2020 Uber Technologies, Inc.
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

// Code generated by MockGen. DO NOT EDIT.
// Source: task_executor.go

// Package replication is a generated GoMock package.
package replication

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	repication "go.temporal.io/server/api/replication/v1"
)

// MockTaskExecutor is a mock of TaskExecutor interface.
type MockTaskExecutor struct {
	ctrl     *gomock.Controller
	recorder *MockTaskExecutorMockRecorder
}

// MockTaskExecutorMockRecorder is the mock recorder for MockTaskExecutor.
type MockTaskExecutorMockRecorder struct {
	mock *MockTaskExecutor
}

// NewMockTaskExecutor creates a new mock instance.
func NewMockTaskExecutor(ctrl *gomock.Controller) *MockTaskExecutor {
	mock := &MockTaskExecutor{ctrl: ctrl}
	mock.recorder = &MockTaskExecutorMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockTaskExecutor) EXPECT() *MockTaskExecutorMockRecorder {
	return m.recorder
}

// Execute mocks base method.
func (m *MockTaskExecutor) Execute(ctx context.Context, replicationTask *repication.ReplicationTask, forceApply bool) (int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Execute", ctx, replicationTask, forceApply)
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Execute indicates an expected call of Execute.
func (mr *MockTaskExecutorMockRecorder) Execute(ctx, replicationTask, forceApply interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Execute", reflect.TypeOf((*MockTaskExecutor)(nil).Execute), ctx, replicationTask, forceApply)
}