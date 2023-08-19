package trips

import (
	context "context"
	"errors"

	"github.com/travelreys/travelreys/pkg/common"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
)

const (
	mongoCollTrips = "trips"

	bsonKeyID        = "id"
	bsonKeyCreatorId = "creator.id"
	bsonKeyDeleted   = "deleted"

	storeLoggerName = "trips.store"
)

var (
	ErrTripNotFound         = errors.New("trips.ErrTripNotFound")
	ErrUnexpectedStoreError = errors.New("trips.ErrUnexpectedStoreError")
)

type Store interface {
	Save(ctx context.Context, trip *Trip) error
	Read(ctx context.Context, ID string) (*Trip, error)
	List(ctx context.Context, ff ListFilter) (TripsList, error)
	Delete(ctx context.Context, ID string) error
}

type store struct {
	db   *mongo.Database
	coll *mongo.Collection

	logger *zap.Logger
}

func NewStore(ctx context.Context, db *mongo.Database, logger *zap.Logger) Store {
	coll := db.Collection(mongoCollTrips)
	coll.Indexes().CreateMany(ctx, []mongo.IndexModel{
		{Keys: bson.M{bsonKeyID: 1}},
		{Keys: bson.M{bsonKeyCreatorId: 1}},
		{Keys: bson.M{"membersId.$**": 1}},
	})
	return &store{db, coll, logger.Named(storeLoggerName)}
}

func (s *store) Save(ctx context.Context, trip *Trip) error {
	saveFF := bson.M{bsonKeyID: trip.ID}
	opts := options.Replace().SetUpsert(true)
	_, err := s.coll.ReplaceOne(ctx, saveFF, trip, opts)
	if err != nil {
		s.logger.Error("Save", zap.Error(err))
		return ErrUnexpectedStoreError
	}
	return nil
}

func (s *store) Read(ctx context.Context, ID string) (*Trip, error) {
	var trip Trip
	err := s.coll.FindOne(ctx, bson.M{bsonKeyID: ID}).Decode(&trip)
	if err == mongo.ErrNoDocuments {
		return nil, ErrTripNotFound
	}
	if err != nil {
		s.logger.Error("Read", zap.String("id", ID), zap.Error(err))
		return nil, ErrUnexpectedStoreError
	}
	return &trip, err
}

func (s *store) List(ctx context.Context, ff ListFilter) (TripsList, error) {
	list := TripsList{}
	bsonM, ok := ff.toBSON()
	if !ok {
		return list, nil
	}

	opts := options.Find().SetSort(bson.M{"startDate": -1})
	cursor, err := s.coll.Find(ctx, bsonM, opts)
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

func (s *store) Delete(ctx context.Context, ID string) error {
	_, err := s.coll.DeleteOne(ctx, bson.M{ID: ID})
	if err != nil {
		s.logger.Error("Delete", zap.String("id", ID), zap.Error(err))
		return ErrUnexpectedStoreError
	}
	return err
}
