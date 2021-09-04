// Code generated by mockery v0.0.0-dev. DO NOT EDIT.

package mocks

import (
	context "context"

	model "github.com/ssup2ket/ssup2ket-auth-service/internal/domain/model"
	mock "github.com/stretchr/testify/mock"

	repo "github.com/ssup2ket/ssup2ket-auth-service/internal/domain/repo"

	uuidmodel "github.com/ssup2ket/ssup2ket-auth-service/pkg/uuidmodel"
)

// UserSecretRepo is an autogenerated mock type for the UserSecretRepo type
type UserSecretRepo struct {
	mock.Mock
}

// Create provides a mock function with given fields: ctx, userSecret
func (_m *UserSecretRepo) Create(ctx context.Context, userSecret *model.UserSecret) error {
	ret := _m.Called(ctx, userSecret)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *model.UserSecret) error); ok {
		r0 = rf(ctx, userSecret)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Delete provides a mock function with given fields: ctx, userUUID
func (_m *UserSecretRepo) Delete(ctx context.Context, userUUID uuidmodel.UUIDModel) error {
	ret := _m.Called(ctx, userUUID)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, uuidmodel.UUIDModel) error); ok {
		r0 = rf(ctx, userUUID)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Get provides a mock function with given fields: ctx, userUUID
func (_m *UserSecretRepo) Get(ctx context.Context, userUUID uuidmodel.UUIDModel) (*model.UserSecret, error) {
	ret := _m.Called(ctx, userUUID)

	var r0 *model.UserSecret
	if rf, ok := ret.Get(0).(func(context.Context, uuidmodel.UUIDModel) *model.UserSecret); ok {
		r0 = rf(ctx, userUUID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.UserSecret)
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

// Update provides a mock function with given fields: ctx, userSecret
func (_m *UserSecretRepo) Update(ctx context.Context, userSecret *model.UserSecret) error {
	ret := _m.Called(ctx, userSecret)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *model.UserSecret) error); ok {
		r0 = rf(ctx, userSecret)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// WithTx provides a mock function with given fields: tx
func (_m *UserSecretRepo) WithTx(tx *repo.DBTx) repo.UserSecretRepo {
	ret := _m.Called(tx)

	var r0 repo.UserSecretRepo
	if rf, ok := ret.Get(0).(func(*repo.DBTx) repo.UserSecretRepo); ok {
		r0 = rf(tx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(repo.UserSecretRepo)
		}
	}

	return r0
}
