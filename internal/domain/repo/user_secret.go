package repo

import (
	"context"

	"github.com/rs/zerolog/log"
	"gorm.io/gorm"

	"github.com/ssup2ket/ssup2ket-auth-service/internal/domain/model"
	modeluuid "github.com/ssup2ket/ssup2ket-auth-service/pkg/model/uuid"
)

// User secret repo
type UserSecretRepo interface {
	WithTx(tx *DBTx) UserSecretRepo

	Create(ctx context.Context, userSecret *model.UserSecret) error
	Get(ctx context.Context, userUUID modeluuid.ModelUUID) (*model.UserSecret, error)
	Update(ctx context.Context, userSecret *model.UserSecret) error
	Delete(ctx context.Context, userUUID modeluuid.ModelUUID) error
}

type UserSecretRepoImp struct {
	db *gorm.DB
}

func NewUserSecretRepoImp(repoDB *gorm.DB) *UserSecretRepoImp {
	return &UserSecretRepoImp{
		db: repoDB,
	}
}

func (u *UserSecretRepoImp) WithTx(tx *DBTx) UserSecretRepo {
	transaction := tx.getTx()
	return NewUserSecretRepoImp(transaction)
}

func (u *UserSecretRepoImp) Create(ctx context.Context, userSecret *model.UserSecret) error {
	result := u.db.Create(userSecret)
	if result.Error != nil {
		if result.Error == gorm.ErrInvalidData {
			log.Ctx(ctx).Error().Err(result.Error).Msg("Failed to create user secret because of duplication")
			return ErrConflict
		}
		log.Ctx(ctx).Error().Err(result.Error).Msg("Failed to create secret")
		return ErrServerError
	}
	return nil
}

func (u *UserSecretRepoImp) Get(ctx context.Context, userUUID modeluuid.ModelUUID) (*model.UserSecret, error) {
	user := model.UserSecret{}
	result := u.db.First(&user, "id = ?", userUUID)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			log.Ctx(ctx).Error().Err(result.Error).Msg("User secret does not exist in primary DB")
			return nil, ErrNotFound
		}
		log.Ctx(ctx).Error().Err(result.Error).Msg("Failed to get user secret from primary DB")
		return nil, ErrServerError
	}
	return &user, nil
}

func (u *UserSecretRepoImp) Update(ctx context.Context, userSecret *model.UserSecret) error {
	result := u.db.Updates(userSecret)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			log.Ctx(ctx).Error().Err(result.Error).Msg("User secret does not exist in primary DB")
			return ErrNotFound
		}
		log.Ctx(ctx).Error().Err(result.Error).Msg("Failed to update user secret in primary DB")
		return ErrServerError
	}
	return nil
}

func (u *UserSecretRepoImp) Delete(ctx context.Context, userUUID modeluuid.ModelUUID) error {
	result := u.db.Delete(&model.UserSecret{}, "id = ?", userUUID)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			log.Ctx(ctx).Error().Err(result.Error).Msg("User secret does not exist in primary DB")
			return ErrNotFound
		}
		log.Ctx(ctx).Error().Err(result.Error).Msg("Failed to delete user secret in primary DB")
		return ErrServerError
	}
	return nil
}
