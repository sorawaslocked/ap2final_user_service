package mongo

import (
	"context"
	"errors"
	"github.com/sorawaslocked/ap2final_user_service/internal/adapter/mongo/dao"
	"github.com/sorawaslocked/ap2final_user_service/internal/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const userCollection = "users"

type User struct {
	col *mongo.Collection
}

func NewUser(conn *mongo.Database) *User {
	return &User{
		col: conn.Collection(userCollection),
	}
}

func (db *User) InsertOne(ctx context.Context, user model.User) (model.User, error) {
	userDao, err := dao.FromUser(user)
	if err != nil {
		return model.User{}, err
	}

	res, err := db.col.InsertOne(ctx, userDao)
	if err != nil {
		return model.User{}, mongoError("InsertOne", err)
	}

	id := res.InsertedID.(primitive.ObjectID).Hex()

	return db.FindOne(ctx, model.UserFilter{ID: &id})
}

func (db *User) FindOne(ctx context.Context, filter model.UserFilter) (model.User, error) {
	var userDao dao.User

	query, err := dao.FromUserFilter(filter)
	if err != nil {
		return model.User{}, err
	}

	err = db.col.FindOne(ctx, query).Decode(&userDao)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return model.User{}, model.ErrNotFound
		}

		return model.User{}, mongoError("FindOne", err)
	}

	return dao.ToUser(userDao), nil
}

func (db *User) Find(ctx context.Context, filter model.UserFilter) ([]model.User, error) {
	var userDaos []dao.User

	query, err := dao.FromUserFilter(filter)
	if err != nil {
		return []model.User{}, err
	}

	cur, err := db.col.Find(ctx, query)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return []model.User{}, model.ErrNotFound
		}

		return []model.User{}, mongoError("Find", err)
	}

	if err = cur.All(ctx, &userDaos); err != nil {
		return []model.User{}, mongoError("Cursor.All", err)
	}

	users := make([]model.User, len(userDaos))

	for i, userDao := range userDaos {
		users[i] = dao.ToUser(userDao)
	}

	return users, nil
}

func (db *User) UpdateOne(ctx context.Context, filter model.UserFilter, update model.UserUpdateData) (model.User, error) {
	query, err := dao.FromUserFilter(filter)
	if err != nil {
		return model.User{}, err
	}

	res, err := db.col.UpdateOne(
		ctx,
		query,
		dao.FromUserUpdateData(update),
	)
	if err != nil {
		return model.User{}, mongoError("UpdateOne", err)
	}

	if res.MatchedCount == 0 {
		return model.User{}, model.ErrNotFound
	}

	return db.FindOne(ctx, filter)
}

func (db *User) DeleteOne(ctx context.Context, filter model.UserFilter) (model.User, error) {
	var userDao dao.User

	query, err := dao.FromUserFilter(filter)
	if err != nil {
		return model.User{}, err
	}

	err = db.col.FindOne(ctx, query).Decode(&userDao)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return model.User{}, model.ErrNotFound
		}

		return model.User{}, mongoError("FindOne", err)
	}

	_, err = db.col.DeleteOne(ctx, query)
	if err != nil {
		return model.User{}, mongoError("DeleteOne", err)
	}

	return dao.ToUser(userDao), nil
}
