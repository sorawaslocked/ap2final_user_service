package dao

import (
	"github.com/sorawaslocked/ap2final_user_service/internal/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type User struct {
	ID           primitive.ObjectID `bson:"_id,omitempty"`
	FirstName    string             `bson:"firstName"`
	LastName     string             `bson:"lastName"`
	Email        string             `bson:"email"`
	PhoneNumber  string             `bson:"phoneNumber"`
	PasswordHash string             `bson:"passwordHash"`
	Role         string             `bson:"role"`
	CreatedAt    time.Time          `bson:"createdAt"`
	UpdatedAt    time.Time          `bson:"updatedAt"`

	IsDeleted bool `bson:"isDeleted"`
	IsActive  bool `bson:"isActive"`
}

func FromUser(user model.User) (User, error) {
	objID, err := primitive.ObjectIDFromHex(user.ID)
	if err != nil && user.ID != "" {
		return User{}, ErrInvalidID
	}

	return User{
		ID:           objID,
		FirstName:    user.FirstName,
		LastName:     user.LastName,
		Email:        user.Email,
		PasswordHash: user.PasswordHash,
		Role:         user.Role,
		CreatedAt:    user.CreatedAt,
		UpdatedAt:    user.UpdatedAt,
		IsDeleted:    user.IsDeleted,
		IsActive:     user.IsActive,
	}, nil
}

func ToUser(user User) model.User {
	return model.User{
		ID:           user.ID.Hex(),
		FirstName:    user.FirstName,
		LastName:     user.LastName,
		Email:        user.Email,
		PasswordHash: user.PasswordHash,
		Role:         user.Role,
		CreatedAt:    user.CreatedAt,
		UpdatedAt:    user.UpdatedAt,
		IsDeleted:    user.IsDeleted,
		IsActive:     user.IsActive,
	}
}

func FromUserFilter(filter model.UserFilter) (bson.M, error) {
	query := bson.M{}

	if filter.ID != nil {
		objID, err := primitive.ObjectIDFromHex(*filter.ID)
		if err != nil {
			return query, ErrInvalidID
		}

		query["_id"] = objID
	}

	if filter.FirstName != nil {
		query["firstName"] = *filter.FirstName
	}

	if filter.LastName != nil {
		query["lastName"] = *filter.LastName
	}

	if filter.Email != nil {
		query["email"] = *filter.Email
	}

	if filter.PhoneNumber != nil {
		query["phoneNumber"] = *filter.PhoneNumber
	}

	if filter.PasswordHash != nil {
		query["passwordHash"] = *filter.PasswordHash
	}

	if filter.Role != nil {
		query["role"] = *filter.Role
	}

	if filter.IsDeleted != nil {
		query["isDeleted"] = *filter.IsDeleted
	}

	if filter.IsActive != nil {
		query["isActive"] = *filter.IsActive
	}

	return query, nil
}

func FromUserUpdateData(update model.UserUpdateData) bson.M {
	query := bson.M{}

	if update.FirstName != nil {
		query["firstName"] = *update.FirstName
	}

	if update.LastName != nil {
		query["lastName"] = *update.LastName
	}

	if update.Email != nil {
		query["email"] = *update.Email
	}

	if update.PhoneNumber != nil {
		query["phoneNumber"] = *update.PhoneNumber
	}

	if update.PasswordHash != nil {
		query["passwordHash"] = *update.PasswordHash
	}

	if update.Role != nil {
		query["role"] = *update.Role
	}

	if update.IsDeleted != nil {
		query["isDeleted"] = *update.IsDeleted
	}

	if update.IsActive != nil {
		query["isActive"] = *update.IsActive
	}

	query["updatedAt"] = update.UpdatedAt

	return bson.M{"$set": query}
}
