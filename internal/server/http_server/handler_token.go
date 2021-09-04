package http_server

import (
	"net/http"

	"github.com/go-chi/render"
	"github.com/rs/zerolog/log"
	"github.com/ssup2ket/ssup2ket-auth-service/internal/domain/service"
)

// Create a token
func (s *ServerHTTP) PostTokens(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	tokenCreate := TokenCreate{}

	// Unmarshal request
	if err := render.Bind(r, &tokenCreate); err != nil {
		log.Ctx(ctx).Error().Err(err).Msg("Wrong ID/password format")
		render.Render(w, r, getErrRendererBadRequest())
		return
	}

	// Create token
	tokenInfo, err := s.domain.Token.CreateToken(ctx, tokenCreate.LoginId, tokenCreate.Password)
	if err != nil {
		if err == service.ErrUnauthorized {
			log.Ctx(ctx).Error().Err(err).Msg("Wrong ID/password")
			render.Render(w, r, getErrRendererUnauthorized())
			return
		}
		log.Ctx(ctx).Error().Err(err).Msg("Failed to create token")
		render.Render(w, r, getErrRendererServerError())
		return
	}

	render.JSON(w, r, TokenInfo{
		Token:     tokenInfo.Token,
		IssuedAt:  tokenInfo.IssuedAt,
		ExpiresAt: tokenInfo.ExpiresAt})
}

// Validate & Bind
func (u *TokenCreate) Bind(r *http.Request) error {
	return nil
}
