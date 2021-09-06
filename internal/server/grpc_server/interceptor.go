package grpc_server

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"

	"github.com/ssup2ket/ssup2ket-auth-service/pkg/authtoken"
)

func icLoggerSetterUnary() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (_ interface{}, err error) {
		// Create logger form global logger and set the logger in the context
		logger := log.With().Logger()
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
		return handler(ctx, req)
	}
}
