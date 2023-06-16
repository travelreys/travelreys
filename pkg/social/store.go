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
	friendsColl          = "friends"
	friendsReqColl       = "friend_requests"
	storeLoggerName      = "social.store"
)

var (
	ErrFriendNotFound       = errors.New("social.store.friendNotFound")
	ErrUnexpectedStoreError = errors.New("social.store.unexpected-error")
)

type ListFriendRequestsFilter struct {
	InitiatorID *string
	TargetID    *string
}

func (ff ListFriendRequestsFilter) toBSON() bson.M {
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
	UpsertFriendRequest(ctx context.Context, freq FriendRequest) error
	GetFriendRequestByID(ctx context.Context, id string) (FriendRequest, error)
	DeleteFriendRequest(ctx context.Context, id string) error
	ListFriendRequests(ctx context.Context, ff ListFriendRequestsFilter) (FriendRequestList, error)

	GetFriend(ctx context.Context, bindingKey string) (Friend, error)
	SaveFriend(ctx context.Context, friend Friend) error
	ListFollowers(ctx context.Context, userID string) (FriendsList, error)
	ListFollowing(ctx context.Context, userID string) (FriendsList, error)
	DeleteFriend(ctx context.Context, bindingKey string) error
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

func (store *store) UpsertFriendRequest(ctx context.Context, freq FriendRequest) error {
	opts := options.Replace().SetUpsert(true)
	_, err := store.friendReqsColl.ReplaceOne(
		ctx, bson.M{"binding": freq.BindingKey}, freq, opts,
	)
	if err != nil {
		store.logger.Error("UpsertFriendRequest", zap.String("freq", common.FmtString(freq)), zap.Error(err))
		return ErrUnexpectedStoreError
	}
	return nil
}

func (store *store) GetFriendRequestByID(ctx context.Context, id string) (FriendRequest, error) {
	var req FriendRequest
	res := store.friendReqsColl.FindOne(ctx, bson.M{bsonKeyID: id})
	if err := res.Decode(&req); err != nil {
		store.logger.Error("GetFriendRequestByID", zap.String("id", id), zap.Error(err))
		return req, ErrUnexpectedStoreError
	}
	return req, nil
}

func (store *store) DeleteFriendRequest(ctx context.Context, id string) error {
	_, err := store.friendReqsColl.DeleteOne(ctx, bson.M{bsonKeyID: id})
	if err != nil {
		store.logger.Error("Delete", zap.String("id", id), zap.Error(err))
		return ErrUnexpectedStoreError
	}
	return nil
}

func (store *store) ListFriendRequests(ctx context.Context, ff ListFriendRequestsFilter) (FriendRequestList, error) {
	freqs := FriendRequestList{}

	bsonFF := ff.toBSON()
	cursor, err := store.friendReqsColl.Find(ctx, bsonFF)
	if err != nil {
		store.logger.Error("ListFriendRequests", zap.String("ff", common.FmtString(ff)), zap.Error(err))
		return freqs, ErrUnexpectedStoreError
	}
	if err := cursor.All(ctx, &freqs); err != nil {
		store.logger.Error("ListFriendRequests", zap.String("ff", common.FmtString(ff)), zap.Error(err))
		return freqs, ErrUnexpectedStoreError
	}
	return freqs, nil
}

func (store *store) SaveFriend(ctx context.Context, friend Friend) error {
	ff := bson.M{bsonKeyBindingKey: friend.BindingKey}
	opts := options.Replace().SetUpsert(true)
	_, err := store.friendsColl.ReplaceOne(ctx, ff, friend, opts)
	if err != nil {
		store.logger.Error("SaveFriend", zap.Error(err))
		return ErrUnexpectedStoreError
	}
	return nil
}

func (store *store) GetFriend(ctx context.Context, bindingKey string) (Friend, error) {
	var friend Friend

	res := store.friendsColl.FindOne(ctx, bson.M{bsonKeyBindingKey: bindingKey})
	if res.Err() == mongo.ErrNoDocuments {
		return Friend{}, ErrFriendNotFound
	}
	if res.Err() != nil {
		store.logger.Error("GetFriend",
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

func (store *store) ListFollowers(ctx context.Context, userID string) (FriendsList, error) {
	list := FriendsList{}
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

func (store *store) ListFollowing(ctx context.Context, userID string) (FriendsList, error) {
	list := FriendsList{}
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

func (store *store) DeleteFriend(ctx context.Context, bindingKey string) error {
	ff := bson.M{bsonKeyBindingKey: bindingKey}
	_, err := store.friendsColl.DeleteOne(ctx, ff)
	if err != nil {
		store.logger.Error("DeleteFriend", zap.Error(err))
		return ErrUnexpectedStoreError
	}
	return nil
}
