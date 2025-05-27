package dto

import (
	"github.com/sorawaslocked/ap2final_protos_gen/events"
	"github.com/sorawaslocked/ap2final_user_service/internal/model"
)

func FromUserToRegisterEvent(user model.User) *events.UserRegisterEvent {
	return &events.UserRegisterEvent{
		UserID: user.ID,
		Email:  user.Email,
	}
}
