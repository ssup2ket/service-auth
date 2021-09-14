package http_server

import (
	"net/http"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	opentracing "github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/uber/jaeger-client-go"

	"github.com/ssup2ket/ssup2ket-auth-service/pkg/authtoken"
)

func mwOpenTracingTracerSetter(t opentracing.Tracer) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()

			// Get spancontext from request headers
			spanCtx, err := t.Extract(opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(r.Header))
			if err != nil {
				if err != opentracing.ErrSpanContextNotFound {
					log.Ctx(ctx).Error().Err(err).Msg("Failed to get opentracing spancontext")
					render.Render(w, r, getErrRendererServerError())
					return
				}
			}

			// Start-Finish span
			span, childCtx := opentracing.StartSpanFromContextWithTracer(ctx, t, "auth-service", ext.RPCServerOption(spanCtx))
			defer span.Finish()

			// Set
			zerolog.Ctx(childCtx).UpdateContext(func(c zerolog.Context) zerolog.Context {
				return c.Str("trace_id", span.Context().(jaeger.SpanContext).TraceID().String()).
					Str("span_id", span.Context().(jaeger.SpanContext).SpanID().String())
			})

			// Call next handler with child context
			next.ServeHTTP(w, r.WithContext(childCtx))
		}
		return http.HandlerFunc(fn)
	}
}

func mwAccessLogger(r *http.Request, status, size int, duration time.Duration) {
	ctx := r.Context()
	zerolog.Ctx(ctx).Info().
		Str("method", r.Method).Str("url", r.URL.String()).Int("status", status).
		Str("client_ip", r.RemoteAddr).Int("response_size", size).Dur("duration", duration).
		Send()
}

func mwAuthTokenValidatorAndSetter(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		// Get auth token
		token := strings.TrimSpace(strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer"))

		// Validate auth token and get auth info
		authInfo, err := authtoken.ValidateAuthToken(token)
		if err != nil {
			log.Ctx(ctx).Error().Err(err).Msg("Auth token isn't valid")
			render.Render(w, r, getErrRendererUnauthorized())
			return
		}

		// Set auth info
		zerolog.Ctx(ctx).UpdateContext(func(c zerolog.Context) zerolog.Context {
			return c.Str("token_user_id", authInfo.UserID).Str("token_user_loginid", authInfo.UserLoginID)
		})

		// Call next handler
		next.ServeHTTP(w, r)
	}

	return http.HandlerFunc(fn)
}

func mwUserIDSetter(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		// Get user ID from request parameter
		userID := chi.URLParam(r, "UserID")
		if userID == "" {
			next.ServeHTTP(w, r)
			return
		}

		// Set User ID to request context
		zerolog.Ctx(ctx).UpdateContext(func(c zerolog.Context) zerolog.Context {
			return c.Str("user_id", userID)
		})

		// Call next handler
		next.ServeHTTP(w, r)
	}

	return http.HandlerFunc(fn)
}
