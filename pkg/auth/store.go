package auth

import (
	"context"
	"errors"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
)

const (
	CollectionUsers = "users"
	BsonKeyID       = "id"
	BsonKeyEmail    = "email"
	StoreLoggerName = "auth.store"
)

var (
	ErrUserNotFound         = errors.New("auth.store.user.notfound")
	ErrUnexpectedStoreError = errors.New("auth.store.unexpectederror")
)

type UpdateUserFilter struct {
	Labels map[string]string `json:"labels" bson:"labels"`
}

type ListUsersFilter struct {
	IDs   []string `json:"ids" bson:"id"`
	Email *string  `json:"email" bson:"email"`
}

func (ff ListUsersFilter) toBson() bson.M {
	bsonM := bson.M{}
	if len(ff.IDs) > 0 {
		bsonM[BsonKeyID] = bson.M{"$in": ff.IDs}
	}
	if ff.Email != nil {
		bsonM[BsonKeyEmail] = ff.Email
	}
	return bsonM
}

type Store interface {
	ReadUserByID(context.Context, string) (User, error)
	ReadUserByEmail(context.Context, string) (User, error)
	ListUsers(context.Context, ListUsersFilter) (UsersList, error)
	UpdateUser(context.Context, string, UpdateUserFilter) error
	SaveUser(context.Context, User) error
}

type store struct {
	db       *mongo.Database
	usrsColl *mongo.Collection
	logger   *zap.Logger
}

func NewStore(db *mongo.Database, logger *zap.Logger) Store {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	usrsColl := db.Collection(CollectionUsers)

	idIdx := mongo.IndexModel{Keys: bson.M{BsonKeyID: 1}}
	usrsColl.Indexes().CreateOne(ctx, idIdx)

	return &store{db, usrsColl, logger.Named(StoreLoggerName)}
}

func (str store) ReadUserByID(ctx context.Context, ID string) (User, error) {
	return str.readUserByFilter(ctx, bson.M{BsonKeyID: ID})
}

func (str store) ReadUserByEmail(ctx context.Context, email string) (User, error) {
	return str.readUserByFilter(ctx, bson.M{BsonKeyEmail: email})
}

func (str store) readUserByFilter(ctx context.Context, ff bson.M) (User, error) {
	var usr User

	err := str.usrsColl.FindOne(ctx, ff).Decode(&usr)
	if err == mongo.ErrNoDocuments {
		return usr, ErrUserNotFound
	}
	if err != nil {
		str.logger.Error(
			"readUserByFilter",
			zap.String("ff", fmt.Sprintf("%+v", ff)),
			zap.Error(err))
		return usr, ErrUnexpectedStoreError
	}
	return usr, err
}

func (str store) ListUsers(ctx context.Context, ff ListUsersFilter) (UsersList, error) {
	list := UsersList{}
	cursor, err := str.usrsColl.Find(ctx, ff.toBson())
	if err != nil {
		str.logger.Error("ListUsers", zap.String("ff", fmt.Sprintf("%+v", ff)), zap.Error(err))
		return list, err
	}
	err = cursor.All(ctx, &list)
	return list, err
}

func (str store) UpdateUser(ctx context.Context, ID string, ff UpdateUserFilter) error {
	update := bson.M{"$set": ff}
	_, err := str.usrsColl.UpdateOne(ctx, bson.M{BsonKeyID: ID}, update)
	if err == mongo.ErrNoDocuments {
		return ErrUserNotFound
	}
	if err != nil {
		str.logger.Error(
			"UpdateUser",
			zap.String("ff", fmt.Sprintf("%+v", ff)),
			zap.String("ID", ID),
			zap.Error(err))
	}
	return err
}

func (str store) SaveUser(ctx context.Context, usr User) error {
	saveFF := bson.M{BsonKeyID: usr.ID}
	opts := options.Replace().SetUpsert(true)
	_, err := str.usrsColl.ReplaceOne(ctx, saveFF, usr, opts)
	if err != nil {
		str.logger.Error(
			"SaveUser",
			zap.String("usr", fmt.Sprintf("%+v", usr)),
			zap.Error(err))
	}
	return err
}
