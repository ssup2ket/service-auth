package http_server

import (
	"net/http"
	"strings"
	"time"

	"github.com/casbin/casbin"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	opentracing "github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	uuid "github.com/satori/go.uuid"
	"github.com/uber/jaeger-client-go"

	"github.com/ssup2ket/ssup2ket-auth-service/internal/server/middleware"
	authtoken "github.com/ssup2ket/ssup2ket-auth-service/pkg/auth/token"
)

func mwRequestIDSetter() func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()

			// Get request ID
			requestID := r.Header.Get(middleware.HeaderRequestID)
			if requestID == "" {
				requestID = uuid.NewV4().String()
			}

			// Set request ID to new context
			newCtx := middleware.SetRequestIDToCtx(ctx, requestID)

			// Set request ID to logger
			zerolog.Ctx(newCtx).UpdateContext(func(c zerolog.Context) zerolog.Context {
				return c.Str("request_id", requestID)
			})

			// Set request ID to response header
			w.Header().Set(middleware.HeaderRequestID, requestID)

			// Call next handler
			next.ServeHTTP(w, r.WithContext(newCtx))
		}

		return http.HandlerFunc(fn)
	}
}

func mwOpenTracingSetter(t opentracing.Tracer) func(next http.Handler) http.Handler {
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

			// Get trace ID and span ID
			traceID := span.Context().(jaeger.SpanContext).TraceID().String()
			spanID := span.Context().(jaeger.SpanContext).SpanID().String()

			// Set trace ID and span ID to logger
			zerolog.Ctx(childCtx).UpdateContext(func(c zerolog.Context) zerolog.Context {
				return c.Str("trace_id", traceID).Str("span_id", spanID)
			})

			// Set trace ID and span ID to response header
			w.Header().Set(middleware.HeaderTraceID, traceID)
			w.Header().Set(middleware.HeaderSpanID, spanID)

			// Call next handler with child context
			next.ServeHTTP(w, r.WithContext(childCtx))
		}
		return http.HandlerFunc(fn)
	}
}

func mwAccessLogger(r *http.Request, status, size int, duration time.Duration) {
	ctx := r.Context()

	// Get request ID
	requestID, err := middleware.GetRequestIDFromCtx(ctx)
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Msg("Failed to get request ID from context")
		return
	}

	// Get span
	span := opentracing.SpanFromContext(ctx)
	if span == nil {
		log.Ctx(ctx).Error().Msg("Failed to get opentracing span from context")
		return
	}

	// Logging
	zerolog.Ctx(ctx).Info().
		Str("request_id", requestID).
		Str("trace_id", span.Context().(jaeger.SpanContext).TraceID().String()).
		Str("span_id", span.Context().(jaeger.SpanContext).SpanID().String()).
		Str("method", r.Method).Str("url", r.URL.String()).Int("status", status).
		Str("client_ip", r.RemoteAddr).Int("response_size", size).Dur("duration", duration).
		Send()
}

func mwAuthTokenValidatorAndSetter() func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()

			// Get auth token
			token := strings.TrimSpace(strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer"))

			// Validate auth token and get auth info
			authInfo, err := authtoken.ValidateToken(token)
			if err != nil {
				log.Ctx(ctx).Error().Err(err).Msg("Auth token isn't valid")
				render.Render(w, r, getErrRendererUnauthorized())
				return
			}

			// Set auth context to context
			newCtx := middleware.SetUserIDToCtx(ctx, authInfo.UserID)
			newCtx = middleware.SetUserLoginIDToCtx(newCtx, authInfo.UserLoginID)
			newCtx = middleware.SetUserRoleToCtx(newCtx, authInfo.UserRole)

			// Set auth info to logger
			zerolog.Ctx(newCtx).UpdateContext(func(c zerolog.Context) zerolog.Context {
				return c.Str("token_user_id", authInfo.UserID).Str("token_user_loginid", authInfo.UserLoginID)
			})

			// Call next handler
			next.ServeHTTP(w, r.WithContext(newCtx))
		}

		return http.HandlerFunc(fn)
	}
}

func mwAuthorizer(e *casbin.Enforcer) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()

			// Get role, method, path from request
			role, err := middleware.GetUserRoleFromCtx(ctx)
			if err != nil {
				log.Ctx(ctx).Error().Msg("No user role in context")
				render.Render(w, r, getErrRendererServerError())
				return
			}
			method := strings.ToLower(r.Method)
			path := r.URL.Path

			// Check authority
			if !e.Enforce(string(role), path, method) {
				log.Ctx(ctx).Error().Msg("This request isn't allowed")
				render.Render(w, r, getErrRendererUnauthorized())
				return
			}

			// Call next handler
			next.ServeHTTP(w, r)
		}

		return http.HandlerFunc(fn)
	}
}

func mwUserIDLoggerSetter() func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
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
}
