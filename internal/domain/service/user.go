package service

import (
	"context"

	"github.com/rs/zerolog/log"
	uuid "github.com/satori/go.uuid"

	"github.com/ssup2ket/ssup2ket-auth-service/internal/domain/model"
	"github.com/ssup2ket/ssup2ket-auth-service/internal/domain/repo"
	"github.com/ssup2ket/ssup2ket-auth-service/pkg/auth"
)

type UserService interface {
	ListUser(ctx context.Context, offset int, limit int) ([]model.UserInfo, error)
	CreateUser(ctx context.Context, userInfo *model.UserInfo, passwd string) (*model.UserInfo, error)
	GetUser(ctx context.Context, userUUID string) (*model.UserInfo, error)
	UpdateUser(ctx context.Context, userInfo *model.UserInfo, passwd string) error
	DeleteUser(ctx context.Context, userUUID string) error
}

type UserServiceImp struct {
	userInfoRepo   repo.UserInfoRepo
	userSecretRepo repo.UserSecretRepo
}

func NewUserServiceImp(userInfo repo.UserInfoRepo, userSecret repo.UserSecretRepo) *UserServiceImp {
	return &UserServiceImp{
		userInfoRepo:   userInfo,
		userSecretRepo: userSecret,
	}
}

func (u *UserServiceImp) ListUser(ctx context.Context, offset int, limit int) ([]model.UserInfo, error) {
	var err error

	// Set default limit
	if limit == 0 {
		limit = 50
	}

	// List users
	users, err := u.userInfoRepo.ListSecondary(ctx, offset, limit)
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Msg("Failed to list user from DB")
		return nil, getReturnErr(err)
	}
	return users, nil
}

func (u *UserServiceImp) CreateUser(ctx context.Context, userInfo *model.UserInfo, passwd string) (*model.UserInfo, error) {
	var err error

	// Begin transaction
	tx := repo.NewDBTx()
	_ = tx.Begin()
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
	uuid := uuid.NewV4().String()

	// Create user info
	userInfo.UUID = uuid
	if err = u.userInfoRepo.WithTx(tx).Create(ctx, userInfo); err != nil {
		log.Ctx(ctx).Error().Err(err).Msg("Failed to create user info to DB")
		return nil, getReturnErr(err)
	}

	// Create user secret
	hash, salt, err := auth.GetPasswordHashAndSalt([]byte(passwd))
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Msg("Failed to create password hash and salt")
		return nil, err
	}
	userSecret := model.UserSecret{
		UUID:       uuid,
		PasswdHash: hash,
		PasswdSalt: salt,
	}
	if err = u.userSecretRepo.WithTx(tx).Create(ctx, &userSecret); err != nil {
		log.Ctx(ctx).Error().Err(err).Msg("Failed to create user secret to DB")
		return nil, getReturnErr(err)
	}

	// Commit transaction
	if err = tx.Commit(); err != nil {
		log.Ctx(ctx).Error().Err(err).Msg("Commit transaction error for creating user")
		return nil, getReturnErr(err)
	}
	return userInfo, nil
}

func (u *UserServiceImp) GetUser(ctx context.Context, userUUID string) (*model.UserInfo, error) {
	var err error

	// Get user info
	userInfo, err := u.userInfoRepo.GetSecondary(ctx, userUUID)
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Msg("Failed to get user from DB")
		return nil, getReturnErr(err)
	}
	return userInfo, nil
}

func (u *UserServiceImp) UpdateUser(ctx context.Context, userInfo *model.UserInfo, passwd string) error {
	var err error

	// Begin transaction
	tx := repo.NewDBTx()
	tx.Begin()
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
	_, err = u.userInfoRepo.WithTx(tx).GetPrimary(ctx, userInfo.UUID)
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Msg("Failed to get user from DB")
		return getReturnErr(err)
	}

	// Update user info
	if err = u.userInfoRepo.WithTx(tx).UpdateUser(ctx, userInfo); err != nil {
		log.Ctx(ctx).Error().Err(err).Msg("Failed to update user from DB")
		return getReturnErr(err)
	}

	// Update user secret
	hash, salt, err := auth.GetPasswordHashAndSalt([]byte(passwd))
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Msg("Failed to create password hash and salt")
		return getReturnErr(err)
	}
	userSecret := model.UserSecret{
		UUID:       userInfo.UUID,
		PasswdHash: hash,
		PasswdSalt: salt,
	}
	if err = u.userSecretRepo.WithTx(tx).Update(ctx, &userSecret); err != nil {
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

func (u *UserServiceImp) DeleteUser(ctx context.Context, userUUID string) error {
	var err error

	// Begin transaction
	tx := repo.NewDBTx()
	tx.Begin()
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
	_, err = u.userInfoRepo.WithTx(tx).GetPrimary(ctx, userUUID)
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Msg("Failed to get user info from DB")
		return getReturnErr(err)
	}

	// Delete user info
	if err := u.userInfoRepo.WithTx(tx).Delete(ctx, userUUID); err != nil {
		log.Ctx(ctx).Error().Err(err).Msg("Failed to delete user info from DB")
		return getReturnErr(err)
	}

	// Delete user secret
	if err := u.userSecretRepo.WithTx(tx).Delete(ctx, userUUID); err != nil {
		log.Ctx(ctx).Error().Err(err).Msg("Failed to delete user secret from DB")
		return getReturnErr(err)
	}

	// Commit transaction
	if err = tx.Commit(); err != nil {
		log.Ctx(ctx).Error().Err(err).Msg("Commit transaction error for deleting user")
		return getReturnErr(err)
	}
	return nil
}
