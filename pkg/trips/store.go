package trips

import (
	context "context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type ListTripPlansFilter struct {
	ID *string `json:"string"`
}

type Store interface {
	SaveTripPlan(ctx context.Context, plan TripPlan) error
	ReadTripPlan(ctx context.Context, ID string) (TripPlan, error)
	ListTripPlans(ctx context.Context, ff ListTripPlansFilter) (TripPlansList, error)
	DeleteTripPlan(ctx context.Context, ID string) error
}

const (
	BsonKeyID = "id"
)

var (
	ErrPlanNotFound         = errors.New("not-found")
	ErrUnexpectedStoreError = errors.New("store-error")
)

type store struct {
	db              *mongo.Database
	tripsCollection *mongo.Collection
}

func NewStore(db *mongo.Database) Store {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	tripsCollection := db.Collection("trips")

	idIdx := mongo.IndexModel{Keys: bson.M{BsonKeyID: 1}}
	tripsCollection.Indexes().CreateOne(ctx, idIdx)

	return &store{db, tripsCollection}
}

func (str *store) SaveTripPlan(ctx context.Context, plan TripPlan) error {
	saveFF := bson.M{BsonKeyID: plan.ID}
	opts := options.Replace().SetUpsert(true)
	_, err := str.tripsCollection.ReplaceOne(ctx, saveFF, plan, opts)
	return err
}

func (str *store) ReadTripPlan(ctx context.Context, ID string) (TripPlan, error) {
	var plan TripPlan

	err := str.tripsCollection.FindOne(ctx, bson.M{ID: ID}).Decode(&plan)
	if err == mongo.ErrNoDocuments {
		return plan, ErrPlanNotFound
	}
	if err != nil {
		return plan, ErrUnexpectedStoreError
	}
	return plan, err
}

func (str *store) ListTripPlans(ctx context.Context, ff ListTripPlansFilter) (TripPlansList, error) {
	var list TripPlansList
	cursor, err := str.tripsCollection.Find(ctx, ff)
	if err != nil {
		return list, err
	}
	err = cursor.All(ctx, &list)
	return list, err
}

func (str *store) DeleteTripPlan(ctx context.Context, ID string) error {
	ff := bson.M{ID: ID}
	_, err := str.tripsCollection.DeleteOne(ctx, ff)
	return err
}
