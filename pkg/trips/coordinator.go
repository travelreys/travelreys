package trips

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"sync"
	"time"

	jsonpatch "github.com/evanphx/json-patch/v5"
	"github.com/go-redis/redis/v9"
	"github.com/google/uuid"
	"github.com/travelreys/travelreys/pkg/common"
	jp "github.com/travelreys/travelreys/pkg/jsonpatch"
	"github.com/travelreys/travelreys/pkg/maps"
	"github.com/travelreys/travelreys/pkg/media"
	"go.uber.org/zap"
)

const (
	defaultRefreshCounterTTL = 1 * time.Minute
)

type Coordinator struct {
	ID     string
	tripID string

	// trip is the coordinators' local copy of the trip
	trip []byte

	// counter is a monotonically increasing integer
	// for maintaining total order broadcast. All clients
	// should apply operations in sequence of the counter.
	counter        uint64
	counterLeaseID int64

	// queue maintains a FIFO Total Order Broadcast together
	// with Counter.
	syncMsgDataQueue chan SyncMsg

	// msgCh recevies request message from clients
	msgCh     <-chan SyncMsg
	msgDoneCh chan<- bool

	refreshCtrTicker    *time.Ticker
	mapsSvc             maps.Service
	mediaSvc            media.Service
	store               Store
	sessStore           SessionStore
	syncMsgControlStore SyncMsgControlStore
	syncMsgDataStore    SyncMsgDataStore

	doneCh chan bool
	logger *zap.Logger
}

func NewCoordinator(
	tripID string,
	mapsSvc maps.Service,
	mediaSvc media.Service,
	store Store,
	sessStore SessionStore,
	syncMsgControlStore SyncMsgControlStore,
	syncMsgDataStore SyncMsgDataStore,
	logger *zap.Logger,
) *Coordinator {
	return &Coordinator{
		ID:                  uuid.New().String(),
		tripID:              tripID,
		trip:                []byte{},
		counter:             1,
		syncMsgDataQueue:    make(chan SyncMsg, common.DefaultChSize),
		refreshCtrTicker:    time.NewTicker(defaultRefreshCounterTTL),
		mapsSvc:             mapsSvc,
		mediaSvc:            mediaSvc,
		store:               store,
		syncMsgControlStore: syncMsgControlStore,
		syncMsgDataStore:    syncMsgDataStore,
		sessStore:           sessStore,
		doneCh:              make(chan bool),
		logger:              logger.Named("sync.coordinator"),
	}
}

func (crd *Coordinator) Init() (<-chan bool, error) {
	crd.logger.Info("new coordinator", zap.String("tripID", crd.tripID))

	// 1. Initialise trip for coordinator
	ctx := context.Background()
	trip, err := crd.store.Read(ctx, crd.tripID)
	if err != nil {
		return nil, err
	}
	tripBytes, _ := json.Marshal(trip)
	crd.trip = tripBytes

	// 2. Subscribe to updates from clients
	ctrlMsgCh, ctrlMsgDoneCh, err := crd.syncMsgControlStore.Subscribe(crd.tripID)
	if err != nil {
		crd.logger.Error("subscribe fail", zap.Error(err))
		return nil, err
	}

	dataMsgCh, dataMsgDoneCh, err := crd.syncMsgDataStore.SubscribeQueue(crd.tripID, GroupCoordinators)
	if err != nil {
		crd.logger.Error("subscribe fail", zap.Error(err))
		return nil, err
	}
	crd.msgCh = msgCh
	crd.msgDoneCh = msgDoneCh

	// 3. See if there is stale state in redis
	ctr, err := crd.sessStore.GetCounter(ctx, crd.tripID)
	if err != nil && err != redis.Nil {
		return nil, err
	}
	if ctr != 0 {
		crd.counter = ctr
	}
	return crd.doneCh, nil
}

func (crd *Coordinator) Stop() {
	crd.sessStore.DeleteCounter(context.Background(), crd.tripID, crd.counterLeaseID)
	crd.refreshCtrTicker.Stop()
	if crd.msgDoneCh != nil {
		crd.msgDoneCh <- true
	}
	close(crd.syncMsgDataQueue)
	crd.doneCh <- true
}

