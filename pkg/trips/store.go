package trips

import (
	context "context"

	"github.com/go-redis/redis/v9"
	"github.com/nats-io/nats.go"
	"go.mongodb.org/mongo-driver/mongo"
)

// Trips Store

type Store interface {
	SaveTripPlan(ctx context.Context, plan TripPlan) error
	ListTripPlans(ctx context.Context) ([]TripPlan, error)
	DeleteTripPlan(ctx context.Context, ID string) error
}

type store struct {
	mongoClient *mongo.Client
}

func NewStore(mongoClient *mongo.Client) Store {
	return store{mongoClient}
}

func (str *store) SaveTripPlan(ctx context.Context, plan TripPlan) error {

}
func (str *store) ListTripPlans(ctx context.Context) ([]TripPlan, error) {

}
func (str *store) DeleteTripPlan(ctx context.Context, ID string) error {

}

// Collaboration Store

type CollabStore interface{}

type collabStore struct {
	natsClient *nats.Conn
	rdb        *redis.Client
}

func NewCollabStore(natsClient *nats.Conn, rdb *redis.Client) CollabStore {
	return collabStore{natsClient, rdb}
}
