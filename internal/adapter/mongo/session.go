package mongo

import (
	"context"
	"errors"
	"github.com/sorawaslocked/ap2final_user_service/internal/adapter/mongo/dao"
	"github.com/sorawaslocked/ap2final_user_service/internal/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

const collectionSessions = "sessions"

type Session struct {
	col *mongo.Collection
}

func NewSession(conn *mongo.Database) *Session {
	return &Session{
		col: conn.Collection(collectionSessions),
	}
}

func (db *Session) InsertOne(ctx context.Context, session model.Session) error {
	_, err := db.col.InsertOne(ctx, session)

	return err
}

func (db *Session) FindOneByToken(ctx context.Context, token string) (model.Session, error) {
	var session dao.Session

	err := db.col.FindOne(ctx, bson.M{"refreshToken": token}).Decode(&session)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return model.Session{}, model.ErrNotFound
		}

		return model.Session{}, err
	}

	return dao.ToSession(session), nil
}

func (db *Session) DeleteByToken(ctx context.Context, token string) error {
	res, err := db.col.DeleteOne(ctx, bson.M{"refreshToken": token})

	if err == nil && res.DeletedCount == 0 {
		return model.ErrNotFound
	}

	return err
}
