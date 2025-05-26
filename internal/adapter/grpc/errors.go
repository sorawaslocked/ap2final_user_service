package grpc

import (
	"errors"
	"fmt"
	"github.com/sorawaslocked/ap2final_base/pkg/logger"
	"github.com/sorawaslocked/ap2final_user_service/internal/adapter/grpc/dto"
	"github.com/sorawaslocked/ap2final_user_service/internal/model"
	"log/slog"
)

func logError(log *slog.Logger, op string, err error) {
	switch {
	case errors.Is(err, dto.ErrMissingPasswordArgument):
		warn(log, op, err)
	case errors.Is(err, dto.ErrMissingLoginCredentials):
		warn(log, op, err)
	case errors.Is(err, model.ErrNotFound):
		warn(log, op, err)
	case errors.Is(err, dto.ErrUnauthenticated):
		warn(log, op, err)
	default:
		log.Error(
			fmt.Sprintf("user %s", op),
			logger.Err(err),
		)
	}
}

func warn(log *slog.Logger, op string, err error) {
	log.Warn(
		fmt.Sprintf("user %s", op),
		logger.Err(err),
	)
}
