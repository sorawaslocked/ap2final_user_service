package model

import "errors"

var (
	ErrNotFound            = errors.New("not found")
	ErrRefreshTokenExpired = errors.New("refresh token expired")
	ErrPasswordsDoNotMatch = errors.New("passwords do not match")
)
