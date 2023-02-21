package tripssync

import (
	context "context"
	"encoding/json"

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
	ID     string
	tripID string

	// plan is the coordinators' local copy of the trip plan
	plan []byte

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

	store     Store
	msgStore  MessageStore
	tobStore  TOBMessageStore
	tripStore trips.Store

	logger *zap.Logger
}

func NewCoordinator(
	tripID string,
	store Store,
	msgStore MessageStore,
	tobStore TOBMessageStore,
	tripStore trips.Store,

	logger *zap.Logger,
) *Coordinator {
	return &Coordinator{
		ID:        uuid.New().String(),
		tripID:    tripID,
		plan:      []byte{},
		counter:   0,
		queue:     make(chan Message, common.DefaultChSize),
		store:     store,
		msgStore:  msgStore,
		tobStore:  tobStore,
		tripStore: tripStore,
		logger:    logger.Named("sync.coordinator"),
	}
}

func (crd *Coordinator) Init() error {
	crd.logger.Info("new coordinator", zap.String("tripID", crd.tripID))

	// 1. Initialise plan for coordinator
	plan, err := crd.tripStore.Read(context.Background(), crd.tripID)
	if err != nil {
		return err
	}
	planBytes, _ := json.Marshal(plan)
	crd.plan = planBytes

	// 2.Subscribe to updates from clients
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

			switch msg.Op {
			case OpJoinSession:
				// Got new player, reset counter
				crd.counter = 0
				sess, _ := crd.store.Read(context.Background(), crd.tripID)
				msg = NewMsgMemberUpdate(msg.TripID, sess.Members)
			case OpLeaveSession:
				// stop coordinator if no users left in session
				sess, _ := crd.store.Read(context.Background(), crd.tripID)
				if len(sess.Members) == 0 {
					crd.logger.Info("stopping coordinator", zap.String("tripID", crd.tripID))
					crd.Stop()
					return
				}
				// Replace with leaving broadcast message
				msg = NewMsgMemberUpdate(msg.TripID, sess.Members)
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
			switch msg.Op {
			case OpUpdateTrip:
				// 4.2 Update local plan and persist the data (if required)
				// Update local copy of plan + validate if the op is valid
				crd.HandleTOBOpUpdateTrip(msg)
			case OpMemberUpdate:
				crd.store.ResetCounter(context.Background(), msg.TripID)
			}
			// 4.3 Broadcasts the tob msg to all other connected clients
			crd.tobStore.Publish(crd.tripID, msg)
		}
	}()
	return nil
}

// HandleOpUpdateTrip handles OpUpdateTrip messages
func (crd *Coordinator) HandleTOBOpUpdateTrip(msg Message) {
	patchOps, _ := json.Marshal(msg.Data.UpdateTrip.Ops)
	patch, _ := jsonpatch.DecodePatch(patchOps)
	modified, err := patch.Apply(crd.plan)
	if err != nil {
		crd.logger.Error("json patch apply", zap.Error(err))
		return
	}
	crd.plan = modified

	var toSave trips.Trip
	if err = json.Unmarshal(crd.plan, &toSave); err != nil {
		crd.logger.Error("json unmarshall fails", zap.Error(err))
	}
	if err = crd.tripStore.Save(context.Background(), toSave); err != nil {
		crd.logger.Error("tripStore save fails", zap.Error(err))
	}
	crd.store.IncrCounter(context.Background(), crd.tripID)
}

/***********
/* Spawner *
/***********/

type CoordinatorSpawner struct {
	store     Store
	msgStore  MessageStore
	tobStore  TOBMessageStore
	tripStore trips.Store

	logger *zap.Logger
}

func NewCoordinatorSpawner(
	store Store,
	msgStore MessageStore,
	tobStore TOBMessageStore,
	tripStore trips.Store,
	logger *zap.Logger,
) *CoordinatorSpawner {
	return &CoordinatorSpawner{
		store:     store,
		msgStore:  msgStore,
		tobStore:  tobStore,
		tripStore: tripStore,
		logger:    logger,
	}
}

// Run listens to SynOpJoinSession requests and spawn a
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
		sess, err := spwn.store.Read(context.Background(), msg.TripID)
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
			msg.TripID,
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
		}
	}
	return nil
}
