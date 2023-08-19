package auth

import (
	"context"
	"errors"
	"fmt"
	"math"
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
)

var (
	ErrInvalidFilter        = errors.New("auth.ErrInvalidFilter")
	ErrUserNotFound         = errors.New("auth.ErrUserNotFound")
	ErrDuplicateUsername    = errors.New("auth.ErrDuplicateUsername")
	ErrUnexpectedStoreError = errors.New("auth.ErrUnexpectedStoreError")
)

type Store interface {
	Read(ctx context.Context, ff ReadFilter) (User, error)
	List(ctx context.Context, ff ListFilter) (UsersList, error)
	Update(ctx context.Context, ID string, ff UpdateFilter) error
	Save(ctx context.Context, usr User) error
	Delete(ctx context.Context, ID string) error

	GetOTP(ctx context.Context, email string) (string, error)
	SaveOTP(ctx context.Context, email string, otp string, dur time.Duration) error
}

type store struct {
	db       *mongo.Database
	usrsColl *mongo.Collection

	rdb    redis.UniversalClient
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
	usrsColl.Indexes().CreateMany(ctx, []mongo.IndexModel{
		{Keys: bson.M{bsonKeyID: 1}},
		{Keys: bson.M{bsonKeyEmail: 1}, Options: options.Index().SetUnique(true)},
		{Keys: bson.M{bsonKeyUsername: 1}, Options: options.Index().SetUnique(true)},
	})
	return &store{db, usrsColl, rdb, logger.Named("auth.store")}
}

func (s store) Read(ctx context.Context, ff ReadFilter) (User, error) {
	var usr User
	err := s.usrsColl.FindOne(ctx, ff).Decode(&usr)
	if err == mongo.ErrNoDocuments {
		return usr, ErrUserNotFound
	}
	if err != nil {
		s.logger.Error(
			"Read",
			zap.String("ff", common.FmtString(ff)),
			zap.Error(err),
		)
		return usr, ErrUnexpectedStoreError
	}
	return usr, err
}

func (s store) List(ctx context.Context, ff ListFilter) (UsersList, error) {
	list := UsersList{}
	bsonM, ok := ff.toBSON()
	if !ok {
		return UsersList{}, nil
	}
	opts := options.Find()
	if ff.PageCount != nil {
		opts.SetLimit(int64(math.Min(float64(*ff.PageCount), 20)))
	}
	cursor, err := s.usrsColl.Find(ctx, bsonM, opts)
	if err != nil {
		s.logger.Error(
			"List",
			zap.String("ff", common.FmtString(ff)),
			zap.Error(err),
		)
		return list, err
	}
	err = cursor.All(ctx, &list)
	return list, err
}

func (s store) Update(ctx context.Context, ID string, ff UpdateFilter) error {
	bsonM, ok := ff.toBSON()
	if !ok {
		return nil
	}
	_, err := s.usrsColl.UpdateOne(
		ctx,
		bson.M{bsonKeyID: ID},
		bson.M{"$set": bsonM},
	)
	if err == mongo.ErrNoDocuments {
		return ErrUserNotFound
	}
	if err != nil {
		s.logger.Error(
			"Update",
			zap.String("ID", ID),
			zap.String("ff", common.FmtString(ff)),
			zap.Error(err),
		)
		if common.MongoIsDupError(err) {
			return ErrDuplicateUsername
		}
	}
	return err
}

func (s store) Delete(ctx context.Context, ID string) error {
	_, err := s.usrsColl.UpdateOne(
		ctx,
		bson.M{bsonKeyID: ID},
		bson.M{"$set": bson.M{
			"email":  "",
			"name":   "",
			"labels": &common.Labels{},
		}},
	)
	return err
}

func (s store) Save(ctx context.Context, usr User) error {
	saveFF := bson.M{bsonKeyID: usr.ID}
	opts := options.Replace().SetUpsert(true)
	_, err := s.usrsColl.ReplaceOne(ctx, saveFF, usr, opts)
	if err != nil {
		s.logger.Error(
			"Save",
			zap.String("usr", common.FmtString(usr)),
			zap.Error(err),
		)
	}
	return err
}

func (s store) GetOTP(ctx context.Context, email string) (string, error) {
	cmd := s.rdb.Get(ctx, fmt.Sprintf("otp:%s", email))
	if cmd.Err() != nil {
		return "", cmd.Err()
	}
	return cmd.Val(), nil
}

func (s store) SaveOTP(ctx context.Context, email, otp string, ttl time.Duration) error {
	cmd := s.rdb.Set(ctx, fmt.Sprintf("otp:%s", email), otp, ttl)
	return cmd.Err()
}
