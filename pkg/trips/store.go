package trips

import (
	context "context"
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
	tripsColl       = "trips"
	storeLoggerName = "trips.store"
)

var (
	ErrTripNotFound         = errors.New("trips.ErrTripNotFound")
	ErrUnexpectedStoreError = errors.New("trips.ErrUnexpectedStoreError")
)

type ListFilter struct {
	UserID     *string
	UserIDs    []string
	OnlyPublic bool
}

func (f ListFilter) toBSON() (bson.M, bool) {
	bsonAnd := bson.A{}
	isSet := false

	if f.UserID != nil {
		bsonUserOr := bson.A{
			bson.M{"creator.id": *f.UserID},
			bson.M{fmt.Sprintf("membersId.%s", *f.UserID): *f.UserID},
		}
		bsonAnd = append(bsonAnd, bson.M{"$or": bsonUserOr})
		isSet = true
	}

	if f.UserIDs != nil && len(f.UserIDs) > 0 {
		bsonUserOr := bson.A{
			bson.M{"creator.id": bson.M{"$in": f.UserIDs}},
		}
		for _, userID := range f.UserIDs {
			bsonUserOr = append(bsonUserOr, bson.M{
				fmt.Sprintf("membersId.%s", userID): userID,
			})
		}
		bsonAnd = append(bsonAnd, bson.M{"$or": bsonUserOr})
		isSet = true
	}

	if f.OnlyPublic {
		sharingAccessLabel := "labels." + LabelSharingAccess
		bsonAnd = append(bsonAnd, bson.M{sharingAccessLabel: SharingAccessViewer})
		isSet = true
	}

	bsonAnd = append(bsonAnd, bson.M{"deleted": false})
	return bson.M{"$and": bsonAnd}, isSet
}

type Store interface {
	Save(ctx context.Context, trip Trip) error
	Read(ctx context.Context, ID string) (Trip, error)
	List(ctx context.Context, ff ListFilter) (TripsList, error)
	Delete(ctx context.Context, ID string) error
}

type store struct {
	db        *mongo.Database
	tripsColl *mongo.Collection
	logger    *zap.Logger
}

func NewStore(ctx context.Context, db *mongo.Database, logger *zap.Logger) Store {
	ctx, cancel := context.WithTimeout(ctx, common.DbReqTimeout)
	defer cancel()

	tripsColl := db.Collection(tripsColl)
	tripsColl.Indexes().CreateMany(ctx, []mongo.IndexModel{
		{Keys: bson.M{bsonKeyID: 1}},
		{Keys: bson.M{"creator.id": 1}},
		{Keys: bson.M{"membersId.$**": 1}},
	})
	return &store{db, tripsColl, logger.Named(storeLoggerName)}
}

func (s *store) Save(ctx context.Context, trip Trip) error {
	saveFF := bson.M{bsonKeyID: trip.ID}
	opts := options.Replace().SetUpsert(true)
	_, err := s.tripsColl.ReplaceOne(ctx, saveFF, trip, opts)
	if err != nil {
		s.logger.Error("Save", zap.Error(err))
		return ErrUnexpectedStoreError
	}
	return nil
}

func (s *store) Read(ctx context.Context, ID string) (Trip, error) {
	var plan Trip
	err := s.tripsColl.FindOne(ctx, bson.M{bsonKeyID: ID}).Decode(&plan)
	if err == mongo.ErrNoDocuments {
		return plan, ErrTripNotFound
	}
	if err != nil {
		s.logger.Error("Read", zap.String("id", ID), zap.Error(err))
		return plan, ErrUnexpectedStoreError
	}
	return plan, err
}

func (s *store) List(ctx context.Context, ff ListFilter) (TripsList, error) {
	list := TripsList{}
	bsonM, ok := ff.toBSON()
	if !ok {
		return list, nil
	}

	opts := options.Find().SetSort(bson.M{"startDate": -1})
	cursor, err := s.tripsColl.Find(ctx, bsonM, opts)
	if err != nil {
		s.logger.Error("List", zap.String("ff", common.FmtString(ff)), zap.Error(err))
		return list, err
	}
	err = cursor.All(ctx, &list)
	return list, err
}

func (s *store) Delete(ctx context.Context, ID string) error {
	_, err := s.tripsColl.DeleteOne(ctx, bson.M{ID: ID})
	if err != nil {
		s.logger.Error("Delete", zap.String("id", ID), zap.Error(err))
		return ErrUnexpectedStoreError
	}
	return err
}
