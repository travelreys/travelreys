package auth

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/go-redis/redis/v9"
	"github.com/travelreys/travelreys/pkg/common"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
)

const (
	bsonKeyID          = "id"
	bsonKeyEmail       = "email"
	bsonKeyName        = "name"
	bsonKeyUsername    = "username"
	bsonKeyPhoneNumber = "phonenumber"
	bsonKeyLabels      = "labels"
	collectionUsers    = "users"
	storeLoggerName    = "auth.store"
)

var (
	ErrUserNotFound         = errors.New("auth.store.user-not-found")
	ErrDuplicateUsername    = errors.New("auth.store.duplicateUsername")
	ErrUnexpectedStoreError = errors.New("auth.store.unexpected-error")
)

type ReadFilter struct {
	ID    string `json:"id" bson:"id,omitempty"`
	Email string `json:"email" bson:"email,omitempty"`
}

type UpdateFilter struct {
	Email       *string `json:"email" bson:"email"`
	Name        *string `json:"name" bson:"name"`
	Username    *string `json:"username" bson:"username"`
	PhoneNumber *PhoneNumber
	Labels      *common.Labels `json:"labels" bson:"labels"`
}

func (ff UpdateFilter) toBSON() bson.M {
	bsonM := bson.M{}

	if ff.Email != nil && *ff.Email != "" {
		bsonM[bsonKeyEmail] = ff.Email
	}
	if ff.Name != nil && *ff.Name != "" {
		bsonM[bsonKeyName] = ff.Name
	}
	if ff.Username != nil && *ff.Username != "" {
		bsonM[bsonKeyUsername] = ff.Username
	}
	if ff.Labels != nil {
		bsonM[bsonKeyLabels] = ff.Labels
	}
	return bsonM
}

type ListFilter struct {
	IDs      []string `json:"ids" bson:"id"`
	Username *string  `json:"username" bson:"username"`
	Email    *string  `json:"email" bson:"email"`
}

func (ff ListFilter) toBSON() bson.M {
	bsonM := bson.M{}
	if len(ff.IDs) > 0 {
		bsonM[bsonKeyID] = bson.M{"$in": ff.IDs}
	}
	if ff.Email != nil && *ff.Username != "" {
		bsonM[bsonKeyEmail] = ff.Email
	}
	if ff.Username != nil && *ff.Username != "" {
		bsonM[bsonKeyUsername] = ff.Username
	}
	return bsonM
}

type Store interface {
	Read(context.Context, ReadFilter) (User, error)
	List(context.Context, ListFilter) (UsersList, error)
	Update(context.Context, string, UpdateFilter) error
	Save(context.Context, User) error

	GetOTP(context.Context, string) (string, error)
	SaveOTP(context.Context, string, string, time.Duration) error
}

type store struct {
	db       *mongo.Database
	usrsColl *mongo.Collection

	rdb redis.UniversalClient

	logger *zap.Logger
}

func NewStore(
	ctx context.Context,
	db *mongo.Database,
	rdb redis.UniversalClient,
	logger *zap.Logger,
) Store {
	ctx, cancel := context.WithTimeout(ctx, common.DbReqTimeout)
	defer cancel()

	usrsColl := db.Collection(collectionUsers)
	if _, err := usrsColl.Indexes().CreateMany(ctx, []mongo.IndexModel{
		{Keys: bson.M{bsonKeyID: 1}},
		{Keys: bson.M{bsonKeyEmail: 1}, Options: options.Index().SetUnique(true)},
		{Keys: bson.M{bsonKeyUsername: 1}, Options: options.Index().SetUnique(true)},
	}); err != nil {
		fmt.Println(err)
	}

	return &store{db, usrsColl, rdb, logger.Named(storeLoggerName)}
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
	cursor, err := s.usrsColl.Find(ctx, ff.toBSON())

	if err != nil {
		s.logger.Error("List", zap.String("ff", common.FmtString(ff)), zap.Error(err))
		return list, err
	}
	err = cursor.All(ctx, &list)
	return list, err
}

func (s store) Update(ctx context.Context, ID string, ff UpdateFilter) error {
	_, err := s.usrsColl.UpdateOne(ctx, bson.M{bsonKeyID: ID}, bson.M{"$set": ff.toBSON()})
	if err == mongo.ErrNoDocuments {
		return ErrUserNotFound
	}
	if err != nil {
		s.logger.Error("Update",
			zap.String("ID", ID),
			zap.String("ff", common.FmtString(ff)),
			zap.Error(err),
		)
		if common.MongoIsDupError(err) {
			return ErrDuplicateUsername
		}
		return ErrUnexpectedStoreError
	}

	return nil
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

func (s store) GetOTP(ctx context.Context, ID string) (string, error) {
	cmd := s.rdb.Get(ctx, fmt.Sprintf("otp:%s", ID))
	if cmd.Err() != nil {
		return "", cmd.Err()
	}
	return cmd.Val(), nil
}

func (s store) SaveOTP(ctx context.Context, ID, otp string, ttl time.Duration) error {
	cmd := s.rdb.Set(ctx, fmt.Sprintf("otp:%s", ID), otp, ttl)
	return cmd.Err()
}
