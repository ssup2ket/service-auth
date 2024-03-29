// Code generated by mockery v2.15.0. DO NOT EDIT.

package mocks

import (
	context "context"

	entity "github.com/ssup2ket/service-auth/internal/domain/entity"
	mock "github.com/stretchr/testify/mock"

	repo "github.com/ssup2ket/service-auth/internal/domain/repo"

	uuid "github.com/ssup2ket/service-auth/pkg/entity/uuid"
)

// UserSecretRepo is an autogenerated mock type for the UserSecretRepo type
type UserSecretRepo struct {
	mock.Mock
}

// Create provides a mock function with given fields: ctx, userSecret
func (_m *UserSecretRepo) Create(ctx context.Context, userSecret *entity.UserSecret) error {
	ret := _m.Called(ctx, userSecret)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *entity.UserSecret) error); ok {
		r0 = rf(ctx, userSecret)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Delete provides a mock function with given fields: ctx, userUUID
func (_m *UserSecretRepo) Delete(ctx context.Context, userUUID uuid.EntityUUID) error {
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
func (_m *UserSecretRepo) Get(ctx context.Context, userUUID uuid.EntityUUID) (*entity.UserSecret, error) {
	ret := _m.Called(ctx, userUUID)

	var r0 *entity.UserSecret
	if rf, ok := ret.Get(0).(func(context.Context, uuid.EntityUUID) *entity.UserSecret); ok {
		r0 = rf(ctx, userUUID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*entity.UserSecret)
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

// Update provides a mock function with given fields: ctx, userSecret
func (_m *UserSecretRepo) Update(ctx context.Context, userSecret *entity.UserSecret) error {
	ret := _m.Called(ctx, userSecret)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *entity.UserSecret) error); ok {
		r0 = rf(ctx, userSecret)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// WithTx provides a mock function with given fields: tx
func (_m *UserSecretRepo) WithTx(tx repo.DBTx) repo.UserSecretRepo {
	ret := _m.Called(tx)

	var r0 repo.UserSecretRepo
	if rf, ok := ret.Get(0).(func(repo.DBTx) repo.UserSecretRepo); ok {
		r0 = rf(tx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(repo.UserSecretRepo)
		}
	}

	return r0
}

type mockConstructorTestingTNewUserSecretRepo interface {
	mock.TestingT
	Cleanup(func())
}

// NewUserSecretRepo creates a new instance of UserSecretRepo. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewUserSecretRepo(t mockConstructorTestingTNewUserSecretRepo) *UserSecretRepo {
	mock := &UserSecretRepo{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