func (crd *Coordinator) Run() error {
	crd.logger.Info("running coordinator", zap.String("tripID", crd.tripID))
	go func() {
		// 3.1. Takes in op msg indicating changes from clients
		for msg := range crd.msgCh {
			ctx := context.Background()
			if msg.Type == SyncMsgTypeControl {
				crd.logger.Debug("recv control msg", zap.String("topic", msg.Control.Topic))

				continue
			}

			if msg.Type == SyncMsgTypeData {
				crd.logger.Debug("recv data msg", zap.String("topic", msg.Data.Topic))

				continue
			}

			continue

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

	go func() {
		for range crd.refreshCtrTicker.C {
			crd.store.RefreshCounterTTL(context.Background(), crd.tripID)
		}
	}()

	return nil
}

// SendFirstMemberJoinMsg sends a memberUpdate message to the very first member
func (crd *Coordinator) SendFirstMemberJoinMsg(ctx context.Context, msg Message, sess Session) {
	msg.TripID = crd.tripID
	crd.AugmentJoinMsgWithTrip(ctx, &msg)
	msg.Data.JoinSession.Members = sess.Members

	msg.Counter = crd.counter
	crd.counter++
	crd.queue <- msg
}

func (crd *Coordinator) HandleMsgOpJoinSession(ctx context.Context, msg *Message) {
	sess, _ := crd.store.Read(ctx, crd.tripID)
	msg.Data.JoinSession.Members = sess.Members
	crd.AugmentJoinMsgWithTrip(ctx, msg)
}

func (crd *Coordinator) AugmentJoinMsgWithTrip(ctx context.Context, msg *Message) {
	var trip trips.Trip
	json.Unmarshal(crd.trip, &trip)

	for key := range trip.MediaItems {
		urls, _ := crd.mediaSvc.GenerateGetSignedURLs(ctx, trip.MediaItems[key])
		for i := 0; i < len(trip.MediaItems[key]); i++ {
			trip.MediaItems[key][i].URLs = urls[i]
		}
	}
	msg.Data.JoinSession.Trip = trip
}

func (crd *Coordinator) HandleMsgOpLeaveSession(ctx context.Context, msg *Message) {
	sess, _ := crd.store.Read(ctx, crd.tripID)
	msg.Data.LeaveSession.Members = sess.Members
}

// HandleTobOpUpdateTrip handles OpUpdateTrip messages on crd.queue
// by applying the json.Op before performing additional processing
// based on message title.
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

	case UpdateTripTitleAddLodging,
		UpdateTripTitleUpdateLodging,
		UpdateTripTitleDeleteLodging:

		mutex := sync.Mutex{}
		wg := sync.WaitGroup{}
		wg.Add(len(toSave.Itineraries))
		for dtKey := range toSave.Itineraries {
			dtKey := dtKey
			go func(dtKey string) {
				itin, ok := toSave.Itineraries[dtKey]
				if !ok {
					return
				}
				routesMap := crd.CalculateRoute(ctx, &itin, &toSave)
				mutex.Lock()
				crd.UpdateRoutes(ctx, dtKey, routesMap, msg, &toSave)
				mutex.Unlock()
				wg.Done()
			}(dtKey)
		}
		wg.Wait()

	case UpdateTripTitleUpdateTripDates:
		crd.ChangeDates(ctx, msg, &toSave)
	case UpdateTripTitleReorderItinerary,
		UpdateTripTitleUpdateActivityPlace,
		UpdateTripTitleDeleteActivity:
		dtKey := crd.FindItineraryDtKey(msg.Data.UpdateTrip.Ops)
		if dtKey == "" {
			break
		}
		if itin, ok := toSave.Itineraries[dtKey]; ok {
			routesMap := crd.CalculateRoute(ctx, &itin, &toSave)
			crd.UpdateRoutes(ctx, dtKey, routesMap, msg, &toSave)
		}
	case UpdateTripTitleReorderActivityToAnotherDay:
		origDtKey := crd.FindItineraryDtKey(msg.Data.UpdateTrip.Ops)
		if origDtKey == "" {
			break
		}
		if itin, ok := toSave.Itineraries[origDtKey]; ok {
			routesMap := crd.CalculateRoute(ctx, &itin, &toSave)
			crd.UpdateRoutes(ctx, origDtKey, routesMap, msg, &toSave)
		}
		lessOneOps := msg.Data.UpdateTrip.Ops[1:]
		destDtKey := crd.FindItineraryDtKey(lessOneOps)
		if destDtKey == "" {
			break
		}
		if itin, ok := toSave.Itineraries[destDtKey]; ok {
			routesMap := crd.CalculateRoute(ctx, &itin, &toSave)
			crd.UpdateRoutes(ctx, destDtKey, routesMap, msg, &toSave)
		}
	case UpdateTripTitleOptimizeItinerary:
		crd.OptimizeRoute(ctx, msg, &toSave)
	case UpdateTripTitleAddMediaItem:
		crd.AugmentMediaItemSignedURL(ctx, msg, &toSave)
	}

	// Persist trip state to database
	if err = crd.tripStore.Save(ctx, toSave); err != nil {
		crd.logger.Error("tripStore save fails", zap.Error(err))
	}

	crd.store.IncrCounter(ctx, crd.tripID)
}

