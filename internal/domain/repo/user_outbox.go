package repo

import (
	"context"

	"github.com/rs/zerolog/log"
	"gorm.io/gorm"

	"github.com/ssup2ket/ssup2ket-auth-service/internal/domain/model"
	modeluuid "github.com/ssup2ket/ssup2ket-auth-service/pkg/model/uuid"
)

// User outbox repo
type UserOutboxRepo interface {
	WithTx(tx *DBTx) UserOutboxRepo

	Create(ctx context.Context, userInfo *model.UserOutbox) error
	Delete(ctx context.Context, userUUID modeluuid.ModelUUID) error
}

type UserOutboxRepoImp struct {
	db *gorm.DB
}

func NewUserOutboxRepoImp(repoDB *gorm.DB) *UserOutboxRepoImp {
	return &UserOutboxRepoImp{
		db: repoDB,
	}
}

func (u *UserOutboxRepoImp) WithTx(tx *DBTx) UserOutboxRepo {
	transaction := tx.getTx()
	return NewUserOutboxRepoImp(transaction)
}

func (u *UserOutboxRepoImp) Create(ctx context.Context, userOutbox *model.UserOutbox) error {
	result := u.db.Create(userOutbox)
	if result.Error != nil {
		log.Ctx(ctx).Error().Err(result.Error).Msg("Failed to create user outbox")
		return ErrServerError
	}
	return nil
}

func (u *UserOutboxRepoImp) Delete(ctx context.Context, userUUID modeluuid.ModelUUID) error {
	result := u.db.Delete(&model.UserOutbox{}, "id = ?", userUUID)
	if result.Error != nil {
		log.Ctx(ctx).Error().Err(result.Error).Msg("Failed to delete user secret in primary DB")
		return ErrServerError
	}
	return nil
}
