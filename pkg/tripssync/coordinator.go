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
	jp "github.com/tiinyplanet/tiinyplanet/pkg/jsonpatch"
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
				crd.HandleMsgOpJoinSession(ctx, &msg)
			case OpLeaveSession:
				// stop coordinator if no users left in session
				crd.HandleMsgOpLeaveSession(ctx, &msg)
				if len(msg.Data.LeaveSession.Members) == 0 {
					crd.logger.Info("stopping coordinator", zap.String("tripID", crd.tripID))
					crd.Stop()
					return
				}
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
				crd.HandleTobOpUpdateTrip(ctx, &msg)
			}
			// 4.3 Broadcasts the tob msg to all other connected clients
			crd.tobStore.Publish(crd.tripID, msg)
		}
	}()
	return nil
}

// SendFirstMemberJoinMsg sends a memberUpdate message to the very first member
func (crd *Coordinator) SendFirstMemberJoinMsg(ctx context.Context, sess Session) {
	msg := NewMsgJoinSession(crd.tripID, sess.Members)
	msg.Counter = crd.counter
	crd.counter++
	crd.queue <- msg
}

func (crd *Coordinator) HandleMsgOpJoinSession(ctx context.Context, msg *Message) {
	sess, _ := crd.store.Read(ctx, crd.tripID)
	msg.Data.JoinSession.Members = sess.Members
}

func (crd *Coordinator) HandleMsgOpLeaveSession(ctx context.Context, msg *Message) {
	sess, _ := crd.store.Read(ctx, crd.tripID)
	msg.Data.LeaveSession.Members = sess.Members
}

// HandleTobOpUpdateTrip handles OpUpdateTrip messages on crd.queue
func (crd *Coordinator) HandleTobOpUpdateTrip(ctx context.Context, msg *Message) {
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

	switch msg.Data.UpdateTrip.Title {
	case MsgUpdateTripTitleReorderItinerary:
		crd.HandleTobOpUpdateTripReorderItinerary(ctx, msg, &toSave)
	case MsgUpdateTripOptimizeItineraryRoute:
		crd.HandleTobOpUpdateTripOptimizeItineraryRoute(ctx, msg, &toSave)
	}

	// Persist trip state to database
	if err = crd.tripStore.Save(ctx, toSave); err != nil {
		crd.logger.Error("tripStore save fails", zap.Error(err))
	}

	crd.store.IncrCounter(ctx, crd.tripID)
}

func (crd *Coordinator) HandleTobOpUpdateTripReorderItinerary(ctx context.Context, msg *Message, toSave *trips.Trip) {
	for _, op := range msg.Data.UpdateTrip.Ops {
		// /itinerary/<id>/...
		if !strings.HasPrefix(op.Path, "/itinerary/") {
			continue
		}
		pathTokens := strings.Split(op.Path, "/")
		if len(pathTokens) < 3 {
			continue
		}
		idx, _ := strconv.Atoi(pathTokens[2])
		itin := toSave.Itinerary[idx]
		pairings := itin.MakeRoutePairings()
		routesToRemove := []string{}

		for pair := range pairings {
			if _, ok := itin.Routes[pair]; ok {
				continue
			}
			itinActIds := strings.Split(pair, trips.LabelDelimeter)
			orig := itin.Activities[itinActIds[0]]
			dest := itin.Activities[itinActIds[1]]
			origAct := toSave.Activities[orig.ActivityListID].Activities[orig.ActivityID]
			destAct := toSave.Activities[dest.ActivityListID].Activities[dest.ActivityID]
			if !(origAct.HasPlace() && destAct.HasPlace()) {
				continue
			}
			routes, err := crd.mapsSvc.Directions(ctx, origAct.Place.PlaceID, destAct.Place.PlaceID, "")
			if err != nil {
				continue
			}

			toSave.Itinerary[idx].Routes[pair] = routes
			jop := jp.MakeAddOp(fmt.Sprintf("/itinerary/%d/routes/%s", idx, pair), routes)
			msg.Data.UpdateTrip.Ops = append(
				msg.Data.UpdateTrip.Ops, jop,
			)
		}

		for pair := range itin.Routes {
			if _, ok := pairings[pair]; ok {
				continue
			}
			routesToRemove = append(routesToRemove, pair)
		}

		for _, pair := range routesToRemove {
			jop := jp.MakeRemoveOp(fmt.Sprintf("/itinerary/%d/routes/%s", idx, pair), "")
			msg.Data.UpdateTrip.Ops = append(msg.Data.UpdateTrip.Ops, jop)
			delete(toSave.Itinerary[idx].Routes, pair)
		}
	}
	crd.trip, _ = json.Marshal(toSave)
}

func (crd *Coordinator) HandleTobOpUpdateTripOptimizeItineraryRoute(ctx context.Context, msg *Message, toSave *trips.Trip) {
	op := msg.Data.UpdateTrip.Ops[0]
	pathTokens := strings.Split(op.Path, "/")
	if len(pathTokens) < 3 {
		return
	}

	idx, _ := strconv.Atoi(pathTokens[2])
	itin := toSave.Itinerary[idx]
	sorted := itin.SortActivities()
	sortedFracIndexes := trips.GetFracIndexes(sorted)
	placeIDs := []string{}
	for _, itinAct := range sorted {
		act := toSave.Activities[itinAct.ActivityListID].Activities[itinAct.ActivityID]
		placeIDs = append(placeIDs, act.Place.PlaceID)
	}

	routes, err := crd.mapsSvc.OptimizeRoute(
		ctx, placeIDs[0], placeIDs[len(placeIDs)-1], placeIDs[1:len(placeIDs)-1],
	)
	if err != nil || len(routes) <= 0 {
		return
	}

	for moveToIdx, currIdx := range routes[0].WaypointOrder {
		actId := sorted[currIdx+1].ID
		newFIdx := sortedFracIndexes[moveToIdx+1]

		itin.Activities[actId].Labels[trips.LabelFractionalIndex] = newFIdx

		jop := jp.MakeRepOp(fmt.Sprintf("/itinerary/%d/activities/%s/labels/fIndex", idx, actId), newFIdx)
		msg.Data.UpdateTrip.Ops = append(msg.Data.UpdateTrip.Ops, jop)
	}

	crd.HandleTobOpUpdateTripReorderItinerary(ctx, msg, toSave)
}
