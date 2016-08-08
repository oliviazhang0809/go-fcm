package gcm

import "github.com/stretchr/testify/mock"

// httpClientMock is an autogenerated mock type for the httpClient type
type httpClientMock struct {
	mock.Mock
}

// getRetryAfter provides a mock function with given fields:
func (_m *httpClientMock) getRetryAfter() string {
	ret := _m.Called()

	var r0 string
	if rf, ok := ret.Get(0).(func() string); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}

// send provides a mock function with given fields: m
func (_m *httpClientMock) send(m HTTPMessage) (*HTTPResponse, error) {
	ret := _m.Called(m)

	var r0 *HTTPResponse
	if rf, ok := ret.Get(0).(func(HTTPMessage) *HTTPResponse); ok {
		r0 = rf(m)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*HTTPResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(HTTPMessage) error); ok {
		r1 = rf(m)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}