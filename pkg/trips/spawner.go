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

	store     Store
	sessStore SessionStore
	msgStore  SyncMsgStore
	logger    *zap.Logger
}

func NewSpawner(
	mapsSvc maps.Service,
	mediaSvc media.Service,
	store Store,
	sessStore SessionStore,
	msgStore SyncMsgStore,
	logger *zap.Logger,
) *Spawner {
	return &Spawner{
		crds:      make(map[string]struct{}),
		mapsSvc:   mapsSvc,
		mediaSvc:  mediaSvc,
		store:     store,
		sessStore: sessStore,
		msgStore:  msgStore,
		logger:    logger.Named("trips.spawner"),
	}
}

func (spwn *Spawner) shouldSpawnCoordinator(msg SyncMsgTOB) bool {
	if msg.Topic != SyncMsgTOBTopicJoin {
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
	return !exist
}

// Run listens to Join requests and spawn a coordinator if there is
// only 1 member in the session (i.e new session)
func (spwn *Spawner) Run() error {
	msgCh, done, err := spwn.msgStore.SubTOBReq("*")
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
			spwn.msgStore,
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

		if err = coord.SendFirstMemberJoinMsg(&msg); err != nil {
			// TODO: spwn.ctrlMsgStore.PubRes(msg.TripID, msg)
			continue
		}

		go func() {
			<-doneCh
			spwn.mu.Lock()
			delete(spwn.crds, coord.tripID)
			spwn.mu.Unlock()
		}()
	}
	return nil
}
