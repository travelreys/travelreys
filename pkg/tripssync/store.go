package tripssync

import (
	context "context"
	"encoding/json"
	"fmt"

	"github.com/go-redis/redis/v9"
	"github.com/nats-io/nats.go"
	"github.com/tiinyplanet/tiinyplanet/pkg/common"
	"github.com/tiinyplanet/tiinyplanet/pkg/trips"
	"go.uber.org/zap"
)

// Pub/Sub Subjects

const (
	GroupSpawners     = "spawners"
	GroupCoordinators = "coordinators"

	sessStoreLogger   = "coordinator.sessStore"
	msgStoreLogger    = "coordinator.msgStore"
	tobMsgStoreLogger = "coordinator.tobStore"
)

// sessConnectionsKey is the Redis key for maintaining session connections
func sessConnectionsKey(tripID string) string {
	return fmt.Sprintf("sync-session:%s:connections", tripID)
}

// sessCounterKey is the Redis for TOB counter
func sessCounterKey(tripID string) string {
	return fmt.Sprintf("sync-session:%s:counter", tripID)
}

// SubjRequest is the NATS.io subj for client -> coordinator communication
func SubjRequest(tripID string) string {
	return fmt.Sprintf("sync-session.requests.%s", tripID)
}

// SubjTOB is the NATS.io subj for coordinator -> client communication
func SubjTOB(tripID string) string {
	return fmt.Sprintf("sync-session.tob.%s", tripID)
}

type Store interface {
	Read(ctx context.Context, tripID string) (Session, error)
	AddConn(ctx context.Context, conn Connection) error
	RemoveConn(ctx context.Context, conn Connection) error

	GetCounter(ctx context.Context, tripID string) (uint64, error)
	IncrCounter(ctx context.Context, tripID string) error
	ResetCounter(ctx context.Context, tripID string) error
	DeleteCounter(ctx context.Context, tripID string) error
}

type store struct {
	rdb    redis.UniversalClient
	logger *zap.Logger
}

func NewStore(rdb redis.UniversalClient, logger *zap.Logger) Store {
	return &store{rdb, logger.Named(sessStoreLogger)}
}

func (s *store) Read(ctx context.Context, tripID string) (Session, error) {
	var strSlice []string
	cmd := s.rdb.HVals(ctx, sessConnectionsKey(tripID))
	err := cmd.ScanSlice(&strSlice)
	if err != nil {
		return Session{}, err
	}

	members := trips.MembersList{}
	for _, str := range strSlice {
		var mem trips.Member
		json.Unmarshal([]byte(str), &mem)
		members = append(members, mem)
	}
	return Session{members}, err
}

func (s *store) AddConn(ctx context.Context, conn Connection) error {
	key := sessConnectionsKey(conn.TripID)
	value, _ := json.Marshal(conn.Member)
	cmd := s.rdb.HSet(ctx, key, conn.ID, value)
	return cmd.Err()
}

func (s *store) RemoveConn(ctx context.Context, conn Connection) error {
	key := sessConnectionsKey(conn.TripID)
	cmd := s.rdb.HDel(ctx, key, conn.ID)
	return cmd.Err()
}

func (s *store) GetCounter(ctx context.Context, tripID string) (uint64, error) {
	key := sessCounterKey(tripID)
	cmd := s.rdb.Get(ctx, key)
	if cmd.Err() != nil {
		return 0, cmd.Err()
	}
	ctr, err := cmd.Int64()
	if err != nil {
		return 0, err
	}
	return uint64(ctr), err
}

func (s *store) IncrCounter(ctx context.Context, tripID string) error {
	key := sessCounterKey(tripID)
	cmd := s.rdb.Incr(ctx, key)
	return cmd.Err()
}

func (s *store) ResetCounter(ctx context.Context, tripID string) error {
	key := sessCounterKey(tripID)
	cmd := s.rdb.Set(ctx, key, "1", 0)
	return cmd.Err()
}

func (s *store) DeleteCounter(ctx context.Context, tripID string) error {
	key := sessCounterKey(tripID)
	cmd := s.rdb.Del(ctx, key)
	return cmd.Err()
}

type MessageStore interface {
	Publish(tripID string, msg Message) error
	Subscribe(tripID string) (<-chan Message, chan<- bool, error)
	SubscribeQueue(tripID, groupName string) (<-chan Message, chan<- bool, error)
}

