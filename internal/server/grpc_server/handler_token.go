package grpc_server

import (
	"context"

	"github.com/rs/zerolog/log"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/ssup2ket/service-auth/internal/domain/service"
	"github.com/ssup2ket/service-auth/internal/server/errors"
)

func (s *ServerGRPC) LoginToken(ctx context.Context, req *TokenLoginRequest) (*TokenInfosResponse, error) {
	// Get loginID, password from request's meta data
	md, okMeta := metadata.FromIncomingContext(ctx)
	if !okMeta {
		log.Ctx(ctx).Error().Msg("Failed to get metadata for id, password")
		return nil, getErrServerError()
	}
	loginIDs, okLoginID := md["username"]
	if !okLoginID {
		log.Ctx(ctx).Error().Msg("Failed to get login ID from metadata")
		return nil, getErrUnauthorized()
	}
	loginID := loginIDs[0]
	passwords, okPassword := md["password"]
	if !okPassword {
		log.Ctx(ctx).Error().Msg("Failed to get password from metadata")
		return nil, getErrUnauthorized()
	}
	password := passwords[0]

	// Create token
	accTokenInfo, refTokenInfo, err := s.domain.Token.CreateTokens(ctx, loginID, password)
	if err != nil {
		if err == service.ErrRepoNotFound {
			log.Ctx(ctx).Error().Err(err).Msg("ID doesn't exists")
			return nil, getErrNotFound(errors.ErrResouceUser)
		}
		if err == service.ErrUnauthorized {
			log.Ctx(ctx).Error().Err(err).Msg("Wrong ID/password")
			return nil, getErrUnauthorized()
		}
		log.Ctx(ctx).Error().Err(err).Msg("Failed to create access, refresh tokens")
		return nil, getErrServerError()
	}

	return &TokenInfosResponse{
		AccessToken: &TokenInfoResponse{
			Token:     accTokenInfo.Token,
			IssuedAt:  timestamppb.New(accTokenInfo.IssuedAt),
			ExpiresAt: timestamppb.New(accTokenInfo.ExpiresAt),
		},
		RefreshToken: &TokenInfoResponse{
			Token:     refTokenInfo.Token,
			IssuedAt:  timestamppb.New(refTokenInfo.IssuedAt),
			ExpiresAt: timestamppb.New(refTokenInfo.ExpiresAt),
		},
	}, nil
}

func (s *ServerGRPC) RefreshToken(ctx context.Context, req *TokenRefreshRequest) (*TokenInfoResponse, error) {
	// Get refresh token
	refreshToken := req.RefreshToken

	// Refresh token
	refTokenInfo, err := s.domain.Token.RefreshToken(ctx, refreshToken)
	if err != nil {
		if err == service.ErrUnauthorized {
			log.Ctx(ctx).Error().Err(err).Msg("Wrong refresh token")
			return nil, getErrUnauthorized()
		}
		log.Ctx(ctx).Error().Err(err).Msg("Failed to refresh token")
		return nil, getErrServerError()
	}

	return &TokenInfoResponse{
		Token:     refTokenInfo.Token,
		IssuedAt:  timestamppb.New(refTokenInfo.IssuedAt),
		ExpiresAt: timestamppb.New(refTokenInfo.ExpiresAt),
	}, nil
}
