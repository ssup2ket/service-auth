package repo

import (
	"context"

	"github.com/rs/zerolog/log"
	"gorm.io/gorm"

	"github.com/ssup2ket/ssup2ket-auth-service/internal/domain/model"
)

// User secret repo
type UserSecretRepo interface {
	WithTx(tx *DBTx) UserSecretRepo

	Create(ctx context.Context, userSecret *model.UserSecret) error
	GetPrimary(ctx context.Context, userUUID string) (*model.UserSecret, error)
	GetSecondary(ctx context.Context, userUUID string) (*model.UserSecret, error)
	Update(ctx context.Context, userSecret *model.UserSecret) error
	Delete(ctx context.Context, userUUID string) error
}

type UserSecretRepoMysql struct {
	primaryDB   *gorm.DB
	secondaryDB *gorm.DB
}

func NewUserSecretRepoMysql() *UserSecretRepoMysql {
	return &UserSecretRepoMysql{
		primaryDB:   primaryMySQL,
		secondaryDB: secondaryMySQL,
	}
}

func (u *UserSecretRepoMysql) WithTx(tx *DBTx) UserSecretRepo {
	gormDB := tx.getTx()

	return &UserSecretRepoMysql{
		primaryDB:   gormDB,
		secondaryDB: u.secondaryDB,
	}
}

func (u *UserSecretRepoMysql) Create(ctx context.Context, userSecret *model.UserSecret) error {
	result := u.primaryDB.Create(userSecret)
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

func (u *UserSecretRepoMysql) GetPrimary(ctx context.Context, userUUID string) (*model.UserSecret, error) {
	user := model.UserSecret{}
	result := u.primaryDB.First(&user, "uuid = ?", userUUID)
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

func (u *UserSecretRepoMysql) GetSecondary(ctx context.Context, userUUID string) (*model.UserSecret, error) {
	user := model.UserSecret{}
	result := u.secondaryDB.First(&user, "uuid = ?", userUUID)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			log.Ctx(ctx).Error().Err(result.Error).Msg("User secret does not exist in second DB")
			return nil, ErrNotFound
		}
		log.Ctx(ctx).Error().Err(result.Error).Msg("Failed to get user secret from second DB")
		return nil, ErrServerError
	}
	return &user, nil
}

func (u *UserSecretRepoMysql) Update(ctx context.Context, userSecret *model.UserSecret) error {
	result := u.primaryDB.Updates(userSecret)
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

func (u *UserSecretRepoMysql) Delete(ctx context.Context, userUUID string) error {
	result := u.primaryDB.Delete(&model.UserSecret{}, "uuid = ?", userUUID)
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
