package model

import "time"

type User struct {
	ID           string
	FirstName    string
	LastName     string
	Email        string
	PhoneNumber  string
	Password     string
	PasswordHash string
	Role         string
	CreatedAt    time.Time
	UpdatedAt    time.Time

	IsDeleted bool
	IsActive  bool
}

type UserFilter struct {
	ID           *string
	FirstName    *string
	LastName     *string
	Email        *string
	PhoneNumber  *string
	PasswordHash *string
	Role         *string

	IsDeleted *bool
	IsActive  *bool
}

type UserCredentialUpdateData struct {
	CurrentPassword string
	NewPassword     string
}

type UserUpdateData struct {
	FirstName    *string
	LastName     *string
	Email        *string
	PhoneNumber  *string
	PasswordHash *string
	Role         *string
	UpdatedAt    time.Time

	IsDeleted *bool
	IsActive  *bool
}
