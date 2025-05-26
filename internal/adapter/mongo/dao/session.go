package dao

import (
	"github.com/sorawaslocked/ap2final_user_service/internal/model"
	"time"
)

type Session struct {
	UserID       string    `bson:"userID"`
	RefreshToken string    `bson:"refreshToken"`
	ExpiresAt    time.Time `bson:"expiresAt"`
	CreatedAt    time.Time `bson:"createdAt"`
}

func FromSession(session model.Session) Session {
	return Session{
		UserID:       session.UserID,
		RefreshToken: session.RefreshToken,
		ExpiresAt:    session.ExpiresAt,
		CreatedAt:    session.CreatedAt,
	}
}

func ToSession(session Session) model.Session {
	return model.Session{
		UserID:       session.UserID,
		RefreshToken: session.RefreshToken,
		ExpiresAt:    session.ExpiresAt,
		CreatedAt:    session.CreatedAt,
	}
}
