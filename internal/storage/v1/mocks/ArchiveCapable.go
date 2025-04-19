// Copyright (c) The Jaeger Authors.
// SPDX-License-Identifier: Apache-2.0
//
// Run 'make generate-mocks' to regenerate.

// Code generated by mockery. DO NOT EDIT.

package mocks

import mock "github.com/stretchr/testify/mock"

// ArchiveCapable is an autogenerated mock type for the ArchiveCapable type
type ArchiveCapable struct {
	mock.Mock
}

type ArchiveCapable_Expecter struct {
	mock *mock.Mock
}

func (_m *ArchiveCapable) EXPECT() *ArchiveCapable_Expecter {
	return &ArchiveCapable_Expecter{mock: &_m.Mock}
}

// IsArchiveCapable provides a mock function with no fields
func (_m *ArchiveCapable) IsArchiveCapable() bool {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for IsArchiveCapable")
	}

	var r0 bool
	if rf, ok := ret.Get(0).(func() bool); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(bool)
	}

	return r0
}

// ArchiveCapable_IsArchiveCapable_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'IsArchiveCapable'
type ArchiveCapable_IsArchiveCapable_Call struct {
	*mock.Call
}

// IsArchiveCapable is a helper method to define mock.On call
func (_e *ArchiveCapable_Expecter) IsArchiveCapable() *ArchiveCapable_IsArchiveCapable_Call {
	return &ArchiveCapable_IsArchiveCapable_Call{Call: _e.mock.On("IsArchiveCapable")}
}

func (_c *ArchiveCapable_IsArchiveCapable_Call) Run(run func()) *ArchiveCapable_IsArchiveCapable_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *ArchiveCapable_IsArchiveCapable_Call) Return(_a0 bool) *ArchiveCapable_IsArchiveCapable_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *ArchiveCapable_IsArchiveCapable_Call) RunAndReturn(run func() bool) *ArchiveCapable_IsArchiveCapable_Call {
	_c.Call.Return(run)
	return _c
}

// NewArchiveCapable creates a new instance of ArchiveCapable. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewArchiveCapable(t interface {
	mock.TestingT
	Cleanup(func())
}) *ArchiveCapable {
	mock := &ArchiveCapable{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
