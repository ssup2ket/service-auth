package service

import (
	"context"

	"github.com/rs/zerolog/log"

	"github.com/ssup2ket/ssup2ket-auth-service/internal/domain/entity"
	"github.com/ssup2ket/ssup2ket-auth-service/internal/domain/repo"
	"github.com/ssup2ket/ssup2ket-auth-service/pkg/auth/hashing"
	"github.com/ssup2ket/ssup2ket-auth-service/pkg/auth/token"
	"github.com/ssup2ket/ssup2ket-auth-service/pkg/entity/uuid"
)

// Token service
type TokenService interface {
	CreateTokens(ctx context.Context, loginID, passwd string) (*token.TokenInfo, *token.TokenInfo, error)
	RefreshToken(ctx context.Context, refreshToken string) (*token.TokenInfo, error)
}

type TokenServiceImp struct {
	userInfoRepoSecondary   repo.UserInfoRepo
	userSecretRepoPrimary   repo.UserSecretRepo
	userSecretRepoSecondary repo.UserSecretRepo
}

func NewTokenServiceImp(userInfoSecondary repo.UserInfoRepo, userSecretPrimary, userSecretSecondary repo.UserSecretRepo) *TokenServiceImp {
	return &TokenServiceImp{
		userInfoRepoSecondary:   userInfoSecondary,
		userSecretRepoPrimary:   userSecretPrimary,
		userSecretRepoSecondary: userSecretSecondary,
	}
}

func (t *TokenServiceImp) CreateTokens(ctx context.Context, loginID, passwd string) (*token.TokenInfo, *token.TokenInfo, error) {
	// Get user info, user secret by loginID
	userInfo, err := t.userInfoRepoSecondary.GetByLoginID(ctx, loginID)
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Msg("Failed to get user info by login ID")
		return nil, nil, getReturnErr(err)
	}
	userSecret, err := t.userSecretRepoSecondary.Get(ctx, userInfo.ID)
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Msg("Failed to get user secret")
		return nil, nil, getReturnErr(err)
	}

	// Validate login ID, password
	if !hashing.ValidateStr(passwd, userSecret.PasswdHash, userSecret.PasswdSalt) {
		return nil, nil, ErrUnauthorized
	}

	// Create access, refresh token
	accTokenInfo, err := token.CreateAccessToken(&token.AuthClaims{UserID: userInfo.ID.String(),
		UserLoginID: userInfo.LoginID, UserRole: userInfo.Role})
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Msg("Failed to create access token")
		return nil, nil, getReturnErr(err)
	}
	refTokenInfo, err := token.CreateRefreshToken(&token.AuthClaims{UserID: userInfo.ID.String(),
		UserLoginID: userInfo.LoginID, UserRole: userInfo.Role})
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Msg("Failed to create refresh token")
		return nil, nil, getReturnErr(err)
	}

	// Update refresh token to DB
	hash, salt, err := hashing.GetStrHashAndSalt(refTokenInfo.Token)
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Msg("Failed to create refresh token's hash and salt")
		return nil, nil, getReturnErr(err)
	}
	t.userSecretRepoPrimary.Update(ctx, &entity.UserSecret{
		ID:               userInfo.ID,
		RefreshTokenHash: hash,
		RefreshTokenSalt: salt,
	})

	return accTokenInfo, refTokenInfo, nil
}

func (t *TokenServiceImp) RefreshToken(ctx context.Context, refreshToken string) (*token.TokenInfo, error) {
	// Validate refresh token and get auth info
	authInfo, err := token.ValidateRefreshToken(refreshToken)
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Msg("Refresh token isn't valid")
		return nil, getReturnErr(err)
	}

	// Get user secret
	userSecret, err := t.userSecretRepoSecondary.Get(ctx, uuid.FromStringOrNil(authInfo.UserID))
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Msg("Failed to get user secret")
		return nil, getReturnErr(err)
	}

	// Check whether the refresh token matches in the DB
	if !hashing.ValidateStr(refreshToken, userSecret.RefreshTokenHash, userSecret.RefreshTokenSalt) {
		log.Ctx(ctx).Error().Err(err).Msg("Refresh token isn't matched")
		return nil, ErrUnauthorized
	}

	// Create access token
	accTokenInfo, err := token.CreateAccessToken(&token.AuthClaims{UserID: authInfo.UserID,
		UserLoginID: authInfo.UserLoginID, UserRole: authInfo.UserRole})
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Msg("Failed to create access token")
		return nil, getReturnErr(err)
	}

	return accTokenInfo, nil
}