type msgStore struct {
	nc  *nats.Conn
	rdb redis.UniversalClient

	logger *zap.Logger
}

func NewMessageStore(nc *nats.Conn, rdb redis.UniversalClient, logger *zap.Logger) MessageStore {
	return &msgStore{nc, rdb, logger.Named(msgStoreLogger)}
}

func (s *msgStore) Publish(tripID string, msg Message) error {
	data, err := json.Marshal(msg)
	if err != nil {
		s.logger.Error("json.Marshal(msg)", zap.Error(err))
	}
	if err = s.nc.Publish(SubjRequest(tripID), data); err != nil {
		return err
	}
	return s.nc.Flush()
}

func (s *msgStore) Subscribe(tripID string) (<-chan Message, chan<- bool, error) {
	natsCh := make(chan *nats.Msg, common.DefaultChSize)
	msgCh := make(chan Message, common.DefaultChSize)
	done := make(chan bool)

	sub, err := s.nc.ChanSubscribe(SubjRequest(tripID), natsCh)
	if err != nil {
		s.logger.Error("s.nc.ChanSubscribe", zap.Error(err))
		return nil, nil, err
	}

	go func() {
		for {
			select {
			case <-done:
				sub.Unsubscribe()
				close(msgCh)
				return
			case natsMsg := <-natsCh:
				var msg Message
				if err := json.Unmarshal(natsMsg.Data, &msg); err != nil {
					s.logger.Error("json.Unmarshal(natsMsg.Data)", zap.Error(err))
					continue
				}
				msgCh <- msg
			}
		}
	}()
	return msgCh, done, nil
}

func (s *msgStore) SubscribeQueue(tripID, groupName string) (<-chan Message, chan<- bool, error) {
	natsCh := make(chan *nats.Msg, common.DefaultChSize)
	msgCh := make(chan Message, common.DefaultChSize)

	done := make(chan bool)

	sub, err := s.nc.QueueSubscribeSyncWithChan(SubjRequest(tripID), groupName, natsCh)
	if err != nil {
		return nil, nil, err
	}

	go func() {
		for {
			select {
			case <-done:
				sub.Unsubscribe()
				close(msgCh)
				return
			case natsMsg := <-natsCh:
				var msg Message
				if err := json.Unmarshal(natsMsg.Data, &msg); err != nil {
					s.logger.Error("json.Unmarshal(natsMsg.Data)", zap.Error(err))
				}
				msgCh <- msg
			}
		}
	}()
	return msgCh, done, nil
}

type TOBMessageStore interface {
	Publish(tripID string, msg Message) error
	Subscribe(tripID string) (<-chan Message, chan<- bool, error)
}

type tobMsgStore struct {
	nc  *nats.Conn
	rdb redis.UniversalClient

	logger *zap.Logger
}

func NewTOBMessageStore(nc *nats.Conn, rdb redis.UniversalClient, logger *zap.Logger) TOBMessageStore {
	return &tobMsgStore{nc, rdb, logger.Named(tobMsgStoreLogger)}
}

func (s *tobMsgStore) Publish(tripID string, msg Message) error {
	subj := SubjTOB(tripID)
	data, _ := json.Marshal(msg)

	if err := s.nc.Publish(subj, data); err != nil {
		s.logger.Error("s.nc.Publish(subj, data)", zap.Error(err))
		return err
	}
	return s.nc.Flush()
}

func (s *tobMsgStore) Subscribe(tripID string) (<-chan Message, chan<- bool, error) {
	subj := SubjTOB(tripID)
	natsCh := make(chan *nats.Msg, common.DefaultChSize)
	msgCh := make(chan Message, common.DefaultChSize)
	done := make(chan bool)

	sub, err := s.nc.ChanSubscribe(subj, natsCh)
	if err != nil {
		return nil, nil, err
	}

	go func() {
		for {
			select {
			case <-done:
				sub.Unsubscribe()
				close(msgCh)
				return
			case natsMsg := <-natsCh:
				var msg Message
				if err := json.Unmarshal(natsMsg.Data, &msg); err != nil {
					s.logger.Error("json.Unmarshal(natsMsg.Data, &msg)", zap.Error(err))
				}
				msgCh <- msg
			}
		}
	}()
	return msgCh, done, nil
}
