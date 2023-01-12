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

// syncSessConnectionsKey is the Redis for maintaining session connections
func syncSessConnectionsKey(planID string) string {
	return fmt.Sprintf("sync-session:%s:connections", planID)
}

// syncSessCounterKey is the Redis for TOB counter
func syncSessCounterKey(planID string) string {
	return fmt.Sprintf("sync-session:%s:counter", planID)
}

// syncSessRequestSub is the NATS.io subj for client -> coordinator communication
func syncSessRequestSubj(planID string) string {
	return fmt.Sprintf("sync-session.requests.%s", planID)
}

// syncSessTOBSubj is the NATS.io subj for coordinator -> client communication
func syncSessTOBSubj(planID string) string {
	return fmt.Sprintf("sync-session.tob.%s", planID)
}

// Session Store

type SessionStore interface {
	Read(ctx context.Context, planID string) (SyncSession, error)
	AddConnToSession(ctx context.Context, conn SyncConnection) error
	RemoveConnFromSession(ctx context.Context, conn SyncConnection) error

	GetSessionCounter(ctx context.Context, planID string) (uint64, error)
	IncrSessionCounter(ctx context.Context, planID string) error
}

type sessionStore struct {
	rdb redis.UniversalClient
}

func NewSessionStore(rdb redis.UniversalClient) SessionStore {
	return &sessionStore{rdb}
}

func (s *sessionStore) Read(ctx context.Context, planID string) (SyncSession, error) {
	var strSlice []string
	key := syncSessConnectionsKey(planID)
	cmd := s.rdb.HVals(ctx, key)
	err := cmd.ScanSlice(&strSlice)
	if err != nil {
		return SyncSession{}, err
	}

	members := trips.TripMembersList{}
	for _, str := range strSlice {
		var mem trips.TripMember
		json.Unmarshal([]byte(str), &mem)
		members = append(members, mem)
	}
	return SyncSession{members}, err
}

func (s *sessionStore) AddConnToSession(ctx context.Context, conn SyncConnection) error {
	key := syncSessConnectionsKey(conn.PlanID)
	value, _ := json.Marshal(conn.Member)
	cmd := s.rdb.HSet(ctx, key, conn.ConnectionID, value)
	return cmd.Err()
}

func (s *sessionStore) RemoveConnFromSession(ctx context.Context, conn SyncConnection) error {
	key := syncSessConnectionsKey(conn.PlanID)
	cmd := s.rdb.HDel(ctx, key, conn.ConnectionID)
	return cmd.Err()
}

func (s *sessionStore) GetSessionCounter(ctx context.Context, planID string) (uint64, error) {
	key := syncSessCounterKey(planID)
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

func (s *sessionStore) IncrSessionCounter(ctx context.Context, planID string) error {
	key := syncSessCounterKey(planID)
	cmd := s.rdb.Incr(ctx, key)
	return cmd.Err()
}

// Sync Message Store

type SyncMessageStore interface {
	Publish(planID string, msg SyncMessage) error
	Subscribe(planID string) (<-chan SyncMessage, chan<- bool, error)
	SubscribeQueue(planID, groupName string) (<-chan SyncMessage, chan<- bool, error)
}

type syncMsgStore struct {
	nc  *nats.Conn
	rdb redis.UniversalClient
}

func NewSyncMessageStore(nc *nats.Conn, rdb redis.UniversalClient) SyncMessageStore {
	return &syncMsgStore{nc, rdb}
}

func (sms *syncMsgStore) Publish(planID string, msg SyncMessage) error {
	subj := syncSessRequestSubj(planID)
	data, err := json.Marshal(msg)
	if err != nil {
		fmt.Println(err)
	}
	if err = sms.nc.Publish(subj, data); err != nil {
		return err
	}
	return sms.nc.Flush()
}

func (sms *syncMsgStore) Subscribe(planID string) (<-chan SyncMessage, chan<- bool, error) {
	subj := syncSessRequestSubj(planID)
	natsCh := make(chan *nats.Msg, common.DefaultChSize)
	msgCh := make(chan SyncMessage, common.DefaultChSize)

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
				var msg SyncMessage
				err := json.Unmarshal(natsMsg.Data, &msg)
				if err == nil {
					msgCh <- msg
				}

			}
		}
	}()
	return msgCh, done, nil
}

func (sms *syncMsgStore) SubscribeQueue(planID, groupName string) (<-chan SyncMessage, chan<- bool, error) {
	subj := syncSessRequestSubj(planID)
	natsCh := make(chan *nats.Msg, common.DefaultChSize)
	msgCh := make(chan SyncMessage, common.DefaultChSize)

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
				var msg SyncMessage
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
	Publish(planID string, msg SyncMessage) error
	Subscribe(planID string) (<-chan SyncMessage, chan<- bool, error)
}

type tobMsgStore struct {
	nc  *nats.Conn
	rdb redis.UniversalClient
}

func NewTOBMessageStore(nc *nats.Conn, rdb redis.UniversalClient) TOBMessageStore {
	return &tobMsgStore{nc, rdb}
}

func (tms *tobMsgStore) Publish(planID string, msg SyncMessage) error {
	subj := syncSessTOBSubj(planID)
	data, _ := json.Marshal(msg)

	if err := tms.nc.Publish(subj, data); err != nil {
		return err
	}
	return tms.nc.Flush()
}

func (tms *tobMsgStore) Subscribe(planID string) (<-chan SyncMessage, chan<- bool, error) {
	subj := syncSessTOBSubj(planID)
	natsCh := make(chan *nats.Msg, common.DefaultChSize)
	msgCh := make(chan SyncMessage, common.DefaultChSize)

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
				var msg SyncMessage
				err := json.Unmarshal(natsMsg.Data, &msg)
				if err == nil {
					msgCh <- msg
				}
			}
		}
	}()
	return msgCh, done, nil
}
