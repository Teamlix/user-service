package grpc_server

import (
	"fmt"
	"net"

	"github.com/Teamlix/proto/gen/go/user_service/v1"
	"github.com/sirupsen/logrus"
	"github.com/teamlix/user-service/internal/service"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
)

type Server struct {
	host     string
	port     string
	server   *grpc.Server
	services *service.Service
	logger   *logrus.Logger
}

func NewServer(
	host string,
	port string,
	service *service.Service,
	logger *logrus.Logger,
) Server {
	return Server{
		host:     host,
		port:     port,
		server:   grpc.NewServer(),
		services: service,
		logger:   logger,
	}
}

func (s *Server) Serve() error {
	lis, err := net.Listen("tcp", net.JoinHostPort(s.host, s.port))
	if err != nil {
		return fmt.Errorf("can't start listening addr: %w", err)
	}
	grpc_health_v1.RegisterHealthServer(s.server, health.NewServer())
	user_service.RegisterUserServiceServer(s.server, newUserServer(s.services, s.logger))
	err = s.server.Serve(lis)

	if err != nil {
		return fmt.Errorf("grpc server: can't start listening: %w", err)
	}
	return nil
}

func (s *Server) Stop() {
	s.server.GracefulStop()
}