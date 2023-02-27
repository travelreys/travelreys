package tripssync

import (
	context "context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	jsonpatch "github.com/evanphx/json-patch/v5"
	"github.com/google/uuid"
	"github.com/tiinyplanet/tiinyplanet/pkg/common"
	"github.com/tiinyplanet/tiinyplanet/pkg/maps"
	"github.com/tiinyplanet/tiinyplanet/pkg/trips"
	"go.uber.org/zap"
)

/***************
/* Coordinator *
/***************/

type Coordinator struct {
	ID     string
	tripID string

	// trip is the coordinators' local copy of the trip trip
	trip []byte

	// counter is a monotonically increasing integer
	// for maintaining total order broadcast. All clients
	// should apply operations in sequence of the counter.
	counter uint64

	// queue maintains a FIFO Total Order Broadcast together
	// with Counter.
	queue chan Message

	// msgCh recevies request message from clients
	msgCh  <-chan Message
	doneCh chan<- bool

	mapsSvc   maps.Service
	store     Store
	msgStore  MessageStore
	tobStore  TOBMessageStore
	tripStore trips.Store

	logger *zap.Logger
}

func NewCoordinator(
	tripID string,
	mapsSvc maps.Service,
	store Store,
	msgStore MessageStore,
	tobStore TOBMessageStore,
	tripStore trips.Store,

	logger *zap.Logger,
) *Coordinator {
	return &Coordinator{
		ID:        uuid.New().String(),
		tripID:    tripID,
		trip:      []byte{},
		counter:   1,
		queue:     make(chan Message, common.DefaultChSize),
		mapsSvc:   mapsSvc,
		store:     store,
		msgStore:  msgStore,
		tobStore:  tobStore,
		tripStore: tripStore,
		logger:    logger.Named("sync.coordinator"),
	}
}

func (crd *Coordinator) Init() error {
	crd.logger.Info("new coordinator", zap.String("tripID", crd.tripID))

	// 1. Initialise trip for coordinator
	trip, err := crd.tripStore.Read(context.Background(), crd.tripID)
	if err != nil {
		return err
	}
	tripBytes, _ := json.Marshal(trip)
	crd.trip = tripBytes

	// 2. Subscribe to updates from clients
	msgCh, done, err := crd.msgStore.SubscribeQueue(crd.tripID, GroupCoordinators)
	if err != nil {
		crd.logger.Error("subscribe fail", zap.Error(err))
		return err
	}
	crd.msgCh = msgCh
	crd.doneCh = done
	return nil
}

func (crd *Coordinator) Stop() {
	crd.store.DeleteCounter(context.Background(), crd.tripID)
	if crd.doneCh != nil {
		crd.doneCh <- true
	}
	close(crd.queue)
}

func (crd *Coordinator) Run() error {
	crd.logger.Info("running coordinator", zap.String("tripID", crd.tripID))
	go func() {
		// 3.1. Takes in op msg indicating changes from clients
		for msg := range crd.msgCh {
			crd.logger.Debug("recv msg", zap.String("msg", common.FmtString(msg)))
			ctx := context.Background()
			switch msg.Op {
			case OpJoinSession:
				sess, _ := crd.store.Read(ctx, crd.tripID)
				msg.Data.JoinSession.Members = sess.Members
			case OpLeaveSession:
				// stop coordinator if no users left in session
				sess, _ := crd.store.Read(ctx, crd.tripID)
				if len(sess.Members) == 0 {
					crd.logger.Info("stopping coordinator", zap.String("tripID", crd.tripID))
					crd.Stop()
					return
				}
				msg.Data.LeaveSession.Members = sess.Members
			}
			// 3.2. Give each message a counter
			msg.Counter = crd.counter
			crd.counter++
			crd.logger.Debug("next counter", zap.Uint64("counter", crd.counter))

			// 3.3. Sends each ordered message back to the FIFO queue
			crd.queue <- msg
		}
	}()

	go func() {
		// 4.1 Read message from FIFO Queue
		for msg := range crd.queue {
			ctx := context.Background()
			switch msg.Op {
			case OpUpdateTrip:
				// 4.2 Update local trip and persist the data (if required)
				// Update local copy of trip + validate if the op is valid
				crd.HandleOpUpdateTrip(ctx, &msg)
			}
			// 4.3 Broadcasts the tob msg to all other connected clients
			crd.tobStore.Publish(crd.tripID, msg)
		}
	}()
	return nil
}

// HandleFirstMember sends a memberUpdate message to the very first member
func (crd *Coordinator) HandleFirstMember(ctx context.Context, sess Session) {
	msg := NewMsgJoinSession(crd.tripID, sess.Members)
	msg.Counter = crd.counter
	crd.counter++
	crd.queue <- msg
}

