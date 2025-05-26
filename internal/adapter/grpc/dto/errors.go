package dto

import (
	"errors"
	"github.com/sorawaslocked/ap2final_user_service/internal/model"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	ErrMissingPasswordArgument = errors.New("provide both new and old passwords")
	ErrMissingLoginCredentials = errors.New("provide login credentials")
	ErrUnauthenticated         = errors.New("unauthenticated")
)

func FromError(err error) error {
	switch {
	case errors.Is(err, ErrMissingPasswordArgument):
		return status.Error(codes.InvalidArgument, err.Error())
	case errors.Is(err, ErrMissingLoginCredentials):
		return status.Error(codes.InvalidArgument, err.Error())
	case errors.Is(err, model.ErrNotFound):
		return status.Error(codes.NotFound, err.Error())
	case errors.Is(err, ErrUnauthenticated):
		return status.Error(codes.Unauthenticated, err.Error())
	case errors.Is(err, model.ErrPasswordsDoNotMatch):
		return status.Error(codes.InvalidArgument, err.Error())
	case errors.Is(err, model.ErrRefreshTokenExpired):
		return status.Error(codes.DeadlineExceeded, err.Error())
	default:
		return status.Error(codes.Internal, "something went wrong")
	}
}
