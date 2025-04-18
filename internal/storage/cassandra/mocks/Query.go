// Copyright (c) The Jaeger Authors.
// SPDX-License-Identifier: Apache-2.0
//
// Run 'make generate-mocks' to regenerate.

// Code generated by mockery. DO NOT EDIT.

package mocks

import (
	cassandra "github.com/jaegertracing/jaeger/internal/storage/cassandra"
	mock "github.com/stretchr/testify/mock"
)

// Query is an autogenerated mock type for the Query type
type Query struct {
	mock.Mock
}

// Bind provides a mock function with given fields: v
func (_m *Query) Bind(v ...any) cassandra.Query {
	ret := _m.Called(v)

	if len(ret) == 0 {
		panic("no return value specified for Bind")
	}

	var r0 cassandra.Query
	if rf, ok := ret.Get(0).(func(...any) cassandra.Query); ok {
		r0 = rf(v...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(cassandra.Query)
		}
	}

	return r0
}

// Consistency provides a mock function with given fields: level
func (_m *Query) Consistency(level cassandra.Consistency) cassandra.Query {
	ret := _m.Called(level)

	if len(ret) == 0 {
		panic("no return value specified for Consistency")
	}

	var r0 cassandra.Query
	if rf, ok := ret.Get(0).(func(cassandra.Consistency) cassandra.Query); ok {
		r0 = rf(level)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(cassandra.Query)
		}
	}

	return r0
}

// Exec provides a mock function with no fields
func (_m *Query) Exec() error {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for Exec")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func() error); ok {
		r0 = rf()
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Iter provides a mock function with no fields
func (_m *Query) Iter() cassandra.Iterator {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for Iter")
	}

	var r0 cassandra.Iterator
	if rf, ok := ret.Get(0).(func() cassandra.Iterator); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(cassandra.Iterator)
		}
	}

	return r0
}

// PageSize provides a mock function with given fields: _a0
func (_m *Query) PageSize(_a0 int) cassandra.Query {
	ret := _m.Called(_a0)

	if len(ret) == 0 {
		panic("no return value specified for PageSize")
	}

	var r0 cassandra.Query
	if rf, ok := ret.Get(0).(func(int) cassandra.Query); ok {
		r0 = rf(_a0)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(cassandra.Query)
		}
	}

	return r0
}

// ScanCAS provides a mock function with given fields: dest
func (_m *Query) ScanCAS(dest ...any) (bool, error) {
	ret := _m.Called(dest)

	if len(ret) == 0 {
		panic("no return value specified for ScanCAS")
	}

	var r0 bool
	var r1 error
	if rf, ok := ret.Get(0).(func(...any) (bool, error)); ok {
		return rf(dest...)
	}
	if rf, ok := ret.Get(0).(func(...any) bool); ok {
		r0 = rf(dest...)
	} else {
		r0 = ret.Get(0).(bool)
	}

	if rf, ok := ret.Get(1).(func(...any) error); ok {
		r1 = rf(dest...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// String provides a mock function with no fields
func (_m *Query) String() string {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for String")
	}

	var r0 string
	if rf, ok := ret.Get(0).(func() string); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}

// NewQuery creates a new instance of Query. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewQuery(t interface {
	mock.TestingT
	Cleanup(func())
}) *Query {
	mock := &Query{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
