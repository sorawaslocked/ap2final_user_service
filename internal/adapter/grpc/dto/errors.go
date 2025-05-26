package dto

import (
	"errors"
	"github.com/sorawaslocked/ap2final_user_service/internal/model"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	ErrMissingPasswordArgument = errors.New("provide both new and old passwords")
)

func FromError(err error) error {
	switch {
	case errors.Is(err, ErrMissingPasswordArgument):
		return status.Error(codes.InvalidArgument, err.Error())
	case errors.Is(err, model.ErrNotFound):
		return status.Error(codes.NotFound, err.Error())
	default:
		return status.Error(codes.Internal, "something went wrong")
	}
}