func (crd *Coordinator) FindItineraryDtKey(ops []jp.Op) string {
	// /itineraries/2023-03-26/activities/9935afee-8bfd-4148-8be8-79fdb2f12b8e
	var dtKey string
	for _, op := range ops {
		if !strings.HasPrefix(op.Path, "/itineraries/") {
			continue
		}
		pathTokens := strings.Split(op.Path, "/")
		if len(pathTokens) < 3 {
			continue
		}
		dtKey = pathTokens[2]
		break
	}
	return dtKey
}

func (crd *Coordinator) ChangeDates(ctx context.Context, msg *Message, toSave *trips.Trip) {
	sortedCurrDates := trips.GetSortedItineraryKeys(toSave)
	numCurrDates := len(sortedCurrDates)
	newItineraries := map[string]trips.Itinerary{}

	numDays := toSave.EndDate.Sub(toSave.StartDate).Hours() / 24
	for i := 0; i <= int(numDays); i++ {
		dt := toSave.StartDate.Add(time.Duration(i*24) * time.Hour)
		key := dt.Format("2006-01-02")
		if i < numCurrDates {
			itin := toSave.Itineraries[sortedCurrDates[i]]
			itin.Date = dt
			newItineraries[key] = itin
		} else {
			newItineraries[key] = trips.NewItinerary(dt)
		}

	}
	toSave.Itineraries = newItineraries
	jop := jp.MakeRepOp("/itineraries", newItineraries)
	msg.Data.UpdateTrip.Ops = append(msg.Data.UpdateTrip.Ops, jop)
	crd.trip, _ = json.Marshal(toSave)
}

func (crd *Coordinator) CalculateRoute(ctx context.Context, itin *trips.Itinerary, toSave *trips.Trip) map[string]maps.RouteList {
	result := map[string]maps.RouteList{}

	lodgings := toSave.Lodgings.GetLodgingsForDate(itin.GetDate())
	pairings := itin.MakeRoutePairings(lodgings)

	for pair := range pairings {
		if _, ok := itin.Routes[pair]; ok {
			result[pair] = itin.Routes[pair]
			continue
		}
		actIds := strings.Split(pair, trips.LabelDelimeter)

		// Origin could be lodging
		orig, ok := itin.Activities[actIds[0]]
		if !ok {
			lod := lodgings[actIds[0]]
			orig = trips.Activity{ID: lod.ID, Place: lod.Place}
		}
		dest := itin.Activities[actIds[1]]
		if !(orig.HasPlace() && dest.HasPlace()) {
			continue
		}

		routes, err := crd.mapsSvc.Directions(
			ctx, orig.Place.PlaceID(), dest.Place.PlaceID(), maps.DirectionModesAllList,
		)
		if err != nil {
			continue
		}
		if len(routes) > 0 {
			shortestRoute, _ := routes.GetMostCommonSenseRoute()
			routes = maps.RouteList{shortestRoute}
		}
		result[pair] = routes
	}
	return result
}

