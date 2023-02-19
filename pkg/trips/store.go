package trips

import (
	context "context"
	"errors"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
)

const (
	BsonKeyID       = "id"
	StoreLoggerName = "trips.store"
	TripsColl       = "trips"
)

var (
	ErrPlanNotFound         = errors.New("trip.store.trip.notfound")
	ErrUnexpectedStoreError = errors.New("trip.store.unexpectederror")
)

type ListTripsFilter struct {
	ID *string `json:"string"`
}

func (ff ListTripsFilter) toBSON() bson.M {
	bsonM := bson.M{}
	if ff.ID != nil {
		bsonM[BsonKeyID] = ff.ID
	}
	return bsonM
}

type Store interface {
	SaveTripPlan(ctx context.Context, plan TripPlan) error
	ReadTrip(ctx context.Context, ID string) (TripPlan, error)
	ListTrips(ctx context.Context, ff ListTripsFilter) (TripPlansList, error)
	DeleteTrip(ctx context.Context, ID string) error
}

type store struct {
	db        *mongo.Database
	tripsColl *mongo.Collection
	logger    *zap.Logger
}

func NewStore(db *mongo.Database, logger *zap.Logger) Store {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	tripsColl := db.Collection(TripsColl)

	idIdx := mongo.IndexModel{Keys: bson.M{BsonKeyID: 1}}
	tripsColl.Indexes().CreateOne(ctx, idIdx)

	return &store{db, tripsColl, logger.Named(StoreLoggerName)}
}

func (store *store) SaveTripPlan(ctx context.Context, plan TripPlan) error {
	saveFF := bson.M{BsonKeyID: plan.ID}
	opts := options.Replace().SetUpsert(true)
	_, err := store.tripsColl.ReplaceOne(ctx, saveFF, plan, opts)
	if err != nil {
		store.logger.Error(
			"SaveTripPlan",
			zap.String("plan", fmt.Sprintf("%+v", plan)),
			zap.Error(err),
		)
		return ErrUnexpectedStoreError
	}
	return nil
}

func (store *store) ReadTrip(ctx context.Context, ID string) (TripPlan, error) {
	var plan TripPlan
	err := store.tripsColl.FindOne(ctx, bson.M{BsonKeyID: ID}).Decode(&plan)
	if err == mongo.ErrNoDocuments {
		return plan, ErrPlanNotFound
	}
	if err != nil {
		store.logger.Error("ReadTrip", zap.String("id", ID), zap.Error(err))
		return plan, ErrUnexpectedStoreError
	}
	return plan, err
}

func (store *store) ListTrips(ctx context.Context, ff ListTripsFilter) (TripPlansList, error) {
	list := TripPlansList{}
	bsonM := ff.toBSON()
	cursor, err := store.tripsColl.Find(ctx, bsonM)
	if err != nil {
		store.logger.Error("ListTrips", zap.String("ff", fmt.Sprintf("%+v", ff)), zap.Error(err))
		return list, err
	}
	err = cursor.All(ctx, &list)
	return list, err
}

func (store *store) DeleteTrip(ctx context.Context, ID string) error {
	ff := bson.M{ID: ID}
	_, err := store.tripsColl.DeleteOne(ctx, ff)
	if err != nil {
		store.logger.Error("DeleteTrip", zap.String("id", ID), zap.Error(err))
		return ErrUnexpectedStoreError
	}
	return err
}
