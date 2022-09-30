package grpc_server

import (
	"context"

	"github.com/Teamlix/proto/gen/go/user_service/v1"
	grpc_auth "github.com/grpc-ecosystem/go-grpc-middleware/auth"
	"github.com/sirupsen/logrus"
	grpc_clients "github.com/teamlix/grpc-clients"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Service interface {
}

type UserServer struct {
	user_service.UnimplementedUserServiceServer
	service       Service
	logger        *logrus.Logger
	clients       *grpc_clients.Clients
	publicMethods map[string]struct{}
}

func newUserServer(service Service, logger *logrus.Logger, clients *grpc_clients.Clients) UserServer {
	publicMethods := make(map[string]struct{})

	return UserServer{service: service, logger: logger, publicMethods: publicMethods, clients: clients}
}

func (us UserServer) AuthFuncOverride(ctx context.Context, fullMethodName string) (context.Context, error) {
	if _, ok := us.publicMethods[fullMethodName]; !ok {
		token, err := grpc_auth.AuthFromMD(ctx, "bearer")
		if err != nil {
			return ctx, err
		}
		ok, err := us.clients.User.CheckAccessToken(ctx, token)
		if err != nil {
			return ctx, err
		}
		if !ok {
			return ctx, status.Errorf(codes.Unauthenticated, "bad authorization string")
		}
	}
	return ctx, nil
}
