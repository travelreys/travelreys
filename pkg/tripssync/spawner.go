package tripssync

import (
	context "context"

	"github.com/tiinyplanet/tiinyplanet/pkg/maps"
	"github.com/tiinyplanet/tiinyplanet/pkg/trips"
	"go.uber.org/zap"
)

type Spawner struct {
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
		coord.SendFirstMemberJoinMsg(ctx, sess)
	}
	return nil
}