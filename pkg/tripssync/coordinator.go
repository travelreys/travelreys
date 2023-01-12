package tripssync

import (
	context "context"
	"encoding/json"
	"fmt"

	jsonpatch "github.com/evanphx/json-patch/v5"
	"github.com/google/uuid"
	"github.com/tiinyplanet/tiinyplanet/pkg/common"
	"github.com/tiinyplanet/tiinyplanet/pkg/reqctx"
	"github.com/tiinyplanet/tiinyplanet/pkg/trips"
	"go.uber.org/zap"
)

/***************
/* Coordinator *
/***************/

type Coordinator struct {
	id     string
	planID string

	// plan is the coordinators' local copy of the trip plan
	plan []byte

	// counter is a monotonically increasing integer
	// for maintaining total order broadcast. All clients
	// should apply operations in sequence of the counter.
	counter uint64

	// queue maintains a FIFO Total Order Broadcast together
	// with Counter.
	queue chan SyncMessage

	// msgCh recevies request message from clients
	msgCh  <-chan SyncMessage
	doneCh chan<- bool

	sesnStore SessionStore
	sms       SyncMessageStore
	tms       TOBMessageStore
	tripStore trips.Store

	logger *zap.Logger
}

func NewCoordinator(
	planID string,
	sesnStore SessionStore,
	sms SyncMessageStore,
	tms TOBMessageStore,
	tripStore trips.Store,
	logger *zap.Logger,
) *Coordinator {
	return &Coordinator{
		id:        uuid.New().String(),
		planID:    planID,
		plan:      []byte{},
		counter:   0,
		queue:     make(chan SyncMessage, common.DefaultChSize),
		sesnStore: sesnStore,
		sms:       sms,
		tms:       tms,
		tripStore: tripStore,
		logger:    logger,
	}
}

func (crd *Coordinator) Init() error {
	crd.logger.Info("new coordinator", zap.String("planID", crd.planID))

	// 1. Initialise plan for coordinator
	plan, err := crd.tripStore.ReadTripPlan(context.Background(), crd.planID)
	if err != nil {
		return err
	}
	planBytes, _ := json.Marshal(plan)
	crd.plan = planBytes

	// 2.Subscribe to updates from clients
	msgCh, done, err := crd.sms.SubscribeQueue(crd.planID, GroupCoordinators)
	if err != nil {
		crd.logger.Error("subscribe fail", zap.Error(err))
		return err
	}
	crd.msgCh = msgCh
	crd.doneCh = done
	return nil
}

func (crd *Coordinator) Stop() {
	if crd.doneCh != nil {
		crd.doneCh <- true
	}
	close(crd.queue)
}

func (crd *Coordinator) Run() error {
	// 3.1. Takes in op msg indicating changes from clients
	// 3.2. Give each message a counter
	// 3.3. Sends each ordered message back to the FIFO queue
	go func() {
		crd.logger.Info("running coordinator", zap.String("planID", crd.planID))
		for msg := range crd.msgCh {
			crd.logger.Debug("recv msg", zap.String("msg", fmt.Sprintf("%+v", msg)))
			if msg.OpType == SyncOpLeaveSession {
				// close coordinator when all users have left the session
				sess, _ := crd.sesnStore.Read(context.Background(), crd.planID)
				if len(sess.Members) == 0 {
					crd.logger.Info("stopping coordinator", zap.String("planID", crd.planID))
					crd.Stop()
					return
				}
				// Replace with leaving broadcast message
				msg = NewSyncMessageLeaveSessionBroadcast(msg.TripPlanID, sess.Members)
			}

			if msg.OpType == SyncOpJoinSession {
				// Got new player, reset counter
				crd.counter = 0

				sess, _ := crd.sesnStore.Read(context.Background(), crd.planID)
				msg = NewSyncMessageJoinSessionBroadcast(msg.TripPlanID, sess.Members)
			}

			msg.Counter = crd.counter
			crd.queue <- msg

			// Need to update the counter
			crd.counter++
			crd.logger.Debug("next counter", zap.Uint64("counter", crd.counter))
		}
	}()

	// 4.1 Broadcasts the operation to all other connected clients
	// 4.2 Update local plan and persist the data
	go func() {
		for msg := range crd.queue {
			// Update local copy of plan + validate if the op is valid
			patchData, _ := json.Marshal(msg.SyncDataUpdateTrip.Value)
			patch, _ := jsonpatch.DecodePatch(patchData)
			modified, err := patch.Apply(crd.plan)
			if err != nil {
				return
			}
			crd.plan = modified
			var plan trips.TripPlan
			json.Unmarshal(crd.plan, &plan)

			// TODO: these ops must be atomic!
			crd.tripStore.SaveTripPlan(reqctx.Context{}, plan)
			crd.sesnStore.IncrSessionCounter(context.Background(), crd.planID)
			crd.tms.Publish(crd.planID, msg)
		}
	}()

	return nil
}

/***********
/* Spawner *
/***********/

type CoordinatorSpawner struct {
	sesnStore SessionStore
	sms       SyncMessageStore
	tms       TOBMessageStore
	tripStore trips.Store

	logger *zap.Logger
}

func NewCoordinatorSpawner(
	sesnStore SessionStore,
	sms SyncMessageStore,
	tms TOBMessageStore,
	tripStore trips.Store,
	logger *zap.Logger,
) *CoordinatorSpawner {
	return &CoordinatorSpawner{
		sesnStore: sesnStore,
		sms:       sms,
		tms:       tms,
		tripStore: tripStore,
		logger:    logger,
	}
}

func (spwn *CoordinatorSpawner) Run() error {
	msgCh, done, err := spwn.sms.Subscribe("*")
	if err != nil {
		return err
	}

	defer func() {
		done <- true
	}()

	for msg := range msgCh {
		if msg.OpType != SyncOpJoinSession {
			continue
		}
		sess, err := spwn.sesnStore.Read(context.Background(), msg.TripPlanID)
		if err != nil {
			// TODO: handle error here
			spwn.logger.Error("unable to read session", zap.Error(err))
			continue
		}
		if len(sess.Members) > 1 {
			// First member would have joined!
			continue
		}

		coord := NewCoordinator(
			msg.TripPlanID,
			spwn.sesnStore,
			spwn.sms,
			spwn.tms,
			spwn.tripStore,
			spwn.logger,
		)
		if err := coord.Init(); err != nil {
			spwn.logger.Error("unable to init coordinator", zap.Error(err))
			continue
		}
		if err := coord.Run(); err != nil {
			spwn.logger.Error("unable to run coordinator", zap.Error(err))
		}
	}
	return nil
}
