package auth

import (
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	CollectionUsers = "users"
	BsonKeyID       = "id"
	BsonKeyEmail    = "email"
)

var (
	ErrUserNotFound         = errors.New("not-found")
	ErrUnexpectedStoreError = errors.New("store-error")
)

type UpdateUserFilter struct {
	Labels map[string]string `json:"labels" bson:"labels"`
}

type Store interface {
	ReadUserByID(context.Context, string) (User, error)
	ReadUserByEmail(ctx context.Context, Email string) (User, error)
	UpdateUser(context.Context, string, UpdateUserFilter) error
	SaveUser(context.Context, User) error
}

type store struct {
	db       *mongo.Database
	usrsColl *mongo.Collection
}

func NewStore(db *mongo.Database) Store {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	usrsColl := db.Collection(CollectionUsers)

	idIdx := mongo.IndexModel{Keys: bson.M{BsonKeyID: 1}}
	usrsColl.Indexes().CreateOne(ctx, idIdx)

	return &store{db, usrsColl}
}

func (str store) ReadUserByID(ctx context.Context, ID string) (User, error) {
	return str.readUserByFilter(ctx, bson.M{BsonKeyID: ID})
}

func (str store) ReadUserByEmail(ctx context.Context, Email string) (User, error) {
	return str.readUserByFilter(ctx, bson.M{BsonKeyEmail: Email})
}

func (str store) readUserByFilter(ctx context.Context, ff bson.M) (User, error) {
	var usr User

	err := str.usrsColl.FindOne(ctx, ff).Decode(&usr)
	if err == mongo.ErrNoDocuments {
		return usr, ErrUserNotFound
	}
	if err != nil {
		return usr, ErrUnexpectedStoreError
	}
	return usr, err
}

func (store store) UpdateUser(ctx context.Context, ID string, ff UpdateUserFilter) error {
	update := bson.M{"$set": ff}
	_, err := store.usrsColl.UpdateOne(ctx, bson.M{BsonKeyID: ID}, update)
	if err == mongo.ErrNoDocuments {
		return ErrUserNotFound
	}
	return err
}

func (store store) SaveUser(ctx context.Context, usr User) error {
	saveFF := bson.M{BsonKeyID: usr.ID}
	opts := options.Replace().SetUpsert(true)
	_, err := store.usrsColl.ReplaceOne(ctx, saveFF, usr, opts)
	return err
}
