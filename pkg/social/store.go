package social

import (
	"context"
	"errors"

	"github.com/travelreys/travelreys/pkg/common"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
)

const (
	defaultPageSize = 10

	bsonKeyID            = "id"
	bsonKeyBindingKey    = "binding"
	bsonKeyRevBindingKey = "revbinding"
	bsonKeyInitiatorID   = "initiatorID"
	bsonKeyTargetID      = "targetID"

	friendsColl     = "friends"
	friendsReqColl  = "friend_requests"
	storeLoggerName = "social.store"
)

var (
	ErrFollowingNotFound    = errors.New("social.ErrFollowingNotFound")
	ErrUnexpectedStoreError = errors.New("social.ErrUnexpectedStoreError")
)

type ListFollowRequestsFilter struct {
	InitiatorID *string
	TargetID    *string
}

func (ff ListFollowRequestsFilter) toBSON() bson.M {
	bsonA := bson.A{}
	if ff.InitiatorID != nil && *ff.InitiatorID != "" {
		bsonA = append(bsonA, bson.M{"initiatorID": *ff.InitiatorID})
	}
	if ff.TargetID != nil && *ff.TargetID != "" {
		bsonA = append(bsonA, bson.M{"targetID": *ff.TargetID})
	}
	return bson.M{"$or": bsonA}
}

type Store interface {
	UpsertFollowRequest(ctx context.Context, freq FollowRequest) error
	GetFollowRequestByID(ctx context.Context, id string) (FollowRequest, error)
	DeleteFollowRequest(ctx context.Context, id string) error
	ListFollowRequests(ctx context.Context, ff ListFollowRequestsFilter) (FollowRequestList, error)

	GetFollowing(ctx context.Context, bindingKey string) (Following, error)
	SaveFollowing(ctx context.Context, following Following) error
	ListFollowers(ctx context.Context, userID string) (FollowingsList, error)
	ListFollowing(ctx context.Context, userID string) (FollowingsList, error)
	DeleteFollowing(ctx context.Context, bindingKey string) error
}

type store struct {
	db             *mongo.Database
	friendsColl    *mongo.Collection
	friendReqsColl *mongo.Collection

	logger *zap.Logger
}

func NewStore(ctx context.Context, db *mongo.Database, logger *zap.Logger) Store {
	ctx, cancel := context.WithTimeout(ctx, common.DbReqTimeout)
	defer cancel()

	friendsColl := db.Collection(friendsColl)
	friendsColl.Indexes().CreateMany(ctx, []mongo.IndexModel{
		{
			Keys:    bson.M{bsonKeyBindingKey: 1},
			Options: options.Index().SetUnique(true),
		},
		{
			Keys: bson.M{"initiatorID": 1, "targetID": 1},
		},
	})
	friendReqColl := db.Collection(friendsReqColl)
	friendReqColl.Indexes().CreateMany(ctx, []mongo.IndexModel{
		{Keys: bson.M{bsonKeyBindingKey: 1}},
	})

	return &store{db, friendsColl, friendReqColl, logger}
}

func (store *store) UpsertFollowRequest(ctx context.Context, freq FollowRequest) error {
	opts := options.Replace().SetUpsert(true)
	_, err := store.friendReqsColl.ReplaceOne(
		ctx, bson.M{"binding": freq.BindingKey}, freq, opts,
	)
	if err != nil {
		store.logger.Error("UpsertFollowRequest", zap.String("freq", common.FmtString(freq)), zap.Error(err))
		return ErrUnexpectedStoreError
	}
	return nil
}

func (store *store) GetFollowRequestByID(ctx context.Context, id string) (FollowRequest, error) {
	var req FollowRequest
	res := store.friendReqsColl.FindOne(ctx, bson.M{bsonKeyID: id})
	if err := res.Decode(&req); err != nil {
		store.logger.Error("GetFollowRequestByID", zap.String("id", id), zap.Error(err))
		return req, ErrUnexpectedStoreError
	}
	return req, nil
}

