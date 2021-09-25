package grpc_server

import (
	"net"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware/v2"
	grpc_recover "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/recovery"
	opentracing "github.com/opentracing/opentracing-go"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"

	"github.com/ssup2ket/ssup2ket-auth-service/internal/domain"
)

// ServerGRPC
type ServerGRPC struct {
	grpcServer *grpc.Server
	domain     *domain.Domain

	UnimplementedTokenServer
	UnimplementedUserServer
}

func New(d *domain.Domain, t opentracing.Tracer) (*ServerGRPC, error) {
	server := ServerGRPC{
		grpcServer: grpc.NewServer(
			grpc_middleware.WithUnaryServerChain(
				grpc_recover.UnaryServerInterceptor(),
				icLoggerSetterUnary(),

				icRequestIdSetterUnary(),
				icOpenTracingSetterUnary(t),
				icAccessLoggerUary(),

				icAuthTokenValidaterAndSetterUary(),
				icUserIDLoggerSetterUary(),
			),
		),
		domain: d,
	}

	RegisterTokenServer(server.grpcServer, &server)
	RegisterUserServer(server.grpcServer, &server)

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
