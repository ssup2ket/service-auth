package grpc_server

import (
	"net"

	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"

	"github.com/ssup2ket/ssup2ket-auth-service/internal/domain"
)

// ServerGRPC
type ServerGRPC struct {
	grpcServer *grpc.Server
	domain     *domain.Domain

	UnimplementedUserServer
	UnimplementedTokenServer
}

func New(d *domain.Domain) (*ServerGRPC, error) {
	server := ServerGRPC{
		grpcServer: grpc.NewServer(),
		domain:     d,
	}

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
