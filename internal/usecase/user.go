package usecase

import (
	"context"
	"github.com/sorawaslocked/ap2final_user_service/internal/model"
)

type User struct {
	repo UserRepository
}

func NewUser(repo UserRepository) *User {
	return &User{
		repo: repo,
	}
}

func (uc *User) Register(ctx context.Context, user model.User) (model.User, error) {
	panic("implement me")
}

func (uc *User) Login(ctx context.Context, user model.User) (model.Token, error) {
	panic("implement me")
}

func (uc *User) Logout(ctx context.Context, token model.Token) error {
	panic("implement me")
}

func (uc *User) GetByID(ctx context.Context, token model.Token, id string) (model.User, error) {
	panic("implement me")
}

func (uc *User) Update(ctx context.Context, token model.Token, user model.User) (model.User, error) {
	panic("implement me")
}

func (uc *User) DeleteByID(ctx context.Context, token model.Token, id string) (model.User, error) {
	panic("implement me")
}
