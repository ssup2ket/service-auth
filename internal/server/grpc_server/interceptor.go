package grpc_server

import (
	"context"
	"fmt"
	"strings"
	"time"

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

	authtoken "github.com/ssup2ket/ssup2ket-auth-service/pkg/auth/token"
	grpcmeta "github.com/ssup2ket/ssup2ket-auth-service/pkg/grpc/meta"
	"github.com/ssup2ket/ssup2ket-auth-service/pkg/header"
)

func icLoggerSetterUnary() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (_ interface{}, err error) {
		// Create logger form global logger and set the logger in the context
		logger := log.With().Logger()

		// Call next handler with logger
		return handler(logger.WithContext(ctx), req)
	}
}

func icAccessLoggerUary() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (_ interface{}, err error) {
		// Get request ID
		requestID, ok := ctx.Value(header.RequestIDKey).(string)
		if !ok {
			log.Ctx(ctx).Error().Msg("Failed to get request ID from context")
			return nil, getErrServerError()
		}

		// Get span
		span := opentracing.SpanFromContext(ctx)
		if span == nil {
			log.Ctx(ctx).Error().Msg("Failed to get opentracing span from context")
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

func icRequestIdSetterUnary() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (_ interface{}, err error) {
		// Get request ID
		md := grpcmeta.ExtractMetaFromContext(ctx)
		requestID := ""
		requestIDs := md[header.RequestIDHeader]
		if len(requestIDs) != 1 {
			requestID = uuid.NewV4().String()
		} else {
			requestID = requestIDs[0]
		}

		// Set request ID to new context
		newCtx := context.WithValue(ctx, header.RequestIDKey, requestID)

		// Set request ID to logger
		zerolog.Ctx(newCtx).UpdateContext(func(c zerolog.Context) zerolog.Context {
			return c.Str("request_id", requestID)
		})

		// Set request ID to response meta
		header := metadata.Pairs(header.RequestIDHeader, requestID)
		grpc.SetHeader(newCtx, header)

		// Call next handler
		return handler(newCtx, req)
	}
}

func icOpenTracingSetterUnary(t opentracing.Tracer) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (_ interface{}, err error) {
		// Get spancontext from request headers
		md := grpcmeta.ExtractMetaFromContext(ctx)
		spanCtx, err := t.Extract(opentracing.HTTPHeaders, grpcmeta.MetadataReaderWriter{MD: md})
		if err != nil {
			if err != opentracing.ErrSpanContextNotFound {
				log.Ctx(ctx).Error().Err(err).Msg("Failed to get opentracing spancontext")
				return nil, getErrServerError()
			}
		}

		// Start-Finish span
		span, childCtx := opentracing.StartSpanFromContextWithTracer(ctx, t, "auth-service", ext.RPCServerOption(spanCtx))
		defer span.Finish()

		// Set trace ID and span ID to logger
		traceID := span.Context().(jaeger.SpanContext).TraceID().String()
		spanID := span.Context().(jaeger.SpanContext).SpanID().String()
		zerolog.Ctx(childCtx).UpdateContext(func(c zerolog.Context) zerolog.Context {
			return c.Str("trace_id", traceID).Str("span_id", spanID)
		})

		// Set trace ID and span ID to logger
		header := metadata.Pairs(header.TraceIDHeader, traceID, header.SpanIDHeader, spanID)
		grpc.SetHeader(childCtx, header)

		// Call next handler
		return handler(childCtx, req)
	}
}

func icAuthTokenValidaterAndSetterUary() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (_ interface{}, err error) {
		// Pass token validation for some requests
		if info.FullMethod == "/Token/CreateToken" || info.FullMethod == "/User/CreateUser" {
			return handler(ctx, req)
		}

		// Get request's meta data
		md, okMeta := metadata.FromIncomingContext(ctx)
		if !okMeta {
			log.Ctx(ctx).Error().Err(err).Msg("Failed to get metadata for auth token")
			return nil, getErrServerError()
		}

		// Get auth token
		token, okToken := md["authorization"]
		if !okToken {
			log.Ctx(ctx).Error().Msg("Failed to get auth token")
			return nil, getErrUnauthorized()
		}

		// Validate auth token and get auth info
		authInfo, err := authtoken.ValidateToken(token[0])
		if err != nil {
			log.Ctx(ctx).Error().Err(err).Msg("Auth token isn't valid")
			return nil, getErrUnauthorized()
		}

		// Set auth info
		zerolog.Ctx(ctx).UpdateContext(func(c zerolog.Context) zerolog.Context {
			return c.Str("token_user_id", authInfo.UserID).Str("token_user_loginid", authInfo.UserLoginID)
		})

		// Call next handler
		return handler(ctx, req)
	}
}

func icUserIDSetterUary() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (_ interface{}, err error) {
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
