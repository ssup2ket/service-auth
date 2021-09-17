package grpc_server

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/ssup2ket/ssup2ket-auth-service/internal/server/errors"
)

func getErrBadRequest() error {
	return status.Error(codes.InvalidArgument, errors.CodeBadRequest)
}

func getErrUnauthorized() error {
	return status.Error(codes.PermissionDenied, errors.CodeUnauthorized)
}

func getErrNotFound(res errors.ErrResouce) error {
	errCode := errors.CodeNotFound
	switch res {
	case errors.ErrResouceUser:
		errCode = errors.CodeNotFoundUser
	}

	return status.Error(codes.NotFound, errCode)
}

func getErrConflict(res errors.ErrResouce) error {
	errCode := errors.CodeNotFound
	switch res {
	case errors.ErrResouceUser:
		errCode = errors.CodeConflictUser
	}

	return status.Error(codes.AlreadyExists, errCode)
}

func getErrServerError() error {
	return status.Error(codes.Unknown, errors.CodeServerError)
}
