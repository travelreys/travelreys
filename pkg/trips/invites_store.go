package trips

import (
	"context"
	"errors"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
)

var (
	ErrInviteNotFound = errors.New("trips.ErrInviteNotFound")
)

type InviteStore interface {
	List(ctx context.Context, ff ListInvitesFilter) (InviteList, error)
	Read(ctx context.Context, ID string) (Invite, error)
	Save(ctx context.Context, invite Invite) error
	Delete(ctx context.Context, ID string) error
}

type inviteStore struct {
	db   *mongo.Database
	coll *mongo.Collection

	logger *zap.Logger
}

func NewInviteStore(
	ctx context.Context,
	db *mongo.Database,
	logger *zap.Logger,
) InviteStore {
	coll := db.Collection("trip_invites")
	coll.Indexes().CreateMany(ctx, []mongo.IndexModel{
		{Keys: bson.M{bsonKeyID: 1}},
		{Keys: bson.M{"tripID": 1}},
		{Keys: bson.M{"userID": 1}},
	})
	return &inviteStore{db, coll, logger.Named("trips.inviteStore")}
}

func (s *inviteStore) List(
	ctx context.Context,
	ff ListInvitesFilter,
) (InviteList, error) {
	list := InviteList{}
	bsonM, _ := ff.toBSON()
	cursor, err := s.coll.Find(ctx, bsonM)
	if err != nil {
		s.logger.Error("ListInvitesForTrip", zap.Error(err))
		return list, err
	}
	err = cursor.All(ctx, &list)
	return list, err
}

func (s *inviteStore) Read(ctx context.Context, ID string) (Invite, error) {
	var invite Invite
	err := s.coll.FindOne(ctx, bson.M{bsonKeyID: ID}).Decode(&invite)
	if err == mongo.ErrNoDocuments {
		return Invite{}, ErrInviteNotFound
	}
	if err != nil {
		s.logger.Error("Read", zap.String("id", ID), zap.Error(err))
		return Invite{}, ErrUnexpectedStoreError
	}
	return invite, err
}

func (s *inviteStore) Save(ctx context.Context, invite Invite) error {
	saveFF := bson.M{bsonKeyID: invite.ID}
	opts := options.Replace().SetUpsert(true)
	_, err := s.coll.ReplaceOne(ctx, saveFF, invite, opts)
	if err != nil {
		s.logger.Error("Save", zap.Error(err))
		return ErrUnexpectedStoreError
	}
	return nil
}

func (s *inviteStore) Delete(ctx context.Context, ID string) error {
	_, err := s.coll.DeleteOne(ctx, bson.M{ID: ID})
	if err != nil {
		s.logger.Error("Delete", zap.String("id", ID), zap.Error(err))
		return ErrUnexpectedStoreError
	}
	return err
}
