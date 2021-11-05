package http_server

import (
	"net/http"

	"github.com/go-chi/render"
	"github.com/rs/zerolog/log"

	"github.com/ssup2ket/ssup2ket-auth-service/internal/domain/service"
)

// Login
func (s *ServerHTTP) PostTokensLogin(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Get login ID and password
	loginID, password, ok := r.BasicAuth()
	if !ok {
		log.Ctx(ctx).Error().Msg("No ID/password info")
		render.Render(w, r, getErrRendererUnauthorized())
		return
	}

	// Create token
	accTokenInfo, refTokenInfo, err := s.domain.Token.CreateTokens(ctx, loginID, password)
	if err != nil {
		if err == service.ErrRepoNotFound {
			log.Ctx(ctx).Error().Err(err).Msg("ID doesn't exists")
			render.Render(w, r, getErrRendererUnauthorized())
			return
		} else if err == service.ErrUnauthorized {
			log.Ctx(ctx).Error().Err(err).Msg("Wrong ID/password")
			render.Render(w, r, getErrRendererUnauthorized())
			return
		}
		log.Ctx(ctx).Error().Err(err).Msg("Failed to create access, refresh tokens")
		render.Render(w, r, getErrRendererServerError())
		return
	}

	render.JSON(w, r, TokenInfos{
		AccessToken: TokenInfo{
			Token:     accTokenInfo.Token,
			IssuedAt:  accTokenInfo.IssuedAt,
			ExpiresAt: accTokenInfo.ExpiresAt,
		},
		RefreshToken: TokenInfo{
			Token:     refTokenInfo.Token,
			IssuedAt:  refTokenInfo.IssuedAt,
			ExpiresAt: refTokenInfo.ExpiresAt,
		},
	})
}

func (s *ServerHTTP) PostTokensRefresh(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	tokenRefresh := TokenRefresh{}

	// Unmarshal request
	if err := render.Bind(r, &tokenRefresh); err != nil {
		log.Ctx(ctx).Error().Err(err).Msg("Wrong create user request")
		render.Render(w, r, getErrRendererBadRequest())
		return
	}

	// Refresh token
	refTokenInfo, err := s.domain.Token.RefreshToken(ctx, tokenRefresh.RefreshToken)
	if err != nil {
		if err == service.ErrUnauthorized {
			log.Ctx(ctx).Error().Err(err).Msg("Wrong refresh token")
			render.Render(w, r, getErrRendererUnauthorized())
			return
		}
		render.Render(w, r, getErrRendererServerError())
		log.Ctx(ctx).Error().Err(err).Msg("Failed to refresh token")
		return
	}

	render.JSON(w, r, TokenInfo{
		Token:     refTokenInfo.Token,
		IssuedAt:  refTokenInfo.IssuedAt,
		ExpiresAt: refTokenInfo.ExpiresAt,
	})
}

func (u *TokenRefresh) Bind(r *http.Request) error {
	return nil
}
