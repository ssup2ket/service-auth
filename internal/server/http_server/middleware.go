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

func mwUserUUIDSetter(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		// Get User UUID from request parameter
		userUUID := chi.URLParam(r, "UserUUID")
		if userUUID == "" {
			next.ServeHTTP(w, r)
			return
		}

		// Set User UUID to request context
		zerolog.Ctx(ctx).UpdateContext(func(c zerolog.Context) zerolog.Context {
			return c.Str("user_uuid", userUUID)
		})
		next.ServeHTTP(w, r)
	}

	return http.HandlerFunc(fn)
}
