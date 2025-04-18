// Copyright (c) The Jaeger Authors.
// SPDX-License-Identifier: Apache-2.0
//
// Run 'make generate-mocks' to regenerate.

// Code generated by mockery. DO NOT EDIT.

package mocks

import (
	context "context"

	metrics "github.com/jaegertracing/jaeger/internal/proto-gen/api_v2/metrics"
	metricstore "github.com/jaegertracing/jaeger/internal/storage/v1/api/metricstore"

	mock "github.com/stretchr/testify/mock"

	time "time"
)

// Reader is an autogenerated mock type for the Reader type
type Reader struct {
	mock.Mock
}

type Reader_Expecter struct {
	mock *mock.Mock
}

func (_m *Reader) EXPECT() *Reader_Expecter {
	return &Reader_Expecter{mock: &_m.Mock}
}

// GetCallRates provides a mock function with given fields: ctx, params
func (_m *Reader) GetCallRates(ctx context.Context, params *metricstore.CallRateQueryParameters) (*metrics.MetricFamily, error) {
	ret := _m.Called(ctx, params)

	if len(ret) == 0 {
		panic("no return value specified for GetCallRates")
	}

	var r0 *metrics.MetricFamily
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *metricstore.CallRateQueryParameters) (*metrics.MetricFamily, error)); ok {
		return rf(ctx, params)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *metricstore.CallRateQueryParameters) *metrics.MetricFamily); ok {
		r0 = rf(ctx, params)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*metrics.MetricFamily)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *metricstore.CallRateQueryParameters) error); ok {
		r1 = rf(ctx, params)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Reader_GetCallRates_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetCallRates'
type Reader_GetCallRates_Call struct {
	*mock.Call
}

// GetCallRates is a helper method to define mock.On call
//   - ctx context.Context
//   - params *metricstore.CallRateQueryParameters
func (_e *Reader_Expecter) GetCallRates(ctx interface{}, params interface{}) *Reader_GetCallRates_Call {
	return &Reader_GetCallRates_Call{Call: _e.mock.On("GetCallRates", ctx, params)}
}

func (_c *Reader_GetCallRates_Call) Run(run func(ctx context.Context, params *metricstore.CallRateQueryParameters)) *Reader_GetCallRates_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(*metricstore.CallRateQueryParameters))
	})
	return _c
}

func (_c *Reader_GetCallRates_Call) Return(_a0 *metrics.MetricFamily, _a1 error) *Reader_GetCallRates_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *Reader_GetCallRates_Call) RunAndReturn(run func(context.Context, *metricstore.CallRateQueryParameters) (*metrics.MetricFamily, error)) *Reader_GetCallRates_Call {
	_c.Call.Return(run)
	return _c
}

// GetErrorRates provides a mock function with given fields: ctx, params
func (_m *Reader) GetErrorRates(ctx context.Context, params *metricstore.ErrorRateQueryParameters) (*metrics.MetricFamily, error) {
	ret := _m.Called(ctx, params)

	if len(ret) == 0 {
		panic("no return value specified for GetErrorRates")
	}

	var r0 *metrics.MetricFamily
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *metricstore.ErrorRateQueryParameters) (*metrics.MetricFamily, error)); ok {
		return rf(ctx, params)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *metricstore.ErrorRateQueryParameters) *metrics.MetricFamily); ok {
		r0 = rf(ctx, params)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*metrics.MetricFamily)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *metricstore.ErrorRateQueryParameters) error); ok {
		r1 = rf(ctx, params)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Reader_GetErrorRates_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetErrorRates'
type Reader_GetErrorRates_Call struct {
	*mock.Call
}

// GetErrorRates is a helper method to define mock.On call
//   - ctx context.Context
//   - params *metricstore.ErrorRateQueryParameters
func (_e *Reader_Expecter) GetErrorRates(ctx interface{}, params interface{}) *Reader_GetErrorRates_Call {
	return &Reader_GetErrorRates_Call{Call: _e.mock.On("GetErrorRates", ctx, params)}
}

func (_c *Reader_GetErrorRates_Call) Run(run func(ctx context.Context, params *metricstore.ErrorRateQueryParameters)) *Reader_GetErrorRates_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(*metricstore.ErrorRateQueryParameters))
	})
	return _c
}

