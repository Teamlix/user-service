package grpc_server

import (
	"context"
	"fmt"
	"net"

	"github.com/Teamlix/proto/gen/go/user_service/v1"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_auth "github.com/grpc-ecosystem/go-grpc-middleware/auth"
	"github.com/sirupsen/logrus"
	grpc_clients "github.com/teamlix/grpc-clients"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
)

type Server struct {
	host    string
	port    string
	server  *grpc.Server
	service Service
	logger  *logrus.Logger
	clients *grpc_clients.Clients
}

func NewServer(
	host string,
	port string,
	service Service,
	logger *logrus.Logger,
	clients *grpc_clients.Clients,
) Server {
	return Server{
		host: host,
		port: port,
		server: grpc.NewServer(
			grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
				grpc_auth.UnaryServerInterceptor(func(ctx context.Context) (context.Context, error) { // that function is overrided by AuthFuncOverride
					return ctx, nil
				}),
			)),
		),
		service: service,
		logger:  logger,
		clients: clients,
	}
}

func (s *Server) Serve() error {
	lis, err := net.Listen("tcp", net.JoinHostPort(s.host, s.port))
	if err != nil {
		return fmt.Errorf("can't start listening addr: %w", err)
	}
	grpc_health_v1.RegisterHealthServer(s.server, health.NewServer())
	user_service.RegisterUserServiceServer(s.server, newUserServer(s.service, s.logger, s.clients))
	err = s.server.Serve(lis)

	if err != nil {
		return fmt.Errorf("grpc server: can't start listening: %w", err)
	}
	return nil
}

func (s *Server) Stop() {
	s.server.GracefulStop()
}
