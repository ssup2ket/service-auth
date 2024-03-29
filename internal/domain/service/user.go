package service

import (
	"context"
	"encoding/json"

	opentracing "github.com/opentracing/opentracing-go"
	"github.com/rs/zerolog/log"

	"github.com/ssup2ket/service-auth/internal/domain/entity"
	"github.com/ssup2ket/service-auth/internal/domain/repo"
	"github.com/ssup2ket/service-auth/pkg/auth/hashing"
	"github.com/ssup2ket/service-auth/pkg/entity/uuid"
	"github.com/ssup2ket/service-auth/pkg/tracing"
)

const (
	AggregateTypeUser    = "User"
	EventTypeUserCreated = "UserCreated"
	EventTypeUserDeleted = "UserDeleted"
)

type userOutboxPayload struct {
	ID      string `json:"id"`
	LoginID string `json:"loginId"`
	Role    string `json:"role"`
}

// User Service
type UserService interface {
	ListUser(ctx context.Context, offset int, limit int) ([]entity.UserInfo, error)
	CreateUser(ctx context.Context, userInfo *entity.UserInfo, passwd string) (*entity.UserInfo, error)
	GetUser(ctx context.Context, userUUID uuid.EntityUUID) (*entity.UserInfo, error)
	UpdateUser(ctx context.Context, userInfo *entity.UserInfo, passwd string) error
	DeleteUser(ctx context.Context, userUUID uuid.EntityUUID) error
}

type UserServiceImp struct {
	repoDBTx repo.DBTx

	outBoxRepoPrimary       repo.OutboxRepo
	userInfoRepoPrimary     repo.UserInfoRepo
	userInfoRepoSecondary   repo.UserInfoRepo
	userSecretRepoPrimary   repo.UserSecretRepo
	userSecretRepoSecondary repo.UserSecretRepo
}

func NewUserServiceImp(dbTx repo.DBTx, userOutBoxPrimary repo.OutboxRepo, userInfoPrimary, userInfoSecondary repo.UserInfoRepo,
	userSecretPrimary, userSecretSecondary repo.UserSecretRepo) *UserServiceImp {
	return &UserServiceImp{
		repoDBTx: dbTx,

		outBoxRepoPrimary:       userOutBoxPrimary,
		userInfoRepoPrimary:     userInfoPrimary,
		userInfoRepoSecondary:   userInfoSecondary,
		userSecretRepoPrimary:   userSecretPrimary,
		userSecretRepoSecondary: userSecretSecondary,
	}
}

func (u *UserServiceImp) ListUser(ctx context.Context, offset int, limit int) ([]entity.UserInfo, error) {
	var err error

	// Set default limit
	if limit == 0 {
		limit = 50
	}

	// List users
	users, err := u.userInfoRepoSecondary.List(ctx, offset, limit)
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Msg("Failed to list user from DB")
		return nil, getReturnErr(err)
	}
	return users, nil
}

func (u *UserServiceImp) CreateUser(ctx context.Context, userInfo *entity.UserInfo, passwd string) (*entity.UserInfo, error) {
	var err error

	// Begin transaction
	tx, _ := u.repoDBTx.Begin()
	defer func() {
		if err != nil {
			if err = tx.Rollback(); err != nil {
				log.Ctx(ctx).Error().Err(err).Msg("Rollback transaction error for creating user")
				return
			}
			log.Ctx(ctx).Error().Err(err).Msg("Create user request is canceled")
			return
		}
	}()

	// Generate UUID to share to userInfo and userSecret
	userUUID := uuid.NewV4()

	// Create user info
	userInfo.ID = userUUID
	if err = u.userInfoRepoPrimary.WithTx(tx).Create(ctx, userInfo); err != nil {
		log.Ctx(ctx).Error().Err(err).Msg("Failed to create user info to DB")
		return nil, getReturnErr(err)
	}

	// Create user secret
	hash, salt, err := hashing.GetStrHashAndSalt(passwd)
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Msg("Failed to create password hash and salt")
		return nil, err
	}
	userSecret := entity.UserSecret{
		ID:         userUUID,
		PasswdHash: hash,
		PasswdSalt: salt,
	}
	if err = u.userSecretRepoPrimary.WithTx(tx).Create(ctx, &userSecret); err != nil {
		log.Ctx(ctx).Error().Err(err).Msg("Failed to create user secret to DB")
		return nil, getReturnErr(err)
	}

	// Get user outbox payload
	userOutboxPayload := userOutboxPayload{
		ID:      userInfo.ID.String(),
		LoginID: userInfo.LoginID,
		Role:    string(userInfo.Role),
	}
	userOutboxPayloadJSON, err := json.Marshal(userOutboxPayload)
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Msg("Failed to marshal user outbox payload")
		return nil, getReturnErr(err)
	}

	// Get span context as JSON
	tracer := opentracing.GlobalTracer()
	span, ctx := opentracing.StartSpanFromContext(ctx, "CreateUser")
	spanContext, err := tracing.GetSpanContextAsJSON(tracer, span)
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Msg("Failed to get user outbox spancontext")
		return nil, getReturnErr(err)
	}

	// Insert created user info to outbox table to public a user create event
	userOutbox := entity.Outbox{
		ID:            uuid.NewV4(),
		AggregateType: AggregateTypeUser,
		AggregateID:   userInfo.ID.String(),
		EventType:     EventTypeUserCreated,
		Payload:       string(userOutboxPayloadJSON),
		SpanContext:   spanContext,
	}
	if err = u.outBoxRepoPrimary.WithTx(tx).Create(ctx, &userOutbox); err != nil {
		log.Ctx(ctx).Error().Err(err).Msg("Failed to insert created user to outbox table")
		return nil, getReturnErr(err)
	}

	// Commit transaction
	if err = tx.Commit(); err != nil {
		log.Ctx(ctx).Error().Err(err).Msg("Commit transaction error for creating user")
		return nil, getReturnErr(err)
	}
	return userInfo, nil
}

