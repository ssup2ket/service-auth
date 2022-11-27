// Code generated by mockery v2.15.0. DO NOT EDIT.

package mocks

import (
	context "context"

	entity "github.com/ssup2ket/service-auth/internal/domain/entity"
	mock "github.com/stretchr/testify/mock"

	repo "github.com/ssup2ket/service-auth/internal/domain/repo"

	uuid "github.com/ssup2ket/service-auth/pkg/entity/uuid"
)

// OutboxRepo is an autogenerated mock type for the OutboxRepo type
type OutboxRepo struct {
	mock.Mock
}

// Create provides a mock function with given fields: ctx, userInfo
func (_m *OutboxRepo) Create(ctx context.Context, userInfo *entity.Outbox) error {
	ret := _m.Called(ctx, userInfo)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *entity.Outbox) error); ok {
		r0 = rf(ctx, userInfo)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Delete provides a mock function with given fields: ctx, userUUID
func (_m *OutboxRepo) Delete(ctx context.Context, userUUID uuid.EntityUUID) error {
	ret := _m.Called(ctx, userUUID)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, uuid.EntityUUID) error); ok {
		r0 = rf(ctx, userUUID)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// WithTx provides a mock function with given fields: tx
func (_m *OutboxRepo) WithTx(tx repo.DBTx) repo.OutboxRepo {
	ret := _m.Called(tx)

	var r0 repo.OutboxRepo
	if rf, ok := ret.Get(0).(func(repo.DBTx) repo.OutboxRepo); ok {
		r0 = rf(tx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(repo.OutboxRepo)
		}
	}

	return r0
}

type mockConstructorTestingTNewOutboxRepo interface {
	mock.TestingT
	Cleanup(func())
}

// NewOutboxRepo creates a new instance of OutboxRepo. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewOutboxRepo(t mockConstructorTestingTNewOutboxRepo) *OutboxRepo {
	mock := &OutboxRepo{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
