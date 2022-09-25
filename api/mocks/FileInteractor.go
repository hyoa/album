// Code generated by mockery v2.14.0. DO NOT EDIT.

package mocks

import mock "github.com/stretchr/testify/mock"

// FileInteractor is an autogenerated mock type for the FileInteractor type
type FileInteractor struct {
	mock.Mock
}

// GetJsonFile provides a mock function with given fields: fileName
func (_m *FileInteractor) GetJsonFile(fileName string) ([]byte, error) {
	ret := _m.Called(fileName)

	var r0 []byte
	if rf, ok := ret.Get(0).(func(string) []byte); ok {
		r0 = rf(fileName)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]byte)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(fileName)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// WriteJsonFile provides a mock function with given fields: fileName, jsonByte
func (_m *FileInteractor) WriteJsonFile(fileName string, jsonByte []byte) error {
	ret := _m.Called(fileName, jsonByte)

	var r0 error
	if rf, ok := ret.Get(0).(func(string, []byte) error); ok {
		r0 = rf(fileName, jsonByte)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

type mockConstructorTestingTNewFileInteractor interface {
	mock.TestingT
	Cleanup(func())
}

// NewFileInteractor creates a new instance of FileInteractor. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewFileInteractor(t mockConstructorTestingTNewFileInteractor) *FileInteractor {
	mock := &FileInteractor{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
