package moodboard

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
	bsonKeyID       = "id"
	moodboardsColl  = "moodboards"
	storeLoggerName = "moodboard.store"
)

var (
	ErrMoodboardNotFound    = errors.New("moodboard.store.not-found")
	ErrUnexpectedStoreError = errors.New("moodboard.store.unexpected-error")
)

type UpdateFilter struct {
	Notes string `json:"notes" bson:"notes"`
}

type Store interface {
	Save(ctx context.Context, mb Moodboard) error
	Read(ctx context.Context, ID string) (Moodboard, error)
	Update(ctx context.Context, ID string, ff UpdateFilter) error
	Delete(ctx context.Context, ID string) error
}

type store struct {
	db     *mongo.Database
	mbColl *mongo.Collection

	logger *zap.Logger
}

func NewStore(ctx context.Context, db *mongo.Database, logger *zap.Logger) Store {
	ctx, cancel := context.WithTimeout(ctx, common.DbReqTimeout)
	defer cancel()

	mbColl := db.Collection(moodboardsColl)
	mbColl.Indexes().CreateMany(ctx, []mongo.IndexModel{
		{Keys: bson.M{bsonKeyID: 1}},
	})
	return &store{db, mbColl, logger.Named(storeLoggerName)}
}

func (s *store) Save(ctx context.Context, mb Moodboard) error {
	saveFF := bson.M{bsonKeyID: mb.ID}
	opts := options.Replace().SetUpsert(true)
	_, err := s.mbColl.ReplaceOne(ctx, saveFF, mb, opts)
	if err != nil {
		s.logger.Error("Save", zap.Error(err))
		return ErrUnexpectedStoreError
	}
	return nil
}

func (s *store) Read(ctx context.Context, ID string) (Moodboard, error) {
	var mb Moodboard
	err := s.mbColl.FindOne(ctx, bson.M{bsonKeyID: ID}).Decode(&mb)
	if err == mongo.ErrNoDocuments {
		return mb, ErrMoodboardNotFound
	}
	if err != nil {
		s.logger.Error("Read", zap.String("id", ID), zap.Error(err))
		return mb, ErrUnexpectedStoreError
	}
	return mb, err
}

func (s *store) Update(ctx context.Context, ID string, ff UpdateFilter) error {
	_, err := s.mbColl.UpdateOne(ctx, bson.M{bsonKeyID: ID}, ff)
	if err == mongo.ErrNoDocuments {
		return ErrMoodboardNotFound
	}
	if err != nil {
		s.logger.Error("Read", zap.String("id", ID), zap.Error(err))
		return ErrUnexpectedStoreError
	}
	return err
}

func (s *store) Delete(ctx context.Context, ID string) error {
	_, err := s.mbColl.DeleteOne(ctx, bson.M{ID: ID})
	if err != nil {
		s.logger.Error("Delete", zap.String("id", ID), zap.Error(err))
		return ErrUnexpectedStoreError
	}
	return err
}
