package grpc_server

import (
	"context"

	"github.com/Teamlix/proto/gen/go/user_service/v1"
	grpc_auth "github.com/grpc-ecosystem/go-grpc-middleware/auth"
	"github.com/sirupsen/logrus"
	grpc_clients "github.com/teamlix/grpc-clients"
	"github.com/teamlix/user-service/internal/domain"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type UserService interface {
	SignUp(ctx context.Context, name, email, password, repeatedPassword string) (domain.Tokens, error)
	SignIn(ctx context.Context, email, password string) (domain.Tokens, error)
	Refresh(ctx context.Context, refreshToken string) (domain.Tokens, error)
	LogOut(ctx context.Context, accessToken, refreshToken string) error
	GetUserByID(ctx context.Context, userID string) (domain.User, error)
	GetUsersList(ctx context.Context, skip, limit int) ([]domain.User, int, error)
	CheckAccessToken(ctx context.Context, accessToken string) error
}

type UserServer struct {
	user_service.UnimplementedUserServiceServer
	service       UserService
	logger        *logrus.Logger
	clients       *grpc_clients.Clients
	publicMethods map[string]struct{}
}

func newUserServer(service UserService, logger *logrus.Logger, clients *grpc_clients.Clients) UserServer {
	publicMethods := make(map[string]struct{})
	publicMethods["/user_service.v1.UserService/SignUp"] = struct{}{}
	publicMethods["/user_service.v1.UserService/SignIn"] = struct{}{}
	publicMethods["/user_service.v1.UserService/GetUserByID"] = struct{}{}
	publicMethods["/user_service.v1.UserService/GetUsersList"] = struct{}{}

	// non gateway methods
	publicMethods["/user_service.v1.UserService/CheckAccessToken"] = struct{}{}

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

func (us UserServer) SignUp(ctx context.Context, req *user_service.SignUpRequest) (*user_service.SignUpResponse, error) {
	name := req.GetName()
	email := req.GetEmail()
	password := req.GetPassword()
	repeatedPassword := req.GetRepeatedPassword()

	us.logger.Debugf("signup request, name: %s, email: %s", name, email)

	t, err := us.service.SignUp(ctx, name, email, password, repeatedPassword)
	if err != nil {
		us.logger.Errorln("signup error: ", err)
		return nil, makeStatusError(err)
	}

	res := user_service.SignUpResponse{
		Result: &user_service.Tokens{
			AccessToken:  t.AccessToken,
			RefreshToken: t.RefreshToken,
		},
	}

	return &res, nil
}

func (us UserServer) SignIn(ctx context.Context, req *user_service.SignInRequest) (*user_service.SignInResponse, error) {
	email := req.GetEmail()
	password := req.GetPassword()

	us.logger.Debugf("signin request, email: %s", email)

	t, err := us.service.SignIn(ctx, email, password)
	if err != nil {
		us.logger.Errorln("signin error: ", err)
		return nil, makeStatusError(err)
	}

	res := user_service.SignInResponse{
		Result: &user_service.Tokens{
			AccessToken:  t.AccessToken,
			RefreshToken: t.RefreshToken,
		},
	}

	return &res, nil
}

func (us UserServer) LogOut(ctx context.Context, req *user_service.LogOutRequest) (*user_service.LogOutResponse, error) {
	at := req.GetAccessToken()
	rt := req.GetRefreshToken()

	err := us.service.LogOut(ctx, at, rt)
	if err != nil {
		us.logger.Errorln("logout error: ", err)
		return nil, makeStatusError(err)
	}

	res := user_service.LogOutResponse{
		Result: true,
	}

	return &res, nil
}

func (us UserServer) Refresh(ctx context.Context, req *user_service.RefreshRequest) (*user_service.RefreshResponse, error) {
	rt := req.GetRefreshToken()

	t, err := us.service.Refresh(ctx, rt)
	if err != nil {
		us.logger.Errorln("refresh error: ", err)
		return nil, makeStatusError(err)
	}

	res := user_service.RefreshResponse{
		Result: &user_service.Tokens{
			AccessToken:  t.AccessToken,
			RefreshToken: t.RefreshToken,
		},
	}

	return &res, nil
}

func (us UserServer) GetUserByID(ctx context.Context, req *user_service.GetUserByIDRequest) (*user_service.GetUserByIDResponse, error) {
	userID := req.GetUserId()

	u, err := us.service.GetUserByID(ctx, userID)
	if err != nil {
		us.logger.Errorln("get user by id error: ", err)
		return nil, makeStatusError(err)
	}

	res := user_service.GetUserByIDResponse{
		Result: &user_service.User{
			Id:   u.ID,
			Name: u.Name,
		},
	}

	return &res, nil
}

func (us UserServer) GetUsersList(ctx context.Context, req *user_service.GetUsersListRequest) (*user_service.GetUsersListResponse, error) {
	skip := req.GetSkip()
	limit := req.GetLimit()

	us.logger.Debugf("get users list request, skip: %d, limit: %d", skip, limit)
	ul, cnt, err := us.service.GetUsersList(ctx, int(skip), int(limit))
	if err != nil {
		us.logger.Errorln("get users list error: ", err)
		return nil, makeStatusError(err)
	}

	resUsers := make([]*user_service.User, 0)
	for _, u := range ul {
		el := user_service.User{
			Id:   u.ID,
			Name: u.Name,
		}
		resUsers = append(resUsers, &el)
	}

	res := user_service.GetUsersListResponse{
		Result: &user_service.GetUsersListResponse_Result{
			Data:       resUsers,
			TotalCount: uint32(cnt),
		},
	}

	return &res, nil
}

func (us UserServer) CheckAccessToken(ctx context.Context, req *user_service.CheckAccessTokenRequest) (*user_service.CheckAccessTokenResponse, error) {
	at := req.GetAccessToken()

	err := us.service.CheckAccessToken(ctx, at)
	if err != nil {
		us.logger.Errorln("check access token error: ", err)
		return nil, makeStatusError(err)
	}

	res := user_service.CheckAccessTokenResponse{
		Result: true,
	}

	return &res, nil
}
