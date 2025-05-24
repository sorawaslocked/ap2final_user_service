package usecase

import (
	"context"
	"github.com/sorawaslocked/ap2final_user_service/internal/model"
)

type UserRepository interface {
	InsertOne(ctx context.Context, user model.User) (model.User, error)
	FindOne(ctx context.Context, filter model.UserFilter) (model.User, error)
	FindMany(ctx context.Context, filter model.UserFilter) ([]model.User, error)
	UpdateOne(ctx context.Context, filter model.UserFilter, update model.User) (model.User, error)
	DeleteOne(ctx context.Context, filter model.UserFilter) (model.User, error)
}
