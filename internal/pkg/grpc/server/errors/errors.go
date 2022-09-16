package grpc_errors

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func GetGrpcError(e error) error {
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
	default:
		code = codes.Internal
	}
	err := status.Error(code, e.Error())
	return err
}
