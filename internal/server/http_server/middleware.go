package http_server

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/hlog"
)

func mwAccessLogger(r *http.Request, status, size int, duration time.Duration) {
	hlog.FromRequest(r).Info().
		Str("method", r.Method).Str("url", r.URL.String()).Int("status", status).
		Str("client_ip", r.RemoteAddr).Int("size", size).Dur("duration", duration).
		Send()
}

func mwUserIDSetter(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		// Get User ID from request parameter
		userID := chi.URLParam(r, "UserID")
		if userID == "" {
			next.ServeHTTP(w, r)
			return
		}

		// Set User ID to request context
		zerolog.Ctx(ctx).UpdateContext(func(c zerolog.Context) zerolog.Context {
			return c.Str("user_id", userID)
		})
		next.ServeHTTP(w, r)
	}

	return http.HandlerFunc(fn)
}
