package http_server

import (
	"net/http"

	"github.com/go-chi/render"
	"github.com/rs/zerolog/log"
)

// Create a token
func (s *ServerHTTP) PostTokens(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	tokenCreate := TokenCreate{}

	// Unmarshal request
	if err := render.Bind(r, &tokenCreate); err != nil {
		log.Ctx(ctx).Error().Err(err).Msg("Wrong ID/Password format")
		render.Render(w, r, getErrRendererBadRequest())
		return
	}

	// Create token
	token, err := s.domain.Token.CreateToken(ctx, tokenCreate.LoginId, tokenCreate.Password)
	if err != nil {
	}

	render.JSON(w, r, TokenInfo{Token: token})
}

// Validate & Bind
func (u *TokenCreate) Bind(r *http.Request) error {
	return nil
}
