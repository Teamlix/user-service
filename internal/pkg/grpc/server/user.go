package grpc_server

import (
	"context"

	"github.com/Teamlix/proto/gen/go/user_service/v1"
	"github.com/sirupsen/logrus"
	"github.com/teamlix/user-service/internal/service"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type UserServer struct {
	user_service.UnimplementedUserServiceServer
	services *service.Service
	logger   *logrus.Logger
}

func newUserServer(services *service.Service, logger *logrus.Logger) UserServer {
	return UserServer{services: services, logger: logger}
}

func (us UserServer) SignUp(context.Context, *user_service.SignUpRequest) (*user_service.SignUpResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SignUp not implemented")
}
func (us UserServer) SignIn(context.Context, *user_service.SignInRequest) (*user_service.SignInResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SignIn not implemented")
}
func (us UserServer) LogOut(context.Context, *user_service.LogOutRequest) (*user_service.LogOutResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method LogOut not implemented")
}
func (us UserServer) Refresh(context.Context, *user_service.RefreshRequest) (*user_service.RefreshResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Refresh not implemented")
}
func (us UserServer) GetUserByID(context.Context, *user_service.GetUserByIDRequest) (*user_service.GetUserByIDResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetUserByID not implemented")
}
func (us UserServer) GetUsersList(context.Context, *user_service.GetUsersListRequest) (*user_service.GetUsersListResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetUsersList not implemented")
}
