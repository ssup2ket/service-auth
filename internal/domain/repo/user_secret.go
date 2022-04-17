package repo

import (
	"context"

	"github.com/rs/zerolog/log"
	"gorm.io/gorm"

	"github.com/ssup2ket/ssup2ket-auth-service/internal/domain/entity"
	"github.com/ssup2ket/ssup2ket-auth-service/pkg/entity/uuid"
)

// User secret repo
type UserSecretRepo interface {
	WithTx(tx DBTx) UserSecretRepo

	Create(ctx context.Context, userSecret *entity.UserSecret) error
	Get(ctx context.Context, userUUID uuid.EntityUUID) (*entity.UserSecret, error)
	Update(ctx context.Context, userSecret *entity.UserSecret) error
	Delete(ctx context.Context, userUUID uuid.EntityUUID) error
}

type UserSecretRepoImp struct {
	db *gorm.DB
}

func NewUserSecretRepoImp(repoDB *gorm.DB) *UserSecretRepoImp {
	return &UserSecretRepoImp{
		db: repoDB,
	}
}

func (u *UserSecretRepoImp) WithTx(tx DBTx) UserSecretRepo {
	transaction := tx.GetTx()
	return NewUserSecretRepoImp(transaction)
}

func (u *UserSecretRepoImp) Create(ctx context.Context, userSecret *entity.UserSecret) error {
	result := u.db.Create(userSecret)
	if result.Error != nil {
		log.Ctx(ctx).Error().Err(result.Error).Msg("Failed to create user secret")
		return ErrServerError
	}
	return nil
}

func (u *UserSecretRepoImp) Get(ctx context.Context, userUUID uuid.EntityUUID) (*entity.UserSecret, error) {
	user := entity.UserSecret{}
	result := u.db.First(&user, "id = ?", userUUID)
	if result.Error != nil {
		log.Ctx(ctx).Error().Err(result.Error).Msg("Failed to get user secret from DB")
		return nil, ErrServerError
	}
	return &user, nil
}

func (u *UserSecretRepoImp) Update(ctx context.Context, userSecret *entity.UserSecret) error {
	result := u.db.Updates(userSecret)
	if result.Error != nil {
		log.Ctx(ctx).Error().Err(result.Error).Msg("Failed to update user secret in DB")
		return ErrServerError
	}
	return nil
}

func (u *UserSecretRepoImp) Delete(ctx context.Context, userUUID uuid.EntityUUID) error {
	result := u.db.Delete(&entity.UserSecret{}, "id = ?", userUUID)
	if result.Error != nil {
		log.Ctx(ctx).Error().Err(result.Error).Msg("Failed to delete user secret in DB")
		return ErrServerError
	}
	return nil
}
