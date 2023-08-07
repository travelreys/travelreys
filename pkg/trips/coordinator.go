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
	counter          uint64
	counterLeaseID   int64
	refreshCtrTicker *time.Ticker

	// queue maintains a FIFO Total Order Broadcast together
	// with Counter.
	dataFifoMsgQueue chan SyncMsgData

	// msgCh recevies request message from clients
	ctrlMsgCh     <-chan SyncMsgControl
	ctrlMsgDoneCh chan<- bool

	dataMsgCh     <-chan SyncMsgData
	dataMsgDoneCh chan<- bool

	mapsSvc      maps.Service
	mediaSvc     media.Service
	store        Store
	sessStore    SessionStore
	ctrlMsgStore SyncMsgControlStore
	dataMsgStore SyncMsgDataStore

	doneCh chan bool
	logger *zap.Logger
}

func NewCoordinator(
	tripID string,
	mapsSvc maps.Service,
	mediaSvc media.Service,
	store Store,
	sessStore SessionStore,
	ctrlMsgStore SyncMsgControlStore,
	dataMsgStore SyncMsgDataStore,
	logger *zap.Logger,
) *Coordinator {
	return &Coordinator{
		ID:               uuid.New().String(),
		tripID:           tripID,
		trip:             []byte{},
		counter:          1,
		dataFifoMsgQueue: make(chan SyncMsgData, common.DefaultChSize),
		mapsSvc:          mapsSvc,
		mediaSvc:         mediaSvc,
		store:            store,
		ctrlMsgStore:     ctrlMsgStore,
		dataMsgStore:     dataMsgStore,
		sessStore:        sessStore,
		refreshCtrTicker: time.NewTicker(defaultRefreshCounterTTL),
		doneCh:           make(chan bool),
		logger:           logger.Named("trips.coordinator"),
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
	ctrlMsgCh, ctrlMsgDoneCh, err := crd.ctrlMsgStore.SubReq(crd.tripID)
	if err != nil {
		crd.logger.Error("subscribe fail", zap.Error(err))
		return nil, err
	}

	dataMsgCh, dataMsgDoneCh, err := crd.dataMsgStore.SubReqQueue(
		crd.tripID,
		GroupCoordinators,
	)
	if err != nil {
		crd.logger.Error("subscribe fail", zap.Error(err))
		return nil, err
	}

	crd.ctrlMsgCh = ctrlMsgCh
	crd.ctrlMsgDoneCh = ctrlMsgDoneCh
	crd.dataMsgCh = dataMsgCh
	crd.dataMsgDoneCh = dataMsgDoneCh

	// 3. See if there is stale state in etcd
	ctr, err := crd.sessStore.GetCounter(ctx, crd.tripID)
	if err != nil && err != redis.Nil {
		return nil, err
	}
	if ctr != 0 {
		crd.counter = ctr
	}
	return crd.doneCh, nil
}

func (crd *Coordinator) stopControl() {
	if crd.dataMsgDoneCh != nil {
		crd.dataMsgDoneCh <- true
	}
}

func (crd *Coordinator) stopData() {
	if crd.ctrlMsgDoneCh != nil {
		crd.ctrlMsgDoneCh <- true
	}
	close(crd.dataFifoMsgQueue)
}

func (crd *Coordinator) Stop() {
	crd.sessStore.DeleteCounter(context.Background(), crd.tripID, crd.counterLeaseID)
	crd.refreshCtrTicker.Stop()
	crd.stopControl()
	crd.stopData()
	crd.doneCh <- true
}

func (crd *Coordinator) runControl() error {
	crd.logger.Info("running control loop", zap.String("tripID", crd.tripID))

	go func() {
		for msg := range crd.ctrlMsgCh {
			ctx := context.Background()
			crd.logger.Debug("recv control msg", zap.String("topic", msg.Topic))

			switch msg.Topic {
			case SyncMsgControlTopicJoin:
				if err := crd.handleSyncMsgControlJoin(ctx, &msg); err != nil {
					crd.logger.Error("handleSyncMsgControlJoin", zap.Error(err))
				}
			case SyncMsgControlTopicLeave:
				toEnd, err := crd.handleSyncMsgControlLeave(ctx, &msg)
				if err != nil {
					crd.logger.Error("handleSyncMsgControlLeave", zap.Error(err))
					continue
				}
				if toEnd {
					crd.logger.Info("stopping coordinator", zap.String("tripID", crd.tripID))
					crd.Stop()
					return
				}
			}
		}
	}()

	return nil
}

func (crd *Coordinator) runData() error {
	go func() {
		// Takes in msg indicating changes from clients
		for msg := range crd.dataMsgCh {
			crd.logger.Debug("recv data msg", zap.String("topic", msg.Topic))

			// Give each message a counter
			msg.Counter = crd.counter
			crd.counter++
			crd.logger.Debug("next counter", zap.Uint64("counter", crd.counter))

			// 3.3. Sends each ordered message back to the FIFO queue
			crd.dataFifoMsgQueue <- msg
		}
	}()

	go func() {
		// 4.1 Read message from FIFO Queue
		for msg := range crd.dataFifoMsgQueue {
			ctx := context.Background()

			// 4.2 Update local trip and persist the data (if required)
			// Update local copy of trip + validate if the op is valid
			crd.applyDataFifoMsg(ctx, &msg)

			// 4.3 Broadcasts the tob msg to all other connected clients
			crd.dataMsgStore.PubRes(crd.tripID, msg)
		}
	}()

	return nil
}

func (crd *Coordinator) Run() error {
	crd.logger.Info("running coordinator", zap.String("tripID", crd.tripID))

	crd.runControl()
	crd.runData()

	go func() {
		for range crd.refreshCtrTicker.C {
			crd.sessStore.RefreshCounterTTL(context.Background(), crd.tripID, crd.counterLeaseID)
		}
	}()

	return nil
}

func (crd *Coordinator) handleSyncMsgControlJoin(ctx context.Context, msg *SyncMsgControl) error {
	ctxs, err := crd.sessStore.ReadTripSessCtx(ctx, crd.tripID)
	if err != nil {
		return err
	}
	msg.Join.Members = ctxs.ToMembers()
	json.Unmarshal(crd.trip, &msg.Join.Trip)

	for key := range msg.Join.Trip.MediaItems {
		urls, _ := crd.mediaSvc.GenerateGetSignedURLs(ctx, msg.Join.Trip.MediaItems[key])
		for i := 0; i < len(msg.Join.Trip.MediaItems[key]); i++ {
			msg.Join.Trip.MediaItems[key][i].URLs = urls[i]
		}
	}
	return nil
}

func (crd *Coordinator) handleSyncMsgControlLeave(ctx context.Context, msg *SyncMsgControl) (bool, error) {
	// stop coordinator if no users left in session
	ctxs, err := crd.sessStore.ReadTripSessCtx(ctx, crd.tripID)
	if err != nil {
		return false, err
	}
	return len(ctxs) == 0, nil
}

// applyDataFifoMsg handles data messages on crd.dataFifoMsgQueue
// by applying the json.Op before performing additional processing
// based on message topic.
func (crd *Coordinator) applyDataFifoMsg(ctx context.Context, msg *SyncMsgData) {

	patchOps, _ := json.Marshal(msg.Ops)
	patch, _ := jsonpatch.DecodePatch(patchOps)
	modified, err := patch.Apply(crd.trip)
	if err != nil {
		crd.logger.Error("json patch apply", zap.Error(err))
		return
	}

	crd.trip = modified

	var toSave Trip
	if err = json.Unmarshal(crd.trip, &toSave); err != nil {
		crd.logger.Error("json unmarshall fails", zap.Error(err))
	}

	switch msg.Topic {
	case SyncMsgDataTopicAddLodging,
		SyncMsgDataTopicUpdateLodging,
		SyncMsgDataTopicDeleteLodging:
		crd.processLodgingChanged(ctx, &toSave, msg)
		crd.trip, _ = json.Marshal(toSave)
	case SyncMsgDataTopicUpdateTripDates:
		crd.processDatesChanged(&toSave, msg)
		crd.trip, _ = json.Marshal(toSave)
	case SyncMsgDataTopicReorderItinerary,
		SyncMsgDataTopicUpdateActivityPlace,
		SyncMsgDataTopicDeleteActivity:
		crd.processActivityChangedSameDay(ctx, &toSave, msg)
		crd.trip, _ = json.Marshal(toSave)
	case SyncMsgDataTopicReorderActivityToAnotherDay:
		crd.processActivityChangedAnotherDay(ctx, &toSave, msg)
		crd.trip, _ = json.Marshal(toSave)
	case SyncMsgDataTopicOptimizeItinerary:
		crd.processOptimizeRoute(ctx, &toSave, msg)
		crd.trip, _ = json.Marshal(toSave)
	case SyncMsgDataTopicAddMediaItem:
		crd.processAugmentMediaItemSignedURL(ctx, &toSave, msg)
	}

	// Persist trip state to database
	if err := crd.store.Save(ctx, &toSave); err != nil {
		crd.logger.Error("save fails", zap.Error(err))
	}

	leaseID, _ := crd.sessStore.IncrCounter(ctx, crd.tripID, crd.counterLeaseID)
	crd.counterLeaseID = leaseID

}

func (crd *Coordinator) processLodgingChanged(
	ctx context.Context,
	toSave *Trip,
	msg *SyncMsgData,
) {
	mutex := sync.Mutex{}
	wg := sync.WaitGroup{}
	wg.Add(len(toSave.Itineraries))
	for dtKey := range toSave.Itineraries {
		_dtKey := dtKey
		go func(dtKey string) {
			itin, ok := toSave.Itineraries[dtKey]
			if !ok {
				return
			}
			routesMap := crd.calculateRoute(ctx, itin, toSave)
			mutex.Lock()
			crd.UpdateRoutes(ctx, dtKey, routesMap, msg, toSave)
			mutex.Unlock()
			wg.Done()
		}(_dtKey)
	}
	wg.Wait()
}

func (crd *Coordinator) processDatesChanged(
	toSave *Trip,
	msg *SyncMsgData,
) {
	sortedCurrDates := GetSortedItineraryKeys(toSave)
	newItineraries := ItineraryMap{}
	numDays := toSave.EndDate.Sub(toSave.StartDate).Hours() / 24

	for i := 0; i <= int(numDays); i++ {
		dt := toSave.StartDate.Add(time.Duration(i*24) * time.Hour)
		key := dt.Format(ItineraryDtKeyFormat)
		var itin *Itinerary
		if i < len(sortedCurrDates) {
			itin = toSave.Itineraries[sortedCurrDates[i]]
			itin.Date = dt
		} else {
			itin = NewItinerary(dt)
		}
		newItineraries[key] = itin

	}
	toSave.Itineraries = newItineraries
	msg.Ops = append(msg.Ops, jp.MakeRepOp(JSONPathItineraryRoot, newItineraries))

}

func (crd *Coordinator) processActivityChangedSameDay(
	ctx context.Context,
	toSave *Trip,
	msg *SyncMsgData,
) {
	dtKey := crd.parseItinDtKeyFromOps(msg.Ops)
	if dtKey == "" {
		return
	}
	if itin, ok := toSave.Itineraries[dtKey]; ok {
		routesMap := crd.calculateRoute(ctx, itin, toSave)
		crd.UpdateRoutes(ctx, dtKey, routesMap, msg, toSave)
	}
}

func (crd *Coordinator) processActivityChangedAnotherDay(
	ctx context.Context,
	toSave *Trip,
	msg *SyncMsgData,
) {
	origDtKey := crd.parseItinDtKeyFromOps(msg.Ops)
	if origDtKey == "" {
		return
	}
	if itin, ok := toSave.Itineraries[origDtKey]; ok {
		routesMap := crd.calculateRoute(ctx, itin, toSave)
		crd.UpdateRoutes(ctx, origDtKey, routesMap, msg, toSave)
	}
	lessOneOps := msg.Ops[1:]
	destDtKey := crd.parseItinDtKeyFromOps(lessOneOps)
	if destDtKey == "" {
		return
	}
	if itin, ok := toSave.Itineraries[destDtKey]; ok {
		routesMap := crd.calculateRoute(ctx, itin, toSave)
		crd.UpdateRoutes(ctx, destDtKey, routesMap, msg, toSave)
	}
}

func (crd *Coordinator) processAugmentMediaItemSignedURL(
	ctx context.Context,
	toSave *Trip,
	msg *SyncMsgData,
) {
	// /mediaItems/${mediaItemsKey}/-
	tkns := strings.Split(msg.Ops[0].Path, "/")
	if len(tkns) < 4 {
		crd.logger.Error("invalid number of tokens", zap.Int("count", len(tkns)))
		return
	}

	mediaItemKey := tkns[2]
	bytes, _ := json.Marshal(msg.Ops[0].Value)

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

	msg.Ops = append(msg.Ops, jp.MakeAddOp(
		fmt.Sprintf(
			"/mediaItems/%s/%d/urls",
			mediaItemKey,
			len(toSave.MediaItems[mediaItemKey])-1,
		), urls[0],
	))
}

func (crd *Coordinator) processOptimizeRoute(
	ctx context.Context,
	toSave *Trip,
	msg *SyncMsgData,
) {
	// /itineraries/2023-03-26
	op := msg.Ops[0]
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

		itin.Activities[actId].Labels[LabelFractionalIndex] = newFIdx

		jop := jp.MakeRepOp(
			fmt.Sprintf("/itineraries/%s/activities/%s/labels/fIndex", dtKey, actId),
			newFIdx,
		)
		msg.Ops = append(msg.Ops, jop)
	}

	routeMaps := crd.calculateRoute(ctx, itin, toSave)
	crd.UpdateRoutes(ctx, dtKey, routeMaps, msg, toSave)
}

