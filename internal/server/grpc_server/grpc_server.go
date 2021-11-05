package grpc_server

import (
	"net"

	"github.com/casbin/casbin"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware/v2"
	grpc_recover "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/recovery"
	opentracing "github.com/opentracing/opentracing-go"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"github.com/ssup2ket/ssup2ket-auth-service/internal/domain"
)

// ServerGRPC
type ServerGRPC struct {
	grpcServer *grpc.Server
	domain     *domain.Domain

	UnimplementedTokenServer
	UnimplementedUserServer
	UnimplementedUserMeServer
}

func New(d *domain.Domain, e *casbin.Enforcer, t opentracing.Tracer) (*ServerGRPC, error) {
	server := ServerGRPC{
		grpcServer: grpc.NewServer(
			grpc_middleware.WithUnaryServerChain(
				grpc_recover.UnaryServerInterceptor(),
				icLoggerSetterUnary(),

				icRequestIdSetterUnary(),
				icOpenTracingSetterUnary(t),
				icAccessLoggerUnary(),

				icAccessTokenValidaterAndSetterUnary(),
				icAuthorizerUnary(e),
				icUserIDLoggerSetterUnary(),
			),
		),
		domain: d,
	}

	// Regist service
	RegisterTokenServer(server.grpcServer, &server)
	RegisterUserServer(server.grpcServer, &server)
	RegisterUserMeServer(server.grpcServer, &server)

	// Set reflection
	reflection.Register(server.grpcServer)

	return &server, nil
}

func (s *ServerGRPC) ListenAndServe() {
	go func() {
		grpcListen, err := net.Listen("tcp", ":9090")
		if err != nil {
			log.Fatal().Err(err).Msg("Failed to listen to run GRPC server")
		}
		if err := s.grpcServer.Serve(grpcListen); err != nil {
			log.Fatal().Err(err).Msg("Failed to serve GRPC server")
		}
	}()
}

func (s *ServerGRPC) Shutdown() {
	s.grpcServer.GracefulStop()
}
