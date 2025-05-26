package dto

import (
	"github.com/sorawaslocked/ap2final_protos_gen/base"
	svc "github.com/sorawaslocked/ap2final_protos_gen/service/user"
	"github.com/sorawaslocked/ap2final_user_service/internal/model"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func ToUserFromRegisterRequest(req *svc.RegisterRequest) model.User {
	return model.User{
		Email:    req.Email,
		Password: req.Password,
	}
}

func ToUserFromLoginRequest(req *svc.LoginRequest) model.User {
	return model.User{
		Email:    req.Email,
		Password: req.Password,
	}
}

func ToUserUpdateFromUpdateRequest(req *svc.UpdateRequest) (
	string,
	model.UserUpdateData,
	model.UserCredentialUpdateData,
	error,
) {
	credentialsUpdate := model.UserCredentialUpdateData{}

	if req.CurrentPassword != nil && req.NewPassword != nil {
		credentialsUpdate.CurrentPassword = *req.CurrentPassword
		credentialsUpdate.NewPassword = *req.NewPassword
	} else if req.CurrentPassword == nil && req.NewPassword == nil {
		credentialsUpdate.CurrentPassword = ""
		credentialsUpdate.NewPassword = ""
	} else {
		return "", model.UserUpdateData{}, model.UserCredentialUpdateData{}, ErrMissingPasswordArgument
	}

	update := model.UserUpdateData{
		FirstName:   req.FirstName,
		LastName:    req.LastName,
		Email:       req.Email,
		PhoneNumber: req.PhoneNumber,
		Role:        req.Role,
		IsDeleted:   req.IsDeleted,
		IsActive:    req.IsActive,
	}

	return req.ID, update, credentialsUpdate, nil
}

func FromUserToPb(user model.User) *base.User {
	return &base.User{
		ID:           user.ID,
		FirstName:    user.FirstName,
		LastName:     user.LastName,
		Email:        user.Email,
		PhoneNumber:  user.PhoneNumber,
		PasswordHash: user.PasswordHash,
		Role:         user.Role,
		CreatedAt:    timestamppb.New(user.CreatedAt),
		UpdatedAt:    timestamppb.New(user.UpdatedAt),
		IsDeleted:    user.IsDeleted,
		IsActive:     user.IsActive,
	}
}
