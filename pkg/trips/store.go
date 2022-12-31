package trips

import (
	context "context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/awhdesmond/tiinyplanet/pkg/reqctx"
	"github.com/go-redis/redis/v9"
	"github.com/nats-io/nats.go"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type ListTripPlansFilter struct {
	ID *string `json:"string"`
}

// Trips Store

type TripStore interface {
	SaveTripPlan(ctx reqctx.Context, plan TripPlan) error
	ReadTripPlan(ctx reqctx.Context, ID string) (TripPlan, error)
	ListTripPlans(ctx reqctx.Context, ff ListTripPlansFilter) (TripPlansList, error)
	DeleteTripPlan(ctx reqctx.Context, ID string) error
}

const (
	BsonKeyID = "id"
)

var (
	ErrPlanNotFound         = errors.New("not-found")
	ErrUnexpectedStoreError = errors.New("store-error")
)

type store struct {
	db *mongo.Database

	tripsCollection *mongo.Collection
}

func NewStore(db *mongo.Database) TripStore {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	tripsCollection := db.Collection("trips")

	idIdx := mongo.IndexModel{Keys: bson.M{BsonKeyID: 1}}
	tripsCollection.Indexes().CreateOne(ctx, idIdx)

	return &store{db, tripsCollection}
}

func (str *store) SaveTripPlan(ctx reqctx.Context, plan TripPlan) error {
	saveFF := bson.M{BsonKeyID: plan.ID}
	opts := options.Replace().SetUpsert(true)
	_, err := str.tripsCollection.ReplaceOne(ctx, saveFF, plan, opts)
	return err
}

func (str *store) ReadTripPlan(ctx reqctx.Context, ID string) (TripPlan, error) {
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

func (str *store) ListTripPlans(ctx reqctx.Context, ff ListTripPlansFilter) (TripPlansList, error) {
	var list TripPlansList
	cursor, err := str.tripsCollection.Find(ctx, ff)
	if err != nil {
		return list, err
	}
	err = cursor.All(ctx, &list)
	return list, err
}

func (str *store) DeleteTripPlan(ctx reqctx.Context, ID string) error {
	ff := bson.M{ID: ID}
	_, err := str.tripsCollection.DeleteOne(ctx, ff)
	return err
}

// Collaboration Store

type CollabStore interface {
	ReadCollabSession(ctx context.Context, planID string) (CollabSession, error)
	AddMemberToCollabSession(ctx context.Context, planID string, member TripMember) error
	RemoveMemberFromCollabSession(ctx context.Context, planID string, member TripMember) error

	SubscribeCollabOpMessages(ctx context.Context, planID string) (chan<- CollabOpMessage, error)
	PublishCollabOpMessages(ctx context.Context, planID string, msg CollabOpMessage) error
}

type collabStore struct {
	nc  *nats.Conn
	rdb redis.UniversalClient

	done chan bool
}

func NewCollabStore(nc *nats.Conn, rdb redis.UniversalClient) CollabStore {
	doneCh := make(chan bool)
	return &collabStore{nc, rdb, doneCh}
}

func (s *collabStore) collabSessionKey(planID string) string {
	return fmt.Sprintf("collab-session:%s", planID)
}

func (s *collabStore) ReadCollabSession(ctx context.Context, planID string) (CollabSession, error) {
	var members []TripMember
	key := s.collabSessionKey(planID)
	data := s.rdb.SMembers(ctx, key)
	err := data.ScanSlice(&members)
	return CollabSession{members}, err
}

func (s *collabStore) AddMemberToCollabSession(ctx context.Context, planID string, member TripMember) error {
	key := s.collabSessionKey(planID)
	value, _ := json.Marshal(member)
	cmd := s.rdb.HSet(ctx, key, member.MemberID, string(value))
	return cmd.Err()
}

func (s *collabStore) RemoveMemberFromCollabSession(ctx context.Context, planID string, member TripMember) error {
	key := s.collabSessionKey(planID)
	cmd := s.rdb.HDel(ctx, key, member.MemberID)
	return cmd.Err()
}

func (s *collabStore) SubscribeCollabOpMessages(ctx context.Context, planID string) (chan<- CollabOpMessage, error) {
	subj := s.collabSessionKey(planID)
	natsCh := make(chan *nats.Msg, 512)
	msgCh := make(chan CollabOpMessage, 512)

	sub, err := s.nc.ChanSubscribe(subj, natsCh)
	if err != nil {
		return nil, err
	}

	go func() {
		for {
			select {
			case <-s.done:
				sub.Unsubscribe()
				close(msgCh)
				return
			case natsMsg := <-natsCh:
				var msg CollabOpMessage
				err := json.Unmarshal(natsMsg.Data, msg)
				if err == nil {
					msgCh <- msg
				}
			}
		}
	}()
	return msgCh, nil
}

func (s *collabStore) PublishCollabOpMessages(ctx context.Context, planID string, msg CollabOpMessage) error {
	subj := s.collabSessionKey(planID)
	data, _ := json.Marshal(msg)
	return s.nc.Publish(subj, data)
}