// HandleOpUpdateTrip handles OpUpdateTrip messages on crd.queue
func (crd *Coordinator) HandleOpUpdateTrip(ctx context.Context, msg *Message) {
	patchOps, _ := json.Marshal(msg.Data.UpdateTrip.Ops)
	patch, _ := jsonpatch.DecodePatch(patchOps)
	modified, err := patch.Apply(crd.trip)
	if err != nil {
		crd.logger.Error("json patch apply", zap.Error(err))
		return
	}
	crd.trip = modified

	var toSave trips.Trip
	if err = json.Unmarshal(crd.trip, &toSave); err != nil {
		crd.logger.Error("json unmarshall fails", zap.Error(err))
	}

	fmt.Println("msg title", msg.Data.UpdateTrip.Title)

	if msg.Data.UpdateTrip.Title == MsgUpdateTripTitleReorderItinerary {
		// /itinerary/<id>/...

		for _, op := range msg.Data.UpdateTrip.Ops {
			if !strings.HasPrefix(op.Path, "/itinerary/") {
				continue
			}
			pathTokens := strings.Split(op.Path, "/")
			if len(pathTokens) < 3 {
				continue
			}
			idx, _ := strconv.Atoi(pathTokens[2])
			itinList := toSave.Itinerary[idx]
			pairings := itinList.MakeRoutePairings()
			routesToRemove := []string{}
			fmt.Println("here", pathTokens, idx)

			for pair := range pairings {
				if _, ok := itinList.Routes[pair]; ok {
					continue
				}
				actIds := strings.Split(pair, "|")
				orig := itinList.Activities[actIds[0]]
				dest := itinList.Activities[actIds[1]]
				origAct := toSave.Activities[orig.ActivityListID].Activities[orig.ActivityID]
				destAct := toSave.Activities[dest.ActivityListID].Activities[dest.ActivityID]
				if !(origAct.HasPlace() && destAct.HasPlace()) {
					continue
				}
				routes, err := crd.mapsSvc.Directions(ctx, origAct.Place.PlaceID, destAct.Place.PlaceID, "")
				if err != nil {
					fmt.Println(err)
					continue
				}
				fmt.Println(routes)

				toSave.Itinerary[idx].Routes[pair] = routes
				msg.Data.UpdateTrip.Ops = append(
					msg.Data.UpdateTrip.Ops, common.MakeJSONPatchAddOp(
						fmt.Sprintf("/itinerary/%d/routes/%s", idx, pair),
						routes,
					),
				)
			}
			fmt.Println(pairings)

			for pair := range itinList.Routes {
				if _, ok := pairings[pair]; !ok {
					routesToRemove = append(routesToRemove, pair)
					msg.Data.UpdateTrip.Ops = append(
						msg.Data.UpdateTrip.Ops, common.MakeJSONPatchRemoveOp(
							fmt.Sprintf("/itinerary/%d/routes/%s", idx, pair), "",
						),
					)
				}
			}
			for _, pair := range routesToRemove {
				delete(toSave.Itinerary[idx].Routes, pair)
			}

			fmt.Println(msg.Data.UpdateTrip.Ops)

		}

		crd.trip, _ = json.Marshal(toSave)
	}

	// Persist trip state to database
	if err = crd.tripStore.Save(ctx, toSave); err != nil {
		crd.logger.Error("tripStore save fails", zap.Error(err))
	}

	crd.store.IncrCounter(ctx, crd.tripID)
}

/***********
/* Spawner *
/***********/

type CoordinatorSpawner struct {
	mapsSvc   maps.Service
	store     Store
	msgStore  MessageStore
	tobStore  TOBMessageStore
	tripStore trips.Store

	logger *zap.Logger
}

func NewCoordinatorSpawner(
	mapsSvc maps.Service,
	store Store,
	msgStore MessageStore,
	tobStore TOBMessageStore,
	tripStore trips.Store,
	logger *zap.Logger,
) *CoordinatorSpawner {
	return &CoordinatorSpawner{
		mapsSvc:   mapsSvc,
		store:     store,
		msgStore:  msgStore,
		tobStore:  tobStore,
		tripStore: tripStore,
		logger:    logger,
	}
}

// Run listens to OpJoinSession requests and spawn a
// coordinator if there is only 1 member in the session (i.e new session)
func (spwn *CoordinatorSpawner) Run() error {
	msgCh, done, err := spwn.msgStore.Subscribe("*")
	if err != nil {
		return err
	}

	defer func() {
		done <- true
	}()

	for msg := range msgCh {
		if msg.Op != OpJoinSession {
			continue
		}
		ctx := context.Background()
		sess, err := spwn.store.Read(ctx, msg.TripID)
		if err != nil {
			// TODO: handle error here
			spwn.logger.Error("unable to read session", zap.Error(err))
			continue
		}
		if len(sess.Members) > 1 {
			// First member would have joined and coordinator is created.
			continue
		}
		coord := NewCoordinator(
			msg.TripID,
			spwn.mapsSvc,
			spwn.store,
			spwn.msgStore,
			spwn.tobStore,
			spwn.tripStore,
			spwn.logger,
		)
		if err := coord.Init(); err != nil {
			spwn.logger.Error("unable to init coordinator", zap.Error(err))
			continue
		}
		if err := coord.Run(); err != nil {
			spwn.logger.Error("unable to run coordinator", zap.Error(err))
			continue
		}
		coord.HandleFirstMember(ctx, sess)
	}
	return nil
}
