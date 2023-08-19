package invites

import (
	"context"
	"errors"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
)

const (
	bsonKeyID = "id"
)

var (
	ErrInviteNotFound       = errors.New("invites.ErrInviteNotFound")
	ErrUnexpectedStoreError = errors.New("invites.ErrUnexpectedStoreError")
)

type Store interface {
	ListTripInvites(ctx context.Context, ff ListTripInvitesFilter) (TripInviteList, error)
	ReadTripInvite(ctx context.Context, ID string) (TripInvite, error)
	SaveTripInvite(ctx context.Context, invite TripInvite) error
	DeleteTripInvite(ctx context.Context, ID string) error

	// ListAppInvites(ctx context.Context, ff ListAppInvites) (AppTripInviteList, error)
	// ReadAppInvite(ctx context.Context, ID string) AppInvite
}

type store struct {
	db   *mongo.Database
	coll *mongo.Collection

	logger *zap.Logger
}

func NewStore(
	ctx context.Context,
	db *mongo.Database,
	logger *zap.Logger,
) Store {
	coll := db.Collection("trip_invites")
	coll.Indexes().CreateMany(ctx, []mongo.IndexModel{
		{Keys: bson.M{bsonKeyID: 1}},
		{Keys: bson.M{"tripID": 1}},
		{Keys: bson.M{"userID": 1}},
	})
	return &store{db, coll, logger.Named("invites.store")}
}

func (s *store) ListTripInvites(
	ctx context.Context,
	ff ListTripInvitesFilter,
) (TripInviteList, error) {
	list := TripInviteList{}
	bsonM, _ := ff.toBSON()
	cursor, err := s.coll.Find(ctx, bsonM)
	if err != nil {
		s.logger.Error("ListInvitesForTrip", zap.Error(err))
		return list, err
	}
	err = cursor.All(ctx, &list)
	return list, err
}

func (s *store) ReadTripInvite(ctx context.Context, ID string) (TripInvite, error) {
	var invite TripInvite
	err := s.coll.FindOne(ctx, bson.M{bsonKeyID: ID}).Decode(&invite)
	if err == mongo.ErrNoDocuments {
		return TripInvite{}, ErrInviteNotFound
	}
	if err != nil {
		s.logger.Error("Read", zap.String("id", ID), zap.Error(err))
		return TripInvite{}, ErrUnexpectedStoreError
	}
	return invite, err
}

func (s *store) SaveTripInvite(ctx context.Context, invite TripInvite) error {
	saveFF := bson.M{bsonKeyID: invite.ID}
	opts := options.Replace().SetUpsert(true)
	_, err := s.coll.ReplaceOne(ctx, saveFF, invite, opts)
	if err != nil {
		s.logger.Error("Save", zap.Error(err))
		return ErrUnexpectedStoreError
	}
	return nil
}

func (s *store) DeleteTripInvite(ctx context.Context, ID string) error {
	_, err := s.coll.DeleteOne(ctx, bson.M{ID: ID})
	if err != nil {
		s.logger.Error("Delete", zap.String("id", ID), zap.Error(err))
		return ErrUnexpectedStoreError
	}
	return err
}
