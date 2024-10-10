// Code generated by mockery v2.43.2. DO NOT EDIT.

package mocks

import mock "github.com/stretchr/testify/mock"

// MailClient is an autogenerated mock type for the MailClient type
type MailClient struct {
	mock.Mock
}

// Send provides a mock function with given fields:
func (_m *MailClient) SendEmail([]string, []byte) error {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for SendEmail")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func() error); ok {
		r0 = rf()
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// NewMailClient creates a new instance of MailClient. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMailClient(t interface {
	mock.TestingT
	Cleanup(func())
}) *MailClient {
	mock := &MailClient{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
