// Code generated by mockery v0.0.0-dev. DO NOT EDIT.

package mocks

import (
	context "context"

	model "github.com/ssup2ket/ssup2ket-auth-service/internal/domain/model"
	mock "github.com/stretchr/testify/mock"

	uuidmodel "github.com/ssup2ket/ssup2ket-auth-service/pkg/uuidmodel"
)

// UserService is an autogenerated mock type for the UserService type
type UserService struct {
	mock.Mock
}

// CreateUser provides a mock function with given fields: ctx, userInfo, passwd
func (_m *UserService) CreateUser(ctx context.Context, userInfo *model.UserInfo, passwd string) (*model.UserInfo, error) {
	ret := _m.Called(ctx, userInfo, passwd)

	var r0 *model.UserInfo
	if rf, ok := ret.Get(0).(func(context.Context, *model.UserInfo, string) *model.UserInfo); ok {
		r0 = rf(ctx, userInfo, passwd)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.UserInfo)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *model.UserInfo, string) error); ok {
		r1 = rf(ctx, userInfo, passwd)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// DeleteUser provides a mock function with given fields: ctx, userUUID
func (_m *UserService) DeleteUser(ctx context.Context, userUUID uuidmodel.UUIDModel) error {
	ret := _m.Called(ctx, userUUID)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, uuidmodel.UUIDModel) error); ok {
		r0 = rf(ctx, userUUID)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// GetUser provides a mock function with given fields: ctx, userUUID
func (_m *UserService) GetUser(ctx context.Context, userUUID uuidmodel.UUIDModel) (*model.UserInfo, error) {
	ret := _m.Called(ctx, userUUID)

	var r0 *model.UserInfo
	if rf, ok := ret.Get(0).(func(context.Context, uuidmodel.UUIDModel) *model.UserInfo); ok {
		r0 = rf(ctx, userUUID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.UserInfo)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, uuidmodel.UUIDModel) error); ok {
		r1 = rf(ctx, userUUID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ListUser provides a mock function with given fields: ctx, offset, limit
func (_m *UserService) ListUser(ctx context.Context, offset int, limit int) ([]model.UserInfo, error) {
	ret := _m.Called(ctx, offset, limit)

	var r0 []model.UserInfo
	if rf, ok := ret.Get(0).(func(context.Context, int, int) []model.UserInfo); ok {
		r0 = rf(ctx, offset, limit)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]model.UserInfo)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, int, int) error); ok {
		r1 = rf(ctx, offset, limit)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UpdateUser provides a mock function with given fields: ctx, userInfo, passwd
func (_m *UserService) UpdateUser(ctx context.Context, userInfo *model.UserInfo, passwd string) error {
	ret := _m.Called(ctx, userInfo, passwd)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *model.UserInfo, string) error); ok {
		r0 = rf(ctx, userInfo, passwd)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
