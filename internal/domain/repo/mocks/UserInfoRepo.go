// Code generated by mockery v2.15.0. DO NOT EDIT.

package mocks

import (
	context "context"

	entity "github.com/ssup2ket/service-auth/internal/domain/entity"
	mock "github.com/stretchr/testify/mock"

	repo "github.com/ssup2ket/service-auth/internal/domain/repo"

	uuid "github.com/ssup2ket/service-auth/pkg/entity/uuid"
)

// UserInfoRepo is an autogenerated mock type for the UserInfoRepo type
type UserInfoRepo struct {
	mock.Mock
}

// Create provides a mock function with given fields: ctx, userInfo
func (_m *UserInfoRepo) Create(ctx context.Context, userInfo *entity.UserInfo) error {
	ret := _m.Called(ctx, userInfo)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *entity.UserInfo) error); ok {
		r0 = rf(ctx, userInfo)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Delete provides a mock function with given fields: ctx, userUUID
func (_m *UserInfoRepo) Delete(ctx context.Context, userUUID uuid.EntityUUID) error {
	ret := _m.Called(ctx, userUUID)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, uuid.EntityUUID) error); ok {
		r0 = rf(ctx, userUUID)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Get provides a mock function with given fields: ctx, userUUID
func (_m *UserInfoRepo) Get(ctx context.Context, userUUID uuid.EntityUUID) (*entity.UserInfo, error) {
	ret := _m.Called(ctx, userUUID)

	var r0 *entity.UserInfo
	if rf, ok := ret.Get(0).(func(context.Context, uuid.EntityUUID) *entity.UserInfo); ok {
		r0 = rf(ctx, userUUID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*entity.UserInfo)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, uuid.EntityUUID) error); ok {
		r1 = rf(ctx, userUUID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetByLoginID provides a mock function with given fields: ctx, userLoginID
func (_m *UserInfoRepo) GetByLoginID(ctx context.Context, userLoginID string) (*entity.UserInfo, error) {
	ret := _m.Called(ctx, userLoginID)

	var r0 *entity.UserInfo
	if rf, ok := ret.Get(0).(func(context.Context, string) *entity.UserInfo); ok {
		r0 = rf(ctx, userLoginID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*entity.UserInfo)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, userLoginID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// List provides a mock function with given fields: ctx, offset, limit
func (_m *UserInfoRepo) List(ctx context.Context, offset int, limit int) ([]entity.UserInfo, error) {
	ret := _m.Called(ctx, offset, limit)

	var r0 []entity.UserInfo
	if rf, ok := ret.Get(0).(func(context.Context, int, int) []entity.UserInfo); ok {
		r0 = rf(ctx, offset, limit)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]entity.UserInfo)
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

// Update provides a mock function with given fields: ctx, userInfo
func (_m *UserInfoRepo) Update(ctx context.Context, userInfo *entity.UserInfo) error {
	ret := _m.Called(ctx, userInfo)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *entity.UserInfo) error); ok {
		r0 = rf(ctx, userInfo)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// WithTx provides a mock function with given fields: tx
func (_m *UserInfoRepo) WithTx(tx repo.DBTx) repo.UserInfoRepo {
	ret := _m.Called(tx)

	var r0 repo.UserInfoRepo
	if rf, ok := ret.Get(0).(func(repo.DBTx) repo.UserInfoRepo); ok {
		r0 = rf(tx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(repo.UserInfoRepo)
		}
	}

	return r0
}

type mockConstructorTestingTNewUserInfoRepo interface {
	mock.TestingT
	Cleanup(func())
}

// NewUserInfoRepo creates a new instance of UserInfoRepo. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewUserInfoRepo(t mockConstructorTestingTNewUserInfoRepo) *UserInfoRepo {
	mock := &UserInfoRepo{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
