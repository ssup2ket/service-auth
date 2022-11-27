package grpc_server

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/casbin/casbin"
	opentracing "github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	uuid "github.com/satori/go.uuid"
	"github.com/uber/jaeger-client-go"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"

	"github.com/ssup2ket/service-auth/internal/server/middleware"
	authtoken "github.com/ssup2ket/service-auth/pkg/auth/token"
	grpcmeta "github.com/ssup2ket/service-auth/pkg/grpc/meta"
)

func icLoggerSetterUnary() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		// Create logger form global logger and set the logger in the context
		logger := log.With().Logger()

		// Call next handler with logger
		return handler(logger.WithContext(ctx), req)
	}
}

func icRequestIdSetterUnary() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		// Get request ID
		md := grpcmeta.ExtractMetaFromContext(ctx)
		requestID := ""
		requestIDs := md[middleware.HeaderRequestID]
		if len(requestIDs) != 1 {
			requestID = uuid.NewV4().String()
		} else {
			requestID = requestIDs[0]
		}

		// Set request ID to new context
		newCtx := middleware.SetRequestIDToCtx(ctx, requestID)

		// Set request ID to logger
		zerolog.Ctx(newCtx).UpdateContext(func(c zerolog.Context) zerolog.Context {
			return c.Str("request_id", requestID)
		})

		// Set request ID to response meta
		header := metadata.Pairs(middleware.HeaderRequestID, requestID)
		grpc.SetHeader(newCtx, header)

		// Call next handler
		return handler(newCtx, req)
	}
}

func icOpenTracingSetterUnary() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		// Get tracer
		tracer := opentracing.GlobalTracer()

		// Get spancontext from request headers
		md := grpcmeta.ExtractMetaFromContext(ctx)
		spanCtx, err := tracer.Extract(opentracing.HTTPHeaders, grpcmeta.MetadataReaderWriter{MD: md})
		if err != nil {
			if err != opentracing.ErrSpanContextNotFound {
				log.Ctx(ctx).Error().Err(err).Msg("Failed to get opentracing spancontext")
				return nil, getErrServerError()
			}
		}

		// Start entry span
		span, childCtx := opentracing.StartSpanFromContextWithTracer(ctx, tracer, "entry-grpc", ext.RPCServerOption(spanCtx))
		defer span.Finish()

		// Get trace ID and span ID
		traceID := span.Context().(jaeger.SpanContext).TraceID().String()
		spanID := span.Context().(jaeger.SpanContext).SpanID().String()

		// Set trace ID and span ID to logger
		zerolog.Ctx(childCtx).UpdateContext(func(c zerolog.Context) zerolog.Context {
			return c.Str("trace_id", traceID).Str("span_id", spanID)
		})

		// Set trace ID to response meta
		header := metadata.Pairs(middleware.HeaderTraceID, traceID)
		grpc.SetHeader(childCtx, header)

		// Call next handler
		return handler(childCtx, req)
	}
}

func icAccessLoggerUnary() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		// Get request ID
		requestID, err := middleware.GetRequestIDFromCtx(ctx)
		if err != nil {
			log.Ctx(ctx).Error().Err(err).Msg("Failed to get request ID from context")
			return nil, getErrServerError()
		}

		// Get span
		span, ctx := opentracing.StartSpanFromContext(ctx, "icAccessLoggerUnary")
		if span == nil {
			log.Ctx(ctx).Error().Msg("Failed to start opentracing span")
			return nil, getErrServerError()
		}

		// Call next handler and calculate duration
		startTime := time.Now()
		resp, err := handler(ctx, req)
		duration := fmt.Sprintf("%f", time.Since(startTime).Seconds())

		// Get response message size
		respMsg, _ := resp.(proto.Message)
		proto.Size(respMsg)

		// Get client ip
		clientPeer, _ := peer.FromContext(ctx)
		clientIp := clientPeer.Addr.String()

		// Logging
		log.Ctx(ctx).Info().
			Str("request_id", requestID).
			Str("trace_id", span.Context().(jaeger.SpanContext).TraceID().String()).
			Str("span_id", span.Context().(jaeger.SpanContext).SpanID().String()).
			Str("method", info.FullMethod).Str("code", status.Code(err).String()).
			Int("response_size", proto.Size(respMsg)).Str("client_ip", clientIp).
			Str("duration", duration).Send()
		return resp, err
	}
}

func icAccessTokenValidaterAndSetterUnary() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		// Pass token validation for some requests
		if info.FullMethod == "/Token/LoginToken" || info.FullMethod == "/Token/RefreshToken" || info.FullMethod == "/User/CreateUser" {
			return handler(ctx, req)
		}

		// Get request's meta data
		md, okMeta := metadata.FromIncomingContext(ctx)
		if !okMeta {
			log.Ctx(ctx).Error().Msg("Failed to get metadata for access token")
			return nil, getErrServerError()
		}

		// Get access token
		tokens, okToken := md["authorization"]
		if !okToken || len(tokens) != 1 {
			log.Ctx(ctx).Error().Msg("Failed to get access token")
			return nil, getErrUnauthorized()
		}
		token := strings.TrimSpace(strings.TrimPrefix(tokens[0], "Bearer"))

		// Validate access token and get auth info
		authInfo, err := authtoken.ValidateAccessToken(token)
		if err != nil {
			log.Ctx(ctx).Error().Err(err).Msg("Access token isn't valid")
			return nil, getErrUnauthorized()
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
		return handler(newCtx, req)
	}
}

func icAuthorizerUnary(e *casbin.Enforcer) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		// Pass token validation for some requests
		if info.FullMethod == "/Token/LoginToken" || info.FullMethod == "/Token/RefreshToken" || info.FullMethod == "/User/CreateUser" {
			return handler(ctx, req)
		}

		// Get role from context
		role, err := middleware.GetUserRoleFromCtx(ctx)
		if err != nil {
			log.Ctx(ctx).Error().Msg("No user role in context")
			return nil, getErrServerError()
		}

		// Get object and action from method info
		// Method Format : /[Service(Object)]/[Action][Service[Object]] ex) /Token/CreateToken
		tokens := strings.Split(info.FullMethod, "/")
		if len(tokens) != 3 {
			log.Ctx(ctx).Error().Msg("Invalid method format")
			return nil, getErrServerError()
		}
		object := strings.ToLower(tokens[1])
		action := strings.TrimSuffix(strings.ToLower(tokens[2]), object)

		// Check authority
		if !e.Enforce(string(role), object, action) {
			log.Ctx(ctx).Error().Msg("This request isn't allowed")
			return nil, getErrUnauthorized()
		}

		// Call next handler
		return handler(ctx, req)
	}
}

func icUserIDLoggerSetterUnary() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		// If not user service, skip this interceptor
		if !strings.HasPrefix(info.FullMethod, "/User/") {
			return handler(ctx, req)
		}

		// Get user id
		reqMsg, _ := req.(proto.Message)
		userIDField := reqMsg.ProtoReflect().Descriptor().Fields().ByName("id") // id field has user id
		if userIDField == nil {
			// if not exist id field, run next handler
			return handler(ctx, req)
		}

		// Set user ID to request context
		zerolog.Ctx(ctx).UpdateContext(func(c zerolog.Context) zerolog.Context {
			return c.Str("user_id", reqMsg.ProtoReflect().Get(userIDField).String())
		})

		// Call next handler
		return handler(ctx, req)
	}
}
