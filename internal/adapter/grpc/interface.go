package grpc

import (
	"context"
	"github.com/sorawaslocked/ap2final_user_service/internal/model"
)

type UserUseCase interface {
	Register(ctx context.Context, user model.User) (model.User, error)
	Login(ctx context.Context, user model.User) (model.Token, error)
	RefreshToken(ctx context.Context, refreshToken string) (model.Token, error)
	GetByID(ctx context.Context, token model.Token, id string) (model.User, error)
	UpdateByID(
		ctx context.Context,
		token model.Token,
		id string,
		credentialsUpdate model.UserCredentialUpdateData,
		update model.UserUpdateData,
	) (model.User, error)
	DeleteByID(ctx context.Context, token model.Token, id string) (model.User, error)
}
