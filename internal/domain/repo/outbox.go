package repo

import (
	"context"

	"github.com/rs/zerolog/log"
	"gorm.io/gorm"

	"github.com/ssup2ket/ssup2ket-auth-service/internal/domain/entity"
	"github.com/ssup2ket/ssup2ket-auth-service/pkg/entity/uuid"
)

// Outbox repo
type OutboxRepo interface {
	WithTx(tx *DBTx) OutboxRepo

	Create(ctx context.Context, userInfo *entity.Outbox) error
	Delete(ctx context.Context, userUUID uuid.EntityUUID) error
}

type OutboxRepoImp struct {
	db *gorm.DB
}

func NewOutboxRepoImp(repoDB *gorm.DB) *OutboxRepoImp {
	return &OutboxRepoImp{
		db: repoDB,
	}
}

func (u *OutboxRepoImp) WithTx(tx *DBTx) OutboxRepo {
	transaction := tx.getTx()
	return NewOutboxRepoImp(transaction)
}

func (u *OutboxRepoImp) Create(ctx context.Context, outbox *entity.Outbox) error {
	result := u.db.Create(outbox)
	if result.Error != nil {
		log.Ctx(ctx).Error().Err(result.Error).Msg("Failed to create outbox")
		return ErrServerError
	}
	return nil
}

func (u *OutboxRepoImp) Delete(ctx context.Context, userUUID uuid.EntityUUID) error {
	result := u.db.Delete(&entity.Outbox{}, "id = ?", userUUID)
	if result.Error != nil {
		log.Ctx(ctx).Error().Err(result.Error).Msg("Failed to delete outbox in primary DB")
		return ErrServerError
	}
	return nil
}