// Update Trip Helpers

// parseItinDtKeyFromOps gets the itinerary dt from the ops array
// e.g /itineraries/2023-03-26/activities/9935afee-8bfd-4148-8be8-79fdb2f12b8e
func (crd *Coordinator) parseItinDtKeyFromOps(ops []jp.Op) string {
	for _, op := range ops {
		if !strings.HasPrefix(op.Path, JSONPathItineraryRoot) {
			continue
		}
		tkns := strings.Split(op.Path, "/")
		if len(tkns) < 3 {
			continue
		}
		return tkns[2]
	}
	return ""
}

func (crd *Coordinator) calculateRoute(
	ctx context.Context,
	itin *Itinerary,
	toSave *Trip,
) maps.RouteListMap {
	result := maps.RouteListMap{}

	lodgings := toSave.Lodgings.GetLodgingsForDate(itin.GetDate())
	pairings := itin.RoutePairings(lodgings)

	for pair := range pairings {
		if _, ok := itin.Routes[pair]; ok {
			result[pair] = itin.Routes[pair]
			continue
		}
		actIds := strings.Split(pair, LabelDelimeter)

		// Origin could be lodging
		orig, ok := itin.Activities[actIds[0]]
		if !ok {
			lod := lodgings[actIds[0]]
			orig = &Activity{ID: lod.ID, Place: lod.Place}
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
	msg *SyncMsgData,
	toSave *Trip,
) {
	for pair, routes := range routesMap {
		toSave.Itineraries[dtKey].Routes[pair] = routes
		jop := jp.MakeAddOp(fmt.Sprintf("/itineraries/%s/routes/%s", dtKey, pair), routes)
		msg.Ops = append(
			msg.Ops, jop,
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
		msg.Ops = append(msg.Ops, jop)
		delete(toSave.Itineraries[dtKey].Routes, pair)
	}

}

// TODO: we need this still
// SendFirstMemberJoinMsg sends a memberUpdate message to the very first member
// func (crd *Coordinator) SendFirstMemberJoinMsg(ctx context.Context, msg Message, sess Session) {
// 	msg.TripID = crd.tripID
// 	crd.AugmentJoinMsgWithTrip(ctx, &msg)
// 	msg.Data.JoinSession.Members = sess.Members

// 	msg.Counter = crd.counter
// 	crd.counter++
// 	crd.queue <- msg
// }
