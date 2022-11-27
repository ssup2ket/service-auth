// Code generated by mockery v2.15.0. DO NOT EDIT.

package mocks

import (
	context "context"

	mock "github.com/stretchr/testify/mock"

	token "github.com/ssup2ket/service-auth/pkg/auth/token"
)

// TokenService is an autogenerated mock type for the TokenService type
type TokenService struct {
	mock.Mock
}

// CreateTokens provides a mock function with given fields: ctx, loginID, passwd
func (_m *TokenService) CreateTokens(ctx context.Context, loginID string, passwd string) (*token.TokenInfo, *token.TokenInfo, error) {
	ret := _m.Called(ctx, loginID, passwd)

	var r0 *token.TokenInfo
	if rf, ok := ret.Get(0).(func(context.Context, string, string) *token.TokenInfo); ok {
		r0 = rf(ctx, loginID, passwd)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*token.TokenInfo)
		}
	}

	var r1 *token.TokenInfo
	if rf, ok := ret.Get(1).(func(context.Context, string, string) *token.TokenInfo); ok {
		r1 = rf(ctx, loginID, passwd)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(*token.TokenInfo)
		}
	}

	var r2 error
	if rf, ok := ret.Get(2).(func(context.Context, string, string) error); ok {
		r2 = rf(ctx, loginID, passwd)
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}

// RefreshToken provides a mock function with given fields: ctx, refreshToken
func (_m *TokenService) RefreshToken(ctx context.Context, refreshToken string) (*token.TokenInfo, error) {
	ret := _m.Called(ctx, refreshToken)

	var r0 *token.TokenInfo
	if rf, ok := ret.Get(0).(func(context.Context, string) *token.TokenInfo); ok {
		r0 = rf(ctx, refreshToken)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*token.TokenInfo)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, refreshToken)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

type mockConstructorTestingTNewTokenService interface {
	mock.TestingT
	Cleanup(func())
}

// NewTokenService creates a new instance of TokenService. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewTokenService(t mockConstructorTestingTNewTokenService) *TokenService {
	mock := &TokenService{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
