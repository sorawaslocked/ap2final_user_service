package model

import "errors"

var (
	ErrNotFound            = errors.New("not found")
	ErrRefreshTokenExpired = errors.New("refresh token expired")
	ErrPasswordsDoNotMatch = errors.New("passwords do not match")
	ErrUnauthorized        = errors.New("unauthorized")
	ErrEmptyClaims         = errors.New("empty claims")
	ErrInvalidToken        = errors.New("invalid token")
)
