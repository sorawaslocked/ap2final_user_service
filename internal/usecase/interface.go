package usecase

import (
	"context"
	"github.com/sorawaslocked/ap2final_user_service/internal/model"
)

type UserRepository interface {
	InsertOne(ctx context.Context, user model.User) (model.User, error)
	FindOne(ctx context.Context, filter model.UserFilter) (model.User, error)
	FindMany(ctx context.Context, filter model.UserFilter) ([]model.User, error)
	UpdateOne(ctx context.Context, filter model.UserFilter, update model.UserUpdateData) (model.User, error)
	DeleteOne(ctx context.Context, filter model.UserFilter) (model.User, error)
}

type TokenRepository interface {
	InsertOne(ctx context.Context, session model.Session) error
	FindOneByToken(ctx context.Context, token string) (model.Session, error)
	DeleteByToken(ctx context.Context, token string) error
}
