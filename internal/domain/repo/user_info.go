package repo

import (
	"context"

	"github.com/rs/zerolog/log"
	"gorm.io/gorm"

	"github.com/ssup2ket/ssup2ket-auth-service/internal/domain/entity"
	"github.com/ssup2ket/ssup2ket-auth-service/pkg/entity/uuid"
)

// User info repo
type UserInfoRepo interface {
	WithTx(tx *DBTx) UserInfoRepo

	List(ctx context.Context, offset int, limit int) ([]entity.UserInfo, error)
	Create(ctx context.Context, userInfo *entity.UserInfo) error
	Get(ctx context.Context, userUUID uuid.EntityUUID) (*entity.UserInfo, error)
	GetByLoginID(ctx context.Context, userLoginID string) (*entity.UserInfo, error)
	Update(ctx context.Context, userInfo *entity.UserInfo) error
	Delete(ctx context.Context, userUUID uuid.EntityUUID) error
}

type UserInfoRepoImp struct {
	db *gorm.DB
}

func NewUserInfoRepoImp(repoDB *gorm.DB) *UserInfoRepoImp {
	return &UserInfoRepoImp{
		db: repoDB,
	}
}

func (u *UserInfoRepoImp) WithTx(tx *DBTx) UserInfoRepo {
	transaction := tx.getTx()
	return NewUserInfoRepoImp(transaction)
}

func (u *UserInfoRepoImp) List(ctx context.Context, offset int, limit int) ([]entity.UserInfo, error) {
	userInfos := []entity.UserInfo{}
	result := u.db.Offset(offset).Limit(limit).Find(&userInfos)
	if result.Error != nil {
		log.Ctx(ctx).Error().Err(result.Error).Msg("Failed to list user info from DB")
		return nil, getReturnErr(result.Error)
	}
	return userInfos, nil
}

func (u *UserInfoRepoImp) Create(ctx context.Context, userInfo *entity.UserInfo) error {
	result := u.db.Create(userInfo)
	if result.Error != nil {
		log.Ctx(ctx).Error().Err(result.Error).Msg("Failed to create user in DB")
		return getReturnErr(result.Error)
	}
	return nil
}

func (u *UserInfoRepoImp) Get(ctx context.Context, userUUID uuid.EntityUUID) (*entity.UserInfo, error) {
	userInfo := entity.UserInfo{}
	result := u.db.First(&userInfo, "id = ?", userUUID)
	if result.Error != nil {
		log.Ctx(ctx).Error().Err(result.Error).Msg("Failed to get user info from DB")
		return nil, getReturnErr(result.Error)
	}
	return &userInfo, nil
}

func (u *UserInfoRepoImp) GetByLoginID(ctx context.Context, userLoginID string) (*entity.UserInfo, error) {
	userInfo := entity.UserInfo{}
	result := u.db.First(&userInfo, "login_id = ?", userLoginID)
	if result.Error != nil {
		log.Ctx(ctx).Error().Err(result.Error).Msg("Failed to get user info from DB by user login ID")
		return nil, getReturnErr(result.Error)
	}
	return &userInfo, nil
}

func (u *UserInfoRepoImp) Update(ctx context.Context, userInfo *entity.UserInfo) error {
	result := u.db.Updates(userInfo)
	if result.Error != nil {
		log.Ctx(ctx).Error().Err(result.Error).Msg("Failed to update user info in DB")
		return getReturnErr(result.Error)
	}
	return nil
}

func (u *UserInfoRepoImp) Delete(ctx context.Context, userUUID uuid.EntityUUID) error {
	result := u.db.Delete(&entity.UserInfo{}, "id = ?", userUUID)
	if result.Error != nil {
		log.Ctx(ctx).Error().Err(result.Error).Msg("Failed to delete user info in DB")
		return getReturnErr(result.Error)
	}
	return nil
}
