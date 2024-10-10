// Code generated by mockery v2.43.2. DO NOT EDIT.

package mocks

import mock "github.com/stretchr/testify/mock"

// MailSender is an autogenerated mock type for the MailSender type
type MailSender struct {
	mock.Mock
}

// Send provides a mock function with given fields:
func (_m *MailSender) SendEmail([]string, []byte) error {
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

// NewMailSender creates a new instance of MailSender. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMailSender(t interface {
	mock.TestingT
	Cleanup(func())
}) *MailSender {
	mock := &MailSender{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
