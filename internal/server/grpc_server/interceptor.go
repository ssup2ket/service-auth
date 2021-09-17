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
	"github.com/uber/jaeger-client-go"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"

	"github.com/ssup2ket/ssup2ket-auth-service/pkg/authtoken"
	"github.com/ssup2ket/ssup2ket-auth-service/pkg/grpcmeta"
)

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
		zerolog.Ctx(childCtx).UpdateContext(func(c zerolog.Context) zerolog.Context {
			return c.Str("trace_id", span.Context().(jaeger.SpanContext).TraceID().String()).
				Str("span_id", span.Context().(jaeger.SpanContext).SpanID().String())
		})

		// Call next handler
		return handler(ctx, req)
	}
}

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
		log.Ctx(ctx).Info().Str("method", info.FullMethod).Str("code", status.Code(err).String()).
			Int("response_size", proto.Size(respMsg)).Str("client_ip", clientIp).
			Str("duration", duration).Send()

		// Call next handler
		return resp, err
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
		authInfo, err := authtoken.ValidateAuthToken(token[0])
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
