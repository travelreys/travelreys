package tripssync

import (
	context "context"
	"sync"

	"github.com/travelreys/travelreys/pkg/maps"
	"github.com/travelreys/travelreys/pkg/trips"
	"go.uber.org/zap"
)

type Spawner struct {
	crds map[string]struct{} // map of coordinators by tripIDs
	mu   sync.Mutex

	mapsSvc   maps.Service
	store     Store
	msgStore  MessageStore
	tobStore  TOBMessageStore
	tripStore trips.Store

	logger *zap.Logger
}

func NewSpawner(
	mapsSvc maps.Service,
	store Store,
	msgStore MessageStore,
	tobStore TOBMessageStore,
	tripStore trips.Store,
	logger *zap.Logger,
) *Spawner {
	return &Spawner{
		crds:      make(map[string]struct{}),
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
func (spwn *Spawner) Run() error {
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
		// if len(sess.Members) > 1 {
		// 	continue
		// }

		exist := false
		spwn.mu.Lock()
		_, ok := spwn.crds[msg.TripID]
		exist = ok
		if !ok {
			spwn.crds[msg.TripID] = struct{}{}
		}
		spwn.mu.Unlock()
		if exist {
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
		doneCh, err := coord.Init()
		if err != nil {
			spwn.logger.Error("unable to init coordinator", zap.Error(err))
			continue
		}
		if err := coord.Run(); err != nil {
			spwn.logger.Error("unable to run coordinator", zap.Error(err))
			continue
		}
		coord.SendFirstMemberJoinMsg(ctx, sess)

		go func() {
			<-doneCh
			spwn.mu.Lock()
			delete(spwn.crds, coord.tripID)
			spwn.mu.Unlock()
		}()
	}
	return nil
}
