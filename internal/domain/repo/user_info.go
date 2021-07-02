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

	ListPrimary(ctx context.Context, offset int, limit int) ([]model.UserInfo, error)
	ListSecondary(ctx context.Context, offset int, limit int) ([]model.UserInfo, error)
	Create(ctx context.Context, userInfo *model.UserInfo) error
	GetPrimary(ctx context.Context, uuid string) (*model.UserInfo, error)
	GetSecondary(ctx context.Context, uuid string) (*model.UserInfo, error)
	UpdateUser(ctx context.Context, userInfo *model.UserInfo) error
	Delete(ctx context.Context, userUUID string) error
}

type UserInfoRepoMysql struct {
	primaryDB   *gorm.DB
	secondaryDB *gorm.DB
}

func NewUserInfoRepoMysql() *UserInfoRepoMysql {
	return &UserInfoRepoMysql{
		primaryDB:   primaryMySQL,
		secondaryDB: secondaryMySQL,
	}
}

func (u *UserInfoRepoMysql) WithTx(tx *DBTx) UserInfoRepo {
	gormDB := tx.getTx()
	return &UserInfoRepoMysql{
		primaryDB:   gormDB,
		secondaryDB: u.secondaryDB,
	}
}

func (u *UserInfoRepoMysql) ListPrimary(ctx context.Context, offset int, limit int) ([]model.UserInfo, error) {
	userInfos := []model.UserInfo{}
	result := u.primaryDB.Offset(offset).Limit(limit).Find(&userInfos)
	if result.Error != nil {
		log.Ctx(ctx).Error().Err(result.Error).Msg("Failed to list user info from primary DB")
		return nil, getReturnErr(result.Error)
	}
	return userInfos, nil
}

func (u *UserInfoRepoMysql) ListSecondary(ctx context.Context, offset int, limit int) ([]model.UserInfo, error) {
	userInfos := []model.UserInfo{}
	result := u.secondaryDB.Offset(offset).Limit(limit).Find(&userInfos)
	if result.Error != nil {
		log.Ctx(ctx).Error().Err(result.Error).Msg("Failed to list user info from secondary DB")
		return nil, getReturnErr(result.Error)
	}
	return userInfos, nil
}

func (u *UserInfoRepoMysql) Create(ctx context.Context, userInfo *model.UserInfo) error {
	result := u.primaryDB.Create(userInfo)
	if result.Error != nil {
		log.Ctx(ctx).Error().Err(result.Error).Msg("Failed to create user")
		return getReturnErr(result.Error)
	}
	return nil
}

func (u *UserInfoRepoMysql) GetPrimary(ctx context.Context, uuid string) (*model.UserInfo, error) {
	userInfo := model.UserInfo{}
	result := u.primaryDB.First(&userInfo, "uuid = ?", uuid)
	if result.Error != nil {
		log.Ctx(ctx).Error().Err(result.Error).Msg("Failed to get user info from primary DB")
		return nil, getReturnErr(result.Error)
	}
	return &userInfo, nil
}

func (u *UserInfoRepoMysql) GetSecondary(ctx context.Context, uuid string) (*model.UserInfo, error) {
	userInfo := model.UserInfo{}
	result := u.secondaryDB.First(&userInfo, "uuid = ?", uuid)
	if result.Error != nil {
		log.Ctx(ctx).Error().Err(result.Error).Msg("Failed to get user info from secondary DB")
		return nil, getReturnErr(result.Error)
	}
	return &userInfo, nil
}

func (u *UserInfoRepoMysql) UpdateUser(ctx context.Context, userInfo *model.UserInfo) error {
	result := u.primaryDB.Updates(userInfo)
	if result.Error != nil {
		log.Ctx(ctx).Error().Err(result.Error).Msg("Failed to update user info in primary DB")
		return getReturnErr(result.Error)
	}
	return nil
}

func (u *UserInfoRepoMysql) Delete(ctx context.Context, uuid string) error {
	result := u.primaryDB.Delete(&model.UserInfo{}, "uuid = ?", uuid)
	if result.Error != nil {
		log.Ctx(ctx).Error().Err(result.Error).Msg("Failed to delete user info in primary DB")
		return getReturnErr(result.Error)
	}
	return nil
}
