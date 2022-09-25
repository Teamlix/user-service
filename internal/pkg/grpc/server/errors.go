package grpc_server

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func makeStatusError(e error) error {
	var code codes.Code
	switch e.Error() {
	case "username has wrong format":
		code = codes.InvalidArgument
	case "email has wrong format":
		code = codes.InvalidArgument
	case "password has wrong length":
		code = codes.InvalidArgument
	case "passwords are not equal":
		code = codes.InvalidArgument
	case "provided email already exists":
		code = codes.AlreadyExists
	case "provided user name already exists":
		code = codes.AlreadyExists
	case "unauthorized":
		code = codes.Unauthenticated
	case "token expired":
		code = codes.Unauthenticated
	case "user not found":
		code = codes.NotFound
	case "wrong userID":
		code = codes.InvalidArgument
	case "invalid payload":
		code = codes.InvalidArgument
	default:
		code = codes.Internal
	}
	err := status.Error(code, e.Error())
	return err
}
