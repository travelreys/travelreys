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
	bsonKeyID            = "id"
	CollTripInvites      = "trip_invites"
	CollEmailTripInvites = "email_trip_invites"
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

	ListEmailTripInvites(ctx context.Context, ff ListEmailTripInvitesFilter) (EmailTripInviteList, error)
	ReadEmailTripInvite(ctx context.Context, ID string) (EmailTripInvite, error)
	SaveEmailTripInvite(ctx context.Context, invite EmailTripInvite) error
	DeleteEmailTripInvite(ctx context.Context, ID string) error
}

type store struct {
	db                  *mongo.Database
	tripInviteColl      *mongo.Collection
	emailTripInviteColl *mongo.Collection

	logger *zap.Logger
}

func NewStore(
	ctx context.Context,
	db *mongo.Database,
	logger *zap.Logger,
) Store {
	tripInviteColl := db.Collection(CollTripInvites)
	tripInviteColl.Indexes().CreateMany(ctx, []mongo.IndexModel{
		{Keys: bson.M{bsonKeyID: 1}},
		{Keys: bson.M{"tripID": 1}},
		{Keys: bson.M{"userID": 1}},
	})

	emailTripInviteColl := db.Collection(CollEmailTripInvites)
	emailTripInviteColl.Indexes().CreateMany(ctx, []mongo.IndexModel{
		{Keys: bson.M{bsonKeyID: 1}},
		{Keys: bson.M{"tripID": 1}},
		{Keys: bson.M{"userID": 1}},
		{
			Keys: bson.M{"createdAt": 1},
			Options: options.Index().SetExpireAfterSeconds(
				int32(emailInviteDuration.Seconds()),
			),
		},
	})

	return &store{
		db,
		tripInviteColl,
		emailTripInviteColl,
		logger.Named("invites.store"),
	}
}

func (s *store) ListTripInvites(
	ctx context.Context,
	ff ListTripInvitesFilter,
) (TripInviteList, error) {
	list := TripInviteList{}
	bsonM, _ := ff.toBSON()
	cursor, err := s.tripInviteColl.Find(ctx, bsonM)
	if err != nil {
		s.logger.Error("ListTripInvites", zap.Error(err))
		return list, err
	}
	err = cursor.All(ctx, &list)
	return list, err
}

func (s *store) ReadTripInvite(ctx context.Context, ID string) (TripInvite, error) {
	var invite TripInvite
	err := s.tripInviteColl.FindOne(ctx, bson.M{bsonKeyID: ID}).Decode(&invite)
	if err == mongo.ErrNoDocuments {
		return TripInvite{}, ErrInviteNotFound
	}
	if err != nil {
		s.logger.Error("ReadTripInvite", zap.String("id", ID), zap.Error(err))
		return TripInvite{}, ErrUnexpectedStoreError
	}
	return invite, err
}

func (s *store) SaveTripInvite(ctx context.Context, invite TripInvite) error {
	saveFF := bson.M{bsonKeyID: invite.ID}
	opts := options.Replace().SetUpsert(true)
	_, err := s.tripInviteColl.ReplaceOne(ctx, saveFF, invite, opts)
	if err != nil {
		s.logger.Error("SaveTripInvite", zap.Error(err))
		return ErrUnexpectedStoreError
	}
	return nil
}

func (s *store) DeleteTripInvite(ctx context.Context, ID string) error {
	_, err := s.tripInviteColl.DeleteOne(ctx, bson.M{ID: ID})
	if err != nil {
		s.logger.Error("DeleteTripInvite", zap.String("id", ID), zap.Error(err))
		return ErrUnexpectedStoreError
	}
	return err
}

func (s *store) ListEmailTripInvites(ctx context.Context, ff ListEmailTripInvitesFilter) (EmailTripInviteList, error) {
	list := EmailTripInviteList{}
	bsonM, _ := ff.toBSON()
	cursor, err := s.emailTripInviteColl.Find(ctx, bsonM)
	if err != nil {
		s.logger.Error("ListEmailTripInvites", zap.Error(err))
		return list, err
	}
	err = cursor.All(ctx, &list)
	return list, err
}

func (s *store) ReadEmailTripInvite(ctx context.Context, ID string) (EmailTripInvite, error) {
	var invite EmailTripInvite
	err := s.emailTripInviteColl.FindOne(
		ctx,
		bson.M{bsonKeyID: ID},
	).Decode(&invite)
	if err == mongo.ErrNoDocuments {
		return EmailTripInvite{}, ErrInviteNotFound
	}
	if err != nil {
		s.logger.Error("ReadEmailTripInvite", zap.String("id", ID), zap.Error(err))
		return EmailTripInvite{}, ErrUnexpectedStoreError
	}
	return invite, err
}

func (s *store) SaveEmailTripInvite(ctx context.Context, invite EmailTripInvite) error {
	saveFF := bson.M{bsonKeyID: invite.ID}
	opts := options.Replace().SetUpsert(true)
	_, err := s.emailTripInviteColl.ReplaceOne(ctx, saveFF, invite, opts)
	if err != nil {
		s.logger.Error("SaveEmailTripInvite", zap.Error(err))
		return ErrUnexpectedStoreError
	}
	return nil
}

func (s *store) DeleteEmailTripInvite(ctx context.Context, ID string) error {
	_, err := s.emailTripInviteColl.DeleteOne(ctx, bson.M{ID: ID})
	if err != nil {
		s.logger.Error("DeleteEmailTripInvite", zap.String("id", ID), zap.Error(err))
		return ErrUnexpectedStoreError
	}
	return err
}
