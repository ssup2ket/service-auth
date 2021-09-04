package grpc_server

import (
	"context"

	"github.com/rs/zerolog/log"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/ssup2ket/ssup2ket-auth-service/internal/domain/service"
)

func (s *ServerGRPC) CreateToken(ctx context.Context, req *TokenCreateRequest) (*TokenInfoResponse, error) {
	// Validate request
	if err := req.validate(); err != nil {
		log.Ctx(ctx).Error().Err(err).Msg("Wrong ID/password format")
		return nil, getErrBadRequest()
	}

	// Create token
	tokenInfo, err := s.domain.Token.CreateToken(ctx, req.LoginId, req.Password)
	if err != nil {
		if err == service.ErrUnauthorized {
			log.Ctx(ctx).Error().Err(err).Msg("Wrong ID/password")
			return nil, getErrUnauthorized()
		}
		log.Ctx(ctx).Error().Err(err).Msg("Failed to create token")
		return nil, getErrServerError()
	}

	return &TokenInfoResponse{Token: tokenInfo.Token,
		IssuedAt:  timestamppb.New(tokenInfo.IssuedAt),
		ExpiresAt: timestamppb.New(tokenInfo.ExpiresAt)}, nil
}

// Request validate
func (t *TokenCreateRequest) validate() error {
	return nil
}
