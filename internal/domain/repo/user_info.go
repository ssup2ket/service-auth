package repo

import (
	"context"

	"github.com/rs/zerolog/log"
	"gorm.io/gorm"

	"github.com/ssup2ket/ssup2ket-auth-service/internal/domain/model"
)

// User info repo
type UserInfoRepo interface {
	WithTx(tx *DBTx) UserInfoRepo

	List(ctx context.Context, offset int, limit int) ([]model.UserInfo, error)
	Create(ctx context.Context, userInfo *model.UserInfo) error
	Get(ctx context.Context, uuid string) (*model.UserInfo, error)
	Update(ctx context.Context, userInfo *model.UserInfo) error
	Delete(ctx context.Context, userUUID string) error
}

type UserInfoRepoMysql struct {
	db *gorm.DB
}

func NewUserInfoRepoImp(repoDB *gorm.DB) *UserInfoRepoMysql {
	return &UserInfoRepoMysql{
		db: repoDB,
	}
}

func (u *UserInfoRepoMysql) WithTx(tx *DBTx) UserInfoRepo {
	transaction := tx.getTx()
	return NewUserInfoRepoImp(transaction)
}

func (u *UserInfoRepoMysql) List(ctx context.Context, offset int, limit int) ([]model.UserInfo, error) {
	userInfos := []model.UserInfo{}
	result := u.db.Offset(offset).Limit(limit).Find(&userInfos)
	if result.Error != nil {
		log.Ctx(ctx).Error().Err(result.Error).Msg("Failed to list user info from primary DB")
		return nil, getReturnErr(result.Error)
	}
	return userInfos, nil
}

func (u *UserInfoRepoMysql) Create(ctx context.Context, userInfo *model.UserInfo) error {
	result := u.db.Create(userInfo)
	if result.Error != nil {
		log.Ctx(ctx).Error().Err(result.Error).Msg("Failed to create user")
		return getReturnErr(result.Error)
	}
	return nil
}

func (u *UserInfoRepoMysql) Get(ctx context.Context, uuid string) (*model.UserInfo, error) {
	userInfo := model.UserInfo{}
	result := u.db.First(&userInfo, "uuid = ?", uuid)
	if result.Error != nil {
		log.Ctx(ctx).Error().Err(result.Error).Msg("Failed to get user info from primary DB")
		return nil, getReturnErr(result.Error)
	}
	return &userInfo, nil
}

func (u *UserInfoRepoMysql) Update(ctx context.Context, userInfo *model.UserInfo) error {
	result := u.db.Updates(userInfo)
	if result.Error != nil {
		log.Ctx(ctx).Error().Err(result.Error).Msg("Failed to update user info in primary DB")
		return getReturnErr(result.Error)
	}
	return nil
}

func (u *UserInfoRepoMysql) Delete(ctx context.Context, uuid string) error {
	result := u.db.Delete(&model.UserInfo{}, "uuid = ?", uuid)
	if result.Error != nil {
		log.Ctx(ctx).Error().Err(result.Error).Msg("Failed to delete user info in primary DB")
		return getReturnErr(result.Error)
	}
	return nil
}
