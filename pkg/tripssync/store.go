package tripssync

import (
	context "context"
	"encoding/json"
	"fmt"

	"github.com/go-redis/redis/v9"
	"github.com/nats-io/nats.go"
	"github.com/tiinyplanet/tiinyplanet/pkg/common"
	"github.com/tiinyplanet/tiinyplanet/pkg/trips"
)

// Pub/Sub Subjects

const (
	GroupSpawners     = "spawners"
	GroupCoordinators = "coordinators"
)

// sessConnectionsKey is the Redis key for maintaining session connections
func sessConnectionsKey(tripID string) string {
	return fmt.Sprintf("sync-session:%s:connections", tripID)
}

// sessCounterKey is the Redis for TOB counter
func sessCounterKey(tripID string) string {
	return fmt.Sprintf("sync-session:%s:counter", tripID)
}

// sessRequestSub is the NATS.io subj for client -> coordinator communication
func sessRequestSubj(tripID string) string {
	return fmt.Sprintf("sync-session.requests.%s", tripID)
}

// sessTOBSubj is the NATS.io subj for coordinator -> client communication
func sessTOBSubj(tripID string) string {
	return fmt.Sprintf("sync-session.tob.%s", tripID)
}

// Session Store

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
	rdb redis.UniversalClient
}

func NewStore(rdb redis.UniversalClient) Store {
	return &store{rdb}
}

func (s *store) Read(ctx context.Context, tripID string) (Session, error) {
	var strSlice []string
	key := sessConnectionsKey(tripID)
	cmd := s.rdb.HVals(ctx, key)
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

// Sync Message Store

type MessageStore interface {
	Publish(tripID string, msg Message) error
	Subscribe(tripID string) (<-chan Message, chan<- bool, error)
	SubscribeQueue(tripID, groupName string) (<-chan Message, chan<- bool, error)
}

type syncMsgStore struct {
	nc  *nats.Conn
	rdb redis.UniversalClient
}

func NewMessageStore(nc *nats.Conn, rdb redis.UniversalClient) MessageStore {
	return &syncMsgStore{nc, rdb}
}

func (sms *syncMsgStore) Publish(tripID string, msg Message) error {
	subj := sessRequestSubj(tripID)
	data, err := json.Marshal(msg)
	if err != nil {
		fmt.Println(err)
	}
	if err = sms.nc.Publish(subj, data); err != nil {
		return err
	}
	return sms.nc.Flush()
}

func (sms *syncMsgStore) Subscribe(tripID string) (<-chan Message, chan<- bool, error) {
	subj := sessRequestSubj(tripID)
	natsCh := make(chan *nats.Msg, common.DefaultChSize)
	msgCh := make(chan Message, common.DefaultChSize)

	done := make(chan bool)

	sub, err := sms.nc.ChanSubscribe(subj, natsCh)
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
				err := json.Unmarshal(natsMsg.Data, &msg)
				if err == nil {
					msgCh <- msg
				}

			}
		}
	}()
	return msgCh, done, nil
}

func (sms *syncMsgStore) SubscribeQueue(tripID, groupName string) (<-chan Message, chan<- bool, error) {
	subj := sessRequestSubj(tripID)
	natsCh := make(chan *nats.Msg, common.DefaultChSize)
	msgCh := make(chan Message, common.DefaultChSize)

	done := make(chan bool)

	sub, err := sms.nc.QueueSubscribeSyncWithChan(subj, groupName, natsCh)
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
				err := json.Unmarshal(natsMsg.Data, &msg)
				if err == nil {
					msgCh <- msg
				}

			}
		}
	}()
	return msgCh, done, nil
}

// TOB Message Store

type TOBMessageStore interface {
	Publish(tripID string, msg Message) error
	Subscribe(tripID string) (<-chan Message, chan<- bool, error)
}

type tobMsgStore struct {
	nc  *nats.Conn
	rdb redis.UniversalClient
}

func NewTOBMessageStore(nc *nats.Conn, rdb redis.UniversalClient) TOBMessageStore {
	return &tobMsgStore{nc, rdb}
}

func (tms *tobMsgStore) Publish(tripID string, msg Message) error {
	subj := sessTOBSubj(tripID)
	data, _ := json.Marshal(msg)

	if err := tms.nc.Publish(subj, data); err != nil {
		return err
	}
	return tms.nc.Flush()
}

func (tms *tobMsgStore) Subscribe(tripID string) (<-chan Message, chan<- bool, error) {
	subj := sessTOBSubj(tripID)
	natsCh := make(chan *nats.Msg, common.DefaultChSize)
	msgCh := make(chan Message, common.DefaultChSize)

	done := make(chan bool)

	sub, err := tms.nc.ChanSubscribe(subj, natsCh)
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
				err := json.Unmarshal(natsMsg.Data, &msg)
				if err == nil {
					msgCh <- msg
				}
			}
		}
	}()
	return msgCh, done, nil
}