func (crd *Coordinator) UpdateRoutes(
	ctx context.Context,
	dtKey string,
	routesMap map[string]maps.RouteList,
	msg *Message,
	toSave *trips.Trip,
) {
	for pair, routes := range routesMap {
		toSave.Itineraries[dtKey].Routes[pair] = routes
		jop := jp.MakeAddOp(fmt.Sprintf("/itineraries/%s/routes/%s", dtKey, pair), routes)
		msg.Data.UpdateTrip.Ops = append(
			msg.Data.UpdateTrip.Ops, jop,
		)
	}

	routesToRemove := []string{}
	for pair := range toSave.Itineraries[dtKey].Routes {
		if _, ok := routesMap[pair]; ok {
			continue
		}
		routesToRemove = append(routesToRemove, pair)
	}

	for _, pair := range routesToRemove {
		jop := jp.MakeRemoveOp(fmt.Sprintf("/itineraries/%s/routes/%s", dtKey, pair), "")
		msg.Data.UpdateTrip.Ops = append(msg.Data.UpdateTrip.Ops, jop)
		delete(toSave.Itineraries[dtKey].Routes, pair)
	}
	crd.trip, _ = json.Marshal(toSave)
}

func (crd *Coordinator) OptimizeRoute(ctx context.Context, msg *Message, toSave *trips.Trip) {
	// /itineraries/2023-03-26
	op := msg.Data.UpdateTrip.Ops[0]
	pathTokens := strings.Split(op.Path, "/")
	if len(pathTokens) < 3 {
		return
	}
	dtKey := pathTokens[2]
	itin := toSave.Itineraries[dtKey]
	sorted := itin.SortActivities()
	sortedFracIndexes := sorted.GetFracIndexes()
	placeIDs := []string{}
	for _, act := range sorted {
		placeIDs = append(placeIDs, act.Place.PlaceID())
	}

	routes, waypointsOrder, err := crd.mapsSvc.OptimizeRoute(
		ctx, placeIDs[0], placeIDs[len(placeIDs)-1], placeIDs[1:len(placeIDs)-1],
	)
	if err != nil || len(routes) <= 0 {
		return
	}

	for moveToIdx, currIdx := range waypointsOrder {
		actId := sorted[currIdx+1].ID
		newFIdx := sortedFracIndexes[moveToIdx+1]

		itin.Activities[actId].Labels[trips.LabelFractionalIndex] = newFIdx

		jop := jp.MakeRepOp(fmt.Sprintf("/itineraries/%s/activities/%s/labels/fIndex", dtKey, actId), newFIdx)
		msg.Data.UpdateTrip.Ops = append(msg.Data.UpdateTrip.Ops, jop)
	}

	routeMaps := crd.CalculateRoute(ctx, &itin, toSave)
	crd.UpdateRoutes(ctx, dtKey, routeMaps, msg, toSave)
}

func (crd *Coordinator) AugmentMediaItemSignedURL(ctx context.Context, msg *Message, toSave *trips.Trip) {
	// /mediaItems/${mediaItemsKey}/-
	tkns := strings.Split(msg.Data.UpdateTrip.Ops[0].Path, "/")
	if len(tkns) < 4 {
		crd.logger.Error("invalid number of tokens", zap.Int("count", len(tkns)))
		return
	}

	mediaItemKey := tkns[2]
	bytes, _ := json.Marshal(msg.Data.UpdateTrip.Ops[0].Value)
	var mediaItem media.MediaItem

	if err := json.Unmarshal(bytes, &mediaItem); err != nil {
		crd.logger.Error("unmarshall media item", zap.Error(err))
		return
	}

	urls, err := crd.mediaSvc.GenerateGetSignedURLs(ctx, media.MediaItemList{mediaItem})
	if err != nil {
		crd.logger.Error("generate signed urls", zap.Error(err))
		return
	}

	msg.Data.UpdateTrip.Ops = append(msg.Data.UpdateTrip.Ops, jp.MakeAddOp(
		fmt.Sprintf(
			"/mediaItems/%s/%d/urls",
			mediaItemKey,
			len(toSave.MediaItems[mediaItemKey])-1,
		), urls[0],
	))
}
