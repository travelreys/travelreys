package tripssync

import (
	context "context"
	"encoding/json"

	jsonpatch "github.com/evanphx/json-patch/v5"
	"github.com/go-redis/redis/v9"
	"github.com/google/uuid"
	"github.com/nats-io/nats.go"
	"github.com/tiinyplanet/tiinyplanet/pkg/common"
	"github.com/tiinyplanet/tiinyplanet/pkg/reqctx"
	"github.com/tiinyplanet/tiinyplanet/pkg/trips"
)

/***************
/* Coordinator *
/***************/

type Coordinator struct {
	store CoordinatorStore

	id string

	planID string

	// plan is the coordinators' local copy of the trip plan
	plan []byte

	// counter is a monotonically increasing integer
	// for maintaining total order broadcast. All clients
	// should apply operations in sequence of the counter.
	counter uint64

	// queue maintains a FIFO Total Order Broadcast together
	// with Counter.
	queue chan CollabOpMessage
}

func NewCoordinator(planID string, store CoordinatorStore) *Coordinator {
	return &Coordinator{
		store:   store,
		id:      uuid.New().String(),
		planID:  planID,
		counter: 0,
		queue:   make(chan CollabOpMessage, common.DefaultChSize),
	}
}

func (crd *Coordinator) Run(ctx context.Context) error {
	// 1. Initialise plan for coordinator
	plan, err := crd.store.tripStore.ReadTripPlan(reqctx.Context{}, crd.planID)
	if err != nil {
		return err
	}
	planBytes, _ := json.Marshal(plan)
	crd.plan = planBytes

	// 2.Subscribe to updates
	msgCh, err := crd.store.SubcribeSessUpdates(ctx, crd.planID)
	if err != nil {
		return err
	}

	// 3.1. Takes in op msg indicating changes from clients
	// 3.2. Give each message a counter
	// 3.3. Sends each ordered message back to the FIFO queue
	go func() {
		for msg := range msgCh {
			if msg.OpType == CollabOpLeaveSession {
				// close channel when all users have left the session
				sess, _ := crd.store.ReadCollabSession(context.Background(), crd.planID)
				if len(sess.Members) == 0 {
					crd.store.done <- true
					return
				}
			}
			msg.Counter = crd.counter
			crd.queue <- msg
		}
	}()

	// 4.1 Broadcasts the operation to all other connected clients
	// 4.2 Update local plan and persist the data
	go func() {
		for msg := range crd.queue {
			go func(coMsg CollabOpMessage) {
				// Update local copy of plan + validate if the op is valid
				patch, _ := jsonpatch.DecodePatch(coMsg.UpdateTripReq.Bytes())
				modified, err := patch.Apply(crd.plan)
				if err != nil {
					return
				}
				crd.plan = modified
				var plan trips.TripPlan
				json.Unmarshal(crd.plan, &plan)

				// TODO: these 2 ops must be atomic!
				crd.store.tripStore.SaveTripPlan(reqctx.Context{}, plan)
				crd.store.IncrSessionCounter(context.Background(), crd.planID)

				crd.store.PublishTOBUpdates(context.Background(), crd.planID, coMsg)
			}(msg)
		}
	}()

	return nil
}

/*********************
/* Coordinator Store *
/*********************/

type CoordinatorStore struct {
	tripStore trips.Store

	nc  *nats.Conn
	rdb redis.UniversalClient

	done chan bool
}

func NewCoordinatorStore(tripStore trips.Store, nc *nats.Conn, rdb redis.UniversalClient) *CoordinatorStore {
	doneCh := make(chan bool)
	return &CoordinatorStore{tripStore, nc, rdb, doneCh}
}

// Subscription

func (s *CoordinatorStore) SubcribeSessUpdates(ctx context.Context, planID string) (<-chan CollabOpMessage, error) {
	subj := collabSessUpdatesKey(planID)
	natsCh := make(chan *nats.Msg, common.DefaultChSize)
	msgCh := make(chan CollabOpMessage, common.DefaultChSize)

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
				err := json.Unmarshal(natsMsg.Data, &msg)
				if err == nil {
					msgCh <- msg
				}
			}
		}
	}()
	return msgCh, nil
}

// Publish

func (s *CoordinatorStore) PublishTOBUpdates(ctx context.Context, planID string, msg CollabOpMessage) error {
	subj := collabSessTOBKey(planID)
	data, _ := json.Marshal(msg)
	return s.nc.Publish(subj, data)
}

// Counter

func (s *CoordinatorStore) GetSessionCounter(ctx context.Context, planID string) (uint64, error) {
	key := collabSessCounterKey(planID)
	cmd := s.rdb.Get(ctx, key)
	if cmd.Err() != nil {
		return 0, cmd.Err()
	}
	ctr, err := cmd.Int64()
	if err != nil {
		return 0, err
	}
	return uint64(ctr), err
}

func (s *CoordinatorStore) IncrSessionCounter(ctx context.Context, planID string) error {
	key := collabSessCounterKey(planID)
	cmd := s.rdb.Incr(ctx, key)
	return cmd.Err()
}

func (s *CoordinatorStore) ReadCollabSession(ctx context.Context, planID string) (CollabSession, error) {
	var members trips.TripMembersList
	key := collabSessMembersKey(planID)
	data := s.rdb.SMembers(ctx, key)
	err := data.ScanSlice(&members)
	return CollabSession{members}, err
}
