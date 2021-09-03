package service

import (
	"context"

	"github.com/rs/zerolog/log"

	"github.com/ssup2ket/ssup2ket-auth-service/internal/domain/repo"
	"github.com/ssup2ket/ssup2ket-auth-service/pkg/authtoken"
	"github.com/ssup2ket/ssup2ket-auth-service/pkg/password"
)

type TokenService interface {
	CreateToken(ctx context.Context, loginID, passwd string) (string, error)
}

type TokenServiceImp struct {
	userInfoRepoSecondary   repo.UserInfoRepo
	userSecretRepoSecondary repo.UserSecretRepo
}

func NewTokenServiceImp(userInfoSecondary repo.UserInfoRepo, userSecretSecondary repo.UserSecretRepo) *TokenServiceImp {
	return &TokenServiceImp{
		userInfoRepoSecondary:   userInfoSecondary,
		userSecretRepoSecondary: userSecretSecondary,
	}
}

func (t *TokenServiceImp) CreateToken(ctx context.Context, loginID, passwd string) (string, error) {
	// Get user info, user secret by loginID
	userInfo, err := t.userInfoRepoSecondary.GetByLoginID(ctx, loginID)
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Msg("Failed to get user info by login ID")
		return "", getReturnErr(err)
	}
	userSecret, err := t.userSecretRepoSecondary.Get(ctx, userInfo.ID)
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Msg("Failed to get user secret")
		return "", getReturnErr(err)
	}

	// Validate login ID, password
	if !password.ValidatePasswd(passwd, userSecret.PasswdHash, userSecret.PasswdSalt) {
		return "", ErrUnauthorized
	}

	// Create auth token
	token, err := authtoken.CreateAuthToken(&authtoken.AuthInfo{UserID: userInfo.ID.String(), UserLoginID: userInfo.LoginID})
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Msg("Failed to get user secret")
		return "", getReturnErr(err)
	}

	return token, nil
}
