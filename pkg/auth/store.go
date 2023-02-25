package auth

import (
	"context"
	"errors"

	"github.com/tiinyplanet/tiinyplanet/pkg/common"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
)

const (
	bsonKeyID       = "id"
	bsonKeyEmail    = "email"
	collectionUsers = "users"
	storeLoggerName = "auth.store"
)

var (
	ErrUserNotFound         = errors.New("auth.store.user-not-found")
	ErrUnexpectedStoreError = errors.New("auth.store.unexpected-error")
)

type ReadFilter struct {
	ID    string `json:"id" bson:"id,omitempty"`
	Email string `json:"email" bson:"email,omitempty"`
}

type UpdateFilter struct {
	Labels common.Labels `json:"labels" bson:"labels"`
}

type ListFilter struct {
	IDs   []string `json:"ids" bson:"id"`
	Email *string  `json:"email" bson:"email"`
}

func (ff ListFilter) toBson() bson.M {
	bsonM := bson.M{}
	if len(ff.IDs) > 0 {
		bsonM[bsonKeyID] = bson.M{"$in": ff.IDs}
	}
	if ff.Email != nil {
		bsonM[bsonKeyEmail] = ff.Email
	}
	return bsonM
}

type Store interface {
	Read(context.Context, ReadFilter) (User, error)
	List(context.Context, ListFilter) (UsersList, error)
	Update(context.Context, string, UpdateFilter) error
	Save(context.Context, User) error
}

type store struct {
	db       *mongo.Database
	usrsColl *mongo.Collection

	logger *zap.Logger
}

func NewStore(ctx context.Context, db *mongo.Database, logger *zap.Logger) Store {
	ctx, cancel := context.WithTimeout(ctx, common.DbReqTimeout)
	defer cancel()

	usrsColl := db.Collection(collectionUsers)

	idIdx := mongo.IndexModel{Keys: bson.M{bsonKeyID: 1}}
	usrsColl.Indexes().CreateOne(ctx, idIdx)

	return &store{db, usrsColl, logger.Named(storeLoggerName)}
}

func (s store) Read(ctx context.Context, ff ReadFilter) (User, error) {
	var usr User

	err := s.usrsColl.FindOne(ctx, ff).Decode(&usr)
	if err == mongo.ErrNoDocuments {
		return usr, ErrUserNotFound
	}
	if err != nil {
		s.logger.Error("Read", zap.String("ff", common.FmtString(ff)), zap.Error(err))
		return usr, ErrUnexpectedStoreError
	}
	return usr, err
}

func (s store) List(ctx context.Context, ff ListFilter) (UsersList, error) {
	list := UsersList{}
	cursor, err := s.usrsColl.Find(ctx, ff.toBson())
	if err != nil {
		s.logger.Error("List", zap.String("ff", common.FmtString(ff)), zap.Error(err))
		return list, err
	}
	err = cursor.All(ctx, &list)
	return list, err
}

func (s store) Update(ctx context.Context, ID string, ff UpdateFilter) error {
	_, err := s.usrsColl.UpdateOne(ctx, bson.M{bsonKeyID: ID}, bson.M{"$set": ff})
	if err == mongo.ErrNoDocuments {
		return ErrUserNotFound
	}
	if err != nil {
		s.logger.Error("Update", zap.String("ID", ID), zap.String("ff", common.FmtString(ff)), zap.Error(err))
	}
	return err
}

func (s store) Save(ctx context.Context, usr User) error {
	saveFF := bson.M{bsonKeyID: usr.ID}
	opts := options.Replace().SetUpsert(true)
	_, err := s.usrsColl.ReplaceOne(ctx, saveFF, usr, opts)
	if err != nil {
		s.logger.Error("Save", zap.String("usr", common.FmtString(usr)), zap.Error(err))
	}
	return err
}
