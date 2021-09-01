package grpc_server

import (
	context "context"

	"github.com/rs/zerolog/log"
)

func (s *ServerGRPC) CreateToken(ctx context.Context, req *TokenCreateRequest) (*TokenInfoResponse, error) {
	// Validate request
	if err := req.validate(); err != nil {
		log.Ctx(ctx).Error().Err(err).Msg("Wrong get user request")
		return nil, getErrBadRequest()
	}

	// Create token
	token, err := s.domain.Token.CreateToken(ctx, req.LoginId, req.Password)
	if err != nil {
	}

	return &TokenInfoResponse{Token: token}, nil
}

// Request validate
func (t *TokenCreateRequest) validate() error {
	return nil
}
