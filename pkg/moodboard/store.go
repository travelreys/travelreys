package moodboard

import (
	"context"
	"errors"
	"fmt"

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

type UpdateBoardFilter struct {
	Title string `json:"title" bson:"title"`
}

func (ff UpdateBoardFilter) toBsonM() bson.M {
	return bson.M{"$set": bson.M{"title": ff.Title}}
}

type UpdatePinFilter struct {
	Notes string `json:"notes" bson:"notes"`
}

func (ff UpdatePinFilter) toBsonM(pinID string) bson.M {
	return bson.M{"$set": bson.M{fmt.Sprintf("pins.%s.notes", pinID): ff.Notes}}
}

type Store interface {
	Save(ctx context.Context, mb Moodboard) error
	Read(ctx context.Context, ID string) (Moodboard, error)
	Update(ctx context.Context, ID string, ff UpdateBoardFilter) error
	SavePin(ctx context.Context, ID string, pin Pin) error
	UpdatePin(ctx context.Context, ID, pinID string, ff UpdatePinFilter) error
	DeletePin(ctx context.Context, ID, pinID string) error
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

func (s *store) Update(ctx context.Context, ID string, ff UpdateBoardFilter) error {
	filterFF := bson.M{bsonKeyID: ID}
	_, err := s.mbColl.UpdateOne(ctx, filterFF, ff.toBsonM())
	if err == mongo.ErrNoDocuments {
		return ErrMoodboardNotFound
	}
	if err != nil {
		s.logger.Error("Read", zap.String("id", ID), zap.Error(err))
		return ErrUnexpectedStoreError
	}
	return err
}

func (s *store) SavePin(ctx context.Context, ID string, pin Pin) error {
	filterFF := bson.M{bsonKeyID: ID}
	updateFF := bson.M{"$set": bson.M{fmt.Sprintf("pin.%s", pin.ID): pin}}

	opts := options.Replace().SetUpsert(true)
	_, err := s.mbColl.ReplaceOne(ctx, filterFF, updateFF, opts)
	if err != nil {
		s.logger.Error("Save", zap.Error(err))
		return ErrUnexpectedStoreError
	}
	return nil
}

func (s *store) UpdatePin(ctx context.Context, ID, pinID string, ff UpdatePinFilter) error {
	filterFF := bson.M{bsonKeyID: ID}
	_, err := s.mbColl.UpdateOne(ctx, filterFF, ff.toBsonM(pinID))
	if err == mongo.ErrNoDocuments {
		return ErrMoodboardNotFound
	}
	if err != nil {
		s.logger.Error("Read", zap.String("id", ID), zap.Error(err))
		return ErrUnexpectedStoreError
	}
	return err
}

func (s *store) DeletePin(ctx context.Context, ID, pinID string) error {
	filterFF := bson.M{bsonKeyID: ID}
	updateFF := bson.M{"$unset": bson.M{fmt.Sprintf("pin.%s", pinID): ""}}
	_, err := s.mbColl.UpdateOne(ctx, filterFF, updateFF)
	if err != nil {
		s.logger.Error("Delete", zap.String("id", ID), zap.Error(err))
		return ErrUnexpectedStoreError
	}
	return err
}
