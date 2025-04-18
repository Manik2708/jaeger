// Copyright (c) The Jaeger Authors.
// SPDX-License-Identifier: Apache-2.0
//
// Run 'make generate-mocks' to regenerate.

// Code generated by mockery. DO NOT EDIT.

package mocks

import mock "github.com/stretchr/testify/mock"

// CollectorService is an autogenerated mock type for the CollectorService type
type CollectorService struct {
	mock.Mock
}

type CollectorService_Expecter struct {
	mock *mock.Mock
}

func (_m *CollectorService) EXPECT() *CollectorService_Expecter {
	return &CollectorService_Expecter{mock: &_m.Mock}
}

// GetSamplingRate provides a mock function with given fields: service, operation
func (_m *CollectorService) GetSamplingRate(service string, operation string) (float64, error) {
	ret := _m.Called(service, operation)

	if len(ret) == 0 {
		panic("no return value specified for GetSamplingRate")
	}

	var r0 float64
	var r1 error
	if rf, ok := ret.Get(0).(func(string, string) (float64, error)); ok {
		return rf(service, operation)
	}
	if rf, ok := ret.Get(0).(func(string, string) float64); ok {
		r0 = rf(service, operation)
	} else {
		r0 = ret.Get(0).(float64)
	}

	if rf, ok := ret.Get(1).(func(string, string) error); ok {
		r1 = rf(service, operation)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// CollectorService_GetSamplingRate_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetSamplingRate'
type CollectorService_GetSamplingRate_Call struct {
	*mock.Call
}

// GetSamplingRate is a helper method to define mock.On call
//   - service string
//   - operation string
func (_e *CollectorService_Expecter) GetSamplingRate(service interface{}, operation interface{}) *CollectorService_GetSamplingRate_Call {
	return &CollectorService_GetSamplingRate_Call{Call: _e.mock.On("GetSamplingRate", service, operation)}
}

func (_c *CollectorService_GetSamplingRate_Call) Run(run func(service string, operation string)) *CollectorService_GetSamplingRate_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string), args[1].(string))
	})
	return _c
}

func (_c *CollectorService_GetSamplingRate_Call) Return(_a0 float64, _a1 error) *CollectorService_GetSamplingRate_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *CollectorService_GetSamplingRate_Call) RunAndReturn(run func(string, string) (float64, error)) *CollectorService_GetSamplingRate_Call {
	_c.Call.Return(run)
	return _c
}

// NewCollectorService creates a new instance of CollectorService. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewCollectorService(t interface {
	mock.TestingT
	Cleanup(func())
}) *CollectorService {
	mock := &CollectorService{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
