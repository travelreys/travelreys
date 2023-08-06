package trips

import (
	"sync"

	"github.com/travelreys/travelreys/pkg/maps"
	"github.com/travelreys/travelreys/pkg/media"
	"go.uber.org/zap"
)

type Spawner struct {
	crds map[string]struct{} // map of coordinators by tripIDs
	mu   sync.Mutex

	mapsSvc  maps.Service
	mediaSvc media.Service

	store        Store
	sessStore    SessionStore
	ctrlMsgStore SyncMsgControlStore
	dataMsgStore SyncMsgDataStore

	logger *zap.Logger
}

func NewSpawner(
	mapsSvc maps.Service,
	mediaSvc media.Service,
	store Store,
	sessStore SessionStore,
	ctrlMsgStore SyncMsgControlStore,
	dataMsgStore SyncMsgDataStore,
	logger *zap.Logger,
) *Spawner {
	return &Spawner{
		crds:         make(map[string]struct{}),
		mapsSvc:      mapsSvc,
		mediaSvc:     mediaSvc,
		store:        store,
		sessStore:    sessStore,
		ctrlMsgStore: ctrlMsgStore,
		dataMsgStore: dataMsgStore,
		logger:       logger.Named("trips.spawner"),
	}
}

func (spwn *Spawner) shouldSpawnCoordinator(msg SyncMsgControl) bool {
	if msg.Topic != SyncMsgControlTopicJoin {
		return false
	}

	exist := false
	spwn.mu.Lock()
	_, ok := spwn.crds[msg.TripID]
	exist = ok
	if !ok {
		spwn.crds[msg.TripID] = struct{}{}
	}
	spwn.mu.Unlock()
	return exist
}

// Run listens to Join requests and spawn a coordinator if there is
// only 1 member in the session (i.e new session)
func (spwn *Spawner) Run() error {
	msgCh, done, err := spwn.ctrlMsgStore.SubReq("*")
	if err != nil {
		return err
	}

	defer func() {
		done <- true
	}()

	for msg := range msgCh {
		if !spwn.shouldSpawnCoordinator(msg) {
			continue
		}
		coord := NewCoordinator(
			msg.TripID,
			spwn.mapsSvc,
			spwn.mediaSvc,
			spwn.store,
			spwn.sessStore,
			spwn.ctrlMsgStore,
			spwn.dataMsgStore,
			spwn.logger,
		)
		doneCh, err := coord.Init()
		if err != nil {
			spwn.logger.Error("unable to init coordinator", zap.Error(err))
			continue
		}
		if err := coord.Run(); err != nil {
			spwn.logger.Error("unable to run coordinator", zap.Error(err))
			continue
		}
		// coord.SendFirstMemberJoinMsg(ctx, msg, sess)

		go func() {
			<-doneCh
			spwn.mu.Lock()
			delete(spwn.crds, coord.tripID)
			spwn.mu.Unlock()
		}()
	}
	return nil
}