func (_c *Reader_GetErrorRates_Call) Return(_a0 *metrics.MetricFamily, _a1 error) *Reader_GetErrorRates_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *Reader_GetErrorRates_Call) RunAndReturn(run func(context.Context, *metricstore.ErrorRateQueryParameters) (*metrics.MetricFamily, error)) *Reader_GetErrorRates_Call {
	_c.Call.Return(run)
	return _c
}

// GetLatencies provides a mock function with given fields: ctx, params
func (_m *Reader) GetLatencies(ctx context.Context, params *metricstore.LatenciesQueryParameters) (*metrics.MetricFamily, error) {
	ret := _m.Called(ctx, params)

	if len(ret) == 0 {
		panic("no return value specified for GetLatencies")
	}

	var r0 *metrics.MetricFamily
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *metricstore.LatenciesQueryParameters) (*metrics.MetricFamily, error)); ok {
		return rf(ctx, params)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *metricstore.LatenciesQueryParameters) *metrics.MetricFamily); ok {
		r0 = rf(ctx, params)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*metrics.MetricFamily)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *metricstore.LatenciesQueryParameters) error); ok {
		r1 = rf(ctx, params)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Reader_GetLatencies_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetLatencies'
type Reader_GetLatencies_Call struct {
	*mock.Call
}

// GetLatencies is a helper method to define mock.On call
//   - ctx context.Context
//   - params *metricstore.LatenciesQueryParameters
func (_e *Reader_Expecter) GetLatencies(ctx interface{}, params interface{}) *Reader_GetLatencies_Call {
	return &Reader_GetLatencies_Call{Call: _e.mock.On("GetLatencies", ctx, params)}
}

func (_c *Reader_GetLatencies_Call) Run(run func(ctx context.Context, params *metricstore.LatenciesQueryParameters)) *Reader_GetLatencies_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(*metricstore.LatenciesQueryParameters))
	})
	return _c
}

func (_c *Reader_GetLatencies_Call) Return(_a0 *metrics.MetricFamily, _a1 error) *Reader_GetLatencies_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *Reader_GetLatencies_Call) RunAndReturn(run func(context.Context, *metricstore.LatenciesQueryParameters) (*metrics.MetricFamily, error)) *Reader_GetLatencies_Call {
	_c.Call.Return(run)
	return _c
}

// GetMinStepDuration provides a mock function with given fields: ctx, params
func (_m *Reader) GetMinStepDuration(ctx context.Context, params *metricstore.MinStepDurationQueryParameters) (time.Duration, error) {
	ret := _m.Called(ctx, params)

	if len(ret) == 0 {
		panic("no return value specified for GetMinStepDuration")
	}

	var r0 time.Duration
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *metricstore.MinStepDurationQueryParameters) (time.Duration, error)); ok {
		return rf(ctx, params)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *metricstore.MinStepDurationQueryParameters) time.Duration); ok {
		r0 = rf(ctx, params)
	} else {
		r0 = ret.Get(0).(time.Duration)
	}

	if rf, ok := ret.Get(1).(func(context.Context, *metricstore.MinStepDurationQueryParameters) error); ok {
		r1 = rf(ctx, params)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Reader_GetMinStepDuration_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetMinStepDuration'
type Reader_GetMinStepDuration_Call struct {
	*mock.Call
}

// GetMinStepDuration is a helper method to define mock.On call
//   - ctx context.Context
//   - params *metricstore.MinStepDurationQueryParameters
func (_e *Reader_Expecter) GetMinStepDuration(ctx interface{}, params interface{}) *Reader_GetMinStepDuration_Call {
	return &Reader_GetMinStepDuration_Call{Call: _e.mock.On("GetMinStepDuration", ctx, params)}
}

func (_c *Reader_GetMinStepDuration_Call) Run(run func(ctx context.Context, params *metricstore.MinStepDurationQueryParameters)) *Reader_GetMinStepDuration_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(*metricstore.MinStepDurationQueryParameters))
	})
	return _c
}

func (_c *Reader_GetMinStepDuration_Call) Return(_a0 time.Duration, _a1 error) *Reader_GetMinStepDuration_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *Reader_GetMinStepDuration_Call) RunAndReturn(run func(context.Context, *metricstore.MinStepDurationQueryParameters) (time.Duration, error)) *Reader_GetMinStepDuration_Call {
	_c.Call.Return(run)
	return _c
}

// NewReader creates a new instance of Reader. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewReader(t interface {
	mock.TestingT
	Cleanup(func())
}) *Reader {
	mock := &Reader{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
