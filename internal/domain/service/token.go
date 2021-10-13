package service

import (
	"context"

	"github.com/rs/zerolog/log"

	"github.com/ssup2ket/ssup2ket-auth-service/internal/domain/repo"
	"github.com/ssup2ket/ssup2ket-auth-service/pkg/auth/password"
	"github.com/ssup2ket/ssup2ket-auth-service/pkg/auth/token"
)

// Token service
type TokenService interface {
	CreateToken(ctx context.Context, loginID, passwd string) (*token.TokenInfo, error)
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

func (t *TokenServiceImp) CreateToken(ctx context.Context, loginID, passwd string) (*token.TokenInfo, error) {
	// Get user info, user secret by loginID
	userInfo, err := t.userInfoRepoSecondary.GetByLoginID(ctx, loginID)
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Msg("Failed to get user info by login ID")
		return nil, getReturnErr(err)
	}
	userSecret, err := t.userSecretRepoSecondary.Get(ctx, userInfo.ID)
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Msg("Failed to get user secret")
		return nil, getReturnErr(err)
	}

	// Validate login ID, password
	if !password.ValidatePasswd(passwd, userSecret.PasswdHash, userSecret.PasswdSalt) {
		return nil, ErrUnauthorized
	}

	// Create auth token
	tokenInfo, err := token.CreateToken(&token.AuthClaims{UserID: userInfo.ID.String(),
		UserLoginID: userInfo.LoginID, UserRole: userInfo.Role})
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Msg("Failed to get user secret")
		return nil, getReturnErr(err)
	}

	return tokenInfo, nil
}
