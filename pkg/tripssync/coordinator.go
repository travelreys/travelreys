package tripssync

import (
	context "context"
	"encoding/json"
	"fmt"

	jsonpatch "github.com/evanphx/json-patch/v5"
	"github.com/google/uuid"
	"github.com/tiinyplanet/tiinyplanet/pkg/common"
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
	crd.logger.Debug("loaded plan", zap.String("plan", string(crd.plan)))

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
	crd.sesnStore.DeleteSessionCounter(context.Background(), crd.planID)
	if crd.doneCh != nil {
		crd.doneCh <- true
	}
	close(crd.queue)
}

func (crd *Coordinator) Run() error {
	crd.logger.Info("running coordinator", zap.String("planID", crd.planID))

	go func() {
		// 3.1. Takes in op msg indicating changes from clients
		for msg := range crd.msgCh {
			crd.logger.Debug("recv msg", zap.String("msg", fmt.Sprintf("%+v", msg)))

			switch msg.OpType {
			case SyncOpJoinSession:
				// Got new player, reset counter
				crd.counter = 0
				sess, _ := crd.sesnStore.Read(context.Background(), crd.planID)
				msg = NewSyncMessageJoinSessionBroadcast(msg.TripPlanID, sess.Members)
			case SyncOpLeaveSession:
				// stop coordinator if no users left in session
				sess, _ := crd.sesnStore.Read(context.Background(), crd.planID)
				if len(sess.Members) == 0 {
					crd.logger.Info("stopping coordinator", zap.String("planID", crd.planID))
					crd.Stop()
					return
				}
				// Replace with leaving broadcast message
				msg = NewSyncMessageLeaveSessionBroadcast(msg.TripPlanID, sess.Members)
			}

			// 3.2. Give each message a counter
			msg.Counter = crd.counter

			// 3.3. Sends each ordered message back to the FIFO queue
			crd.queue <- msg

			// Need to update the counter
			crd.counter++
			crd.logger.Debug("next counter", zap.Uint64("counter", crd.counter))
		}
	}()

	go func() {
		// 4.1 Read message from FIFO Queue
		for msg := range crd.queue {
			// 4.2 Update local plan and persist the data (if required)
			// Update local copy of plan + validate if the op is valid
			if msg.OpType == SyncOpUpdateTrip {
				crd.HandleTOBSyncOpUpdateTrip(msg)
			}
			// 4.3 Broadcasts the tob msg to all other connected clients
			crd.tms.Publish(crd.planID, msg)
		}
	}()

	return nil
}

// HandleSyncOpUpdateTrip handles SyncOpUpdateTrip messages
func (crd *Coordinator) HandleTOBSyncOpUpdateTrip(msg SyncMessage) {
	opList := []interface{}{msg.SyncDataUpdateTrip}
	patchJSON, _ := json.Marshal(opList)
	patch, _ := jsonpatch.DecodePatch(patchJSON)

	modified, err := patch.Apply(crd.plan)
	if err != nil {
		crd.logger.Error("json patch apply", zap.Error(err))
		return
	}
	crd.plan = modified
	crd.logger.Debug("modified plan", zap.String("modified", string(modified)))

	var toSave trips.TripPlan
	json.Unmarshal(crd.plan, &toSave)

	crd.tripStore.SaveTripPlan(context.Background(), toSave)
	crd.sesnStore.IncrSessionCounter(context.Background(), crd.planID)

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

// Run listens to SynOpJoinSession requests and spawn a
// coordinator if there is only 1 member in the session (i.e new session)
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