func (u *UserServiceImp) GetUser(ctx context.Context, userUUID uuid.EntityUUID) (*entity.UserInfo, error) {
	var err error

	// Get user info
	userInfo, err := u.userInfoRepoSecondary.Get(ctx, userUUID)
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Msg("Failed to get user from DB")
		return nil, getReturnErr(err)
	}
	return userInfo, nil
}

func (u *UserServiceImp) UpdateUser(ctx context.Context, userInfo *entity.UserInfo, passwd string) error {
	var err error

	// Begin transaction
	tx, _ := u.repoDBTx.Begin()
	defer func() {
		if err != nil {
			if err = tx.Rollback(); err != nil {
				log.Ctx(ctx).Error().Err(err).Msg("Rollback transaction error for updating user")
				return
			}
			log.Ctx(ctx).Error().Err(err).Msg("Update user request is canceled")
			return
		}
	}()

	// Get user info
	_, err = u.userInfoRepoPrimary.WithTx(tx).Get(ctx, userInfo.ID)
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Msg("Failed to get user from DB")
		return getReturnErr(err)
	}

	// Update user info
	if err = u.userInfoRepoPrimary.WithTx(tx).Update(ctx, userInfo); err != nil {
		log.Ctx(ctx).Error().Err(err).Msg("Failed to update user from DB")
		return getReturnErr(err)
	}

	// Update user secret
	hash, salt, err := hashing.GetStrHashAndSalt(passwd)
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Msg("Failed to create password hash and salt")
		return getReturnErr(err)
	}
	userSecret := entity.UserSecret{
		ID:         userInfo.ID,
		PasswdHash: hash,
		PasswdSalt: salt,
	}
	if err = u.userSecretRepoPrimary.WithTx(tx).Update(ctx, &userSecret); err != nil {
		log.Ctx(ctx).Error().Err(err).Msg("Failed to create user secret to DB")
		return getReturnErr(err)
	}

	// Commit transaction
	if err = tx.Commit(); err != nil {
		log.Ctx(ctx).Error().Err(err).Msg("Commit transaction error for updating user")
		return getReturnErr(err)
	}
	return nil
}

func (u *UserServiceImp) DeleteUser(ctx context.Context, userUUID uuid.EntityUUID) error {
	var err error

	// Begin transaction
	tx, _ := u.repoDBTx.Begin()
	defer func() {
		if err != nil {
			if err = tx.Rollback(); err != nil {
				log.Ctx(ctx).Error().Err(err).Msg("Rollback transaction error for deleting user")
				return
			}
			log.Ctx(ctx).Error().Err(err).Msg("Delete user request is canceled")
			return
		}
	}()

	// Get user info
	userInfo, err := u.userInfoRepoPrimary.WithTx(tx).Get(ctx, userUUID)
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Msg("Failed to get user info from DB")
		return getReturnErr(err)
	}

	// Delete user info
	if err := u.userInfoRepoPrimary.WithTx(tx).Delete(ctx, userUUID); err != nil {
		log.Ctx(ctx).Error().Err(err).Msg("Failed to delete user info from DB")
		return getReturnErr(err)
	}

	// Delete user secret
	if err := u.userSecretRepoPrimary.WithTx(tx).Delete(ctx, userUUID); err != nil {
		log.Ctx(ctx).Error().Err(err).Msg("Failed to delete user secret from DB")
		return getReturnErr(err)
	}

	// Get user outbox payload
	userOutboxPayload := userOutboxPayload{
		ID:      userInfo.ID.String(),
		LoginID: userInfo.LoginID,
		Role:    string(userInfo.Role),
	}
	userOutboxPayloadJSON, err := json.Marshal(userOutboxPayload)
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Msg("Failed to marshal user outbox payload")
		return getReturnErr(err)
	}

	// Get span context as JSON
	tracer := opentracing.GlobalTracer()
	span, ctx := opentracing.StartSpanFromContext(ctx, "DeleteUser")
	spanContext, err := tracing.GetSpanContextAsJSON(tracer, span)
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Msg("Failed to get user outbox spancontext")
		return getReturnErr(err)
	}

	// Insert deleted user info to outbox table to public a user delete event
	userOutbox := entity.Outbox{
		ID:            uuid.NewV4(),
		AggregateType: AggregateTypeUser,
		AggregateID:   userUUID.String(),
		EventType:     EventTypeUserDeleted,
		Payload:       string(userOutboxPayloadJSON),
		SpanContext:   spanContext,
	}
	if err = u.outBoxRepoPrimary.WithTx(tx).Create(ctx, &userOutbox); err != nil {
		log.Ctx(ctx).Error().Err(err).Msg("Failed to insert created user to outbox table")
		return getReturnErr(err)
	}

	// Commit transaction
	if err = tx.Commit(); err != nil {
		log.Ctx(ctx).Error().Err(err).Msg("Commit transaction error for deleting user")
		return getReturnErr(err)
	}
	return nil
}
