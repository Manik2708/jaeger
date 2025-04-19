// Copyright (c) The Jaeger Authors.
// SPDX-License-Identifier: Apache-2.0
//
// Run 'make generate-mocks' to regenerate.

// Code generated by mockery. DO NOT EDIT.

package mocks

import mock "github.com/stretchr/testify/mock"

// ServiceAliasMappingExternalSource is an autogenerated mock type for the ServiceAliasMappingExternalSource type
type ServiceAliasMappingExternalSource struct {
	mock.Mock
}

type ServiceAliasMappingExternalSource_Expecter struct {
	mock *mock.Mock
}

func (_m *ServiceAliasMappingExternalSource) EXPECT() *ServiceAliasMappingExternalSource_Expecter {
	return &ServiceAliasMappingExternalSource_Expecter{mock: &_m.Mock}
}

// Load provides a mock function with no fields
func (_m *ServiceAliasMappingExternalSource) Load() (map[string]string, error) {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for Load")
	}

	var r0 map[string]string
	var r1 error
	if rf, ok := ret.Get(0).(func() (map[string]string, error)); ok {
		return rf()
	}
	if rf, ok := ret.Get(0).(func() map[string]string); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(map[string]string)
		}
	}

	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ServiceAliasMappingExternalSource_Load_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Load'
type ServiceAliasMappingExternalSource_Load_Call struct {
	*mock.Call
}

// Load is a helper method to define mock.On call
func (_e *ServiceAliasMappingExternalSource_Expecter) Load() *ServiceAliasMappingExternalSource_Load_Call {
	return &ServiceAliasMappingExternalSource_Load_Call{Call: _e.mock.On("Load")}
}

func (_c *ServiceAliasMappingExternalSource_Load_Call) Run(run func()) *ServiceAliasMappingExternalSource_Load_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *ServiceAliasMappingExternalSource_Load_Call) Return(_a0 map[string]string, _a1 error) *ServiceAliasMappingExternalSource_Load_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *ServiceAliasMappingExternalSource_Load_Call) RunAndReturn(run func() (map[string]string, error)) *ServiceAliasMappingExternalSource_Load_Call {
	_c.Call.Return(run)
	return _c
}

// NewServiceAliasMappingExternalSource creates a new instance of ServiceAliasMappingExternalSource. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewServiceAliasMappingExternalSource(t interface {
	mock.TestingT
	Cleanup(func())
}) *ServiceAliasMappingExternalSource {
	mock := &ServiceAliasMappingExternalSource{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
