// Code generated by mockery v2.14.0. DO NOT EDIT.

package mocks

import mock "github.com/stretchr/testify/mock"

// Storage is an autogenerated mock type for the Storage type
type Storage struct {
	mock.Mock
}

// MediaExist provides a mock function with given fields: key, bucket
func (_m *Storage) MediaExist(key string, bucket string) (bool, error) {
	ret := _m.Called(key, bucket)

	var r0 bool
	if rf, ok := ret.Get(0).(func(string, string) bool); ok {
		r0 = rf(key, bucket)
	} else {
		r0 = ret.Get(0).(bool)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string, string) error); ok {
		r1 = rf(key, bucket)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// SignDownloadUri provides a mock function with given fields: key, bucket
func (_m *Storage) SignDownloadUri(key string, bucket string) (string, error) {
	ret := _m.Called(key, bucket)

	var r0 string
	if rf, ok := ret.Get(0).(func(string, string) string); ok {
		r0 = rf(key, bucket)
	} else {
		r0 = ret.Get(0).(string)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string, string) error); ok {
		r1 = rf(key, bucket)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// SignUploadUri provides a mock function with given fields: key, bucket
func (_m *Storage) SignUploadUri(key string, bucket string) (string, error) {
	ret := _m.Called(key, bucket)

	var r0 string
	if rf, ok := ret.Get(0).(func(string, string) string); ok {
		r0 = rf(key, bucket)
	} else {
		r0 = ret.Get(0).(string)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string, string) error); ok {
		r1 = rf(key, bucket)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

type mockConstructorTestingTNewStorage interface {
	mock.TestingT
	Cleanup(func())
}

// NewStorage creates a new instance of Storage. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewStorage(t mockConstructorTestingTNewStorage) *Storage {
	mock := &Storage{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