func (store *store) DeleteFollowRequest(ctx context.Context, id string) error {
	_, err := store.friendReqsColl.DeleteOne(ctx, bson.M{bsonKeyID: id})
	if err != nil {
		store.logger.Error("Delete", zap.String("id", id), zap.Error(err))
		return ErrUnexpectedStoreError
	}
	return nil
}

func (store *store) ListFollowRequests(ctx context.Context, ff ListFollowRequestsFilter) (FollowRequestList, error) {
	freqs := FollowRequestList{}

	bsonFF := ff.toBSON()
	cursor, err := store.friendReqsColl.Find(ctx, bsonFF)
	if err != nil {
		store.logger.Error("ListFollowRequests", zap.String("ff", common.FmtString(ff)), zap.Error(err))
		return freqs, ErrUnexpectedStoreError
	}
	if err := cursor.All(ctx, &freqs); err != nil {
		store.logger.Error("ListFollowRequests", zap.String("ff", common.FmtString(ff)), zap.Error(err))
		return freqs, ErrUnexpectedStoreError
	}
	return freqs, nil
}

func (store *store) SaveFollowing(ctx context.Context, friend Following) error {
	ff := bson.M{bsonKeyBindingKey: friend.BindingKey}
	opts := options.Replace().SetUpsert(true)
	_, err := store.friendsColl.ReplaceOne(ctx, ff, friend, opts)
	if err != nil {
		store.logger.Error("SaveFollowing", zap.Error(err))
		return ErrUnexpectedStoreError
	}
	return nil
}

func (store *store) GetFollowing(ctx context.Context, bindingKey string) (Following, error) {
	var friend Following

	res := store.friendsColl.FindOne(ctx, bson.M{bsonKeyBindingKey: bindingKey})
	if res.Err() == mongo.ErrNoDocuments {
		return Following{}, ErrFollowingNotFound
	}
	if res.Err() != nil {
		store.logger.Error("GetFollowing",
			zap.String("bindingKey", bindingKey),
			zap.Error(res.Err()),
		)
		return friend, ErrUnexpectedStoreError
	}
	if err := res.Decode(&friend); err != nil {
		return friend, ErrUnexpectedStoreError
	}
	return friend, nil
}

func (store *store) ListFollowers(ctx context.Context, userID string) (FollowingsList, error) {
	list := FollowingsList{}
	cursor, err := store.friendsColl.Find(ctx, bson.M{"targetID": userID})
	if err != nil {
		store.logger.Error("ListFollowers", zap.String("userID", userID), zap.Error(err))
		return list, ErrUnexpectedStoreError
	}
	if err := cursor.All(ctx, &list); err != nil {
		store.logger.Error("ListFollowers", zap.String("userID", userID), zap.Error(err))
		return list, ErrUnexpectedStoreError
	}
	return list, nil
}

func (store *store) ListFollowing(ctx context.Context, userID string) (FollowingsList, error) {
	list := FollowingsList{}
	cursor, err := store.friendsColl.Find(ctx, bson.M{"initiatorID": userID})
	if err != nil {
		store.logger.Error("ListFollowing", zap.String("userID", userID), zap.Error(err))
		return list, ErrUnexpectedStoreError
	}
	if err := cursor.All(ctx, &list); err != nil {
		store.logger.Error("ListFollowing", zap.String("userID", userID), zap.Error(err))
		return list, ErrUnexpectedStoreError
	}
	return list, nil
}

func (store *store) DeleteFollowing(ctx context.Context, bindingKey string) error {
	ff := bson.M{bsonKeyBindingKey: bindingKey}
	_, err := store.friendsColl.DeleteOne(ctx, ff)
	if err != nil {
		store.logger.Error("DeleteFollowing", zap.Error(err))
		return ErrUnexpectedStoreError
	}
	return nil
}
