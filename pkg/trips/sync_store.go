package trips

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/go-redis/redis/v9"
	"github.com/nats-io/nats.go"
	"github.com/travelreys/travelreys/pkg/common"
	"go.uber.org/zap"
)

const (
	GroupSpawners     = "spawners"
	GroupCoordinators = "coordinators"

	defaultSyncSessionConnTTL = 5 * time.Minute

	sessStoreLogger    = "coordinator.sessStore"
	syncMsgStoreLogger = "coordinator.syncMsgStore"
)

var (
	ErrCounterNotFound = errors.New("trips.ErrCounterNotFound")
)

// sessConnKey is the Redis key for maintaining session connections
func sessConnKey(tripID, connID string) string {
	return fmt.Sprintf("sync-session.%s.connections.%s", tripID, connID)
}

// sessCounterKey is the Redis for TOB counter
func sessCounterKey(tripID string) string {
	return fmt.Sprintf("sync-session.%s.counter", tripID)
}

type SessionStore interface {
	AddSessCtx(ctx context.Context, sessCtx SessionContext) error
	RemoveSessCtx(ctx context.Context, sessCtx SessionContext) error
	ReadTripSessCtx(ctx context.Context, tripID string) (SessionContextList, error)

	GetCounter(ctx context.Context, tripID string) (uint64, error)
	IncrCounter(ctx context.Context, tripID string) error
	DeleteCounter(ctx context.Context, tripID string) error
	RefreshCounterTTL(ctx context.Context, tripID string) error
}

type sessionStore struct {
	rdb    redis.UniversalClient
	logger *zap.Logger
}

func NewSessionStore(rdb redis.UniversalClient, logger *zap.Logger) SessionStore {
	return &sessionStore{rdb, logger.Named(sessStoreLogger)}
}

func (s *sessionStore) AddSessCtx(ctx context.Context, sessCtx SessionContext) error {
	key := sessConnKey(sessCtx.TripID, sessCtx.ConnID)
	value, _ := common.MsgpackMarshal(sessCtx)
	cmd := s.rdb.Set(ctx, key, string(value), defaultSyncSessionConnTTL)
	return cmd.Err()
}

func (s *sessionStore) RemoveSessCtx(ctx context.Context, sessCtx SessionContext) error {
	cmd := s.rdb.Del(ctx, sessConnKey(sessCtx.TripID, sessCtx.ConnID), sessCtx.ConnID)
	return cmd.Err()
}

func (s *sessionStore) ReadTripSessCtx(
	ctx context.Context,
	tripID string,
) (SessionContextList, error) {
	var cursor uint64
	keys := []string{}
	for {
		skeys, cursor, err := s.rdb.Scan(
			ctx, cursor, sessConnKey(tripID, "*"), 10,
		).Result()
		if err != nil {
			return nil, err
		}
		for _, key := range skeys {
			keys = append(keys, key)
		}
		if cursor == 0 {
			break
		}
	}

	var l SessionContextList
	for _, key := range keys {
		str, err := s.rdb.Get(ctx, key).Result()
		if err != nil {
			return nil, err
		}
		var sessCtx SessionContext
		if err := common.MsgpackUnmarshal([]byte(str), &sessCtx); err != nil {
			continue
		}
		l = append(l, sessCtx)
	}
	return l, nil
}

func (s *sessionStore) GetCounter(ctx context.Context, tripID string) (uint64, error) {
	cmd := s.rdb.Get(ctx, sessCounterKey(tripID))
	if cmd.Err() != nil && cmd.Err() == redis.Nil {
		return 0, ErrCounterNotFound
	}
	ctr, err := cmd.Int64()
	if err != nil {
		return 0, nil
	}
	return uint64(ctr), err
}

func (s *sessionStore) IncrCounter(ctx context.Context, tripID string) error {
	incCmd := s.rdb.Incr(ctx, sessCounterKey(tripID))
	if incCmd.Err() != nil {
		return incCmd.Err()
	}
	return s.RefreshCounterTTL(ctx, tripID)
}

func (s *sessionStore) DeleteCounter(ctx context.Context, tripID string) error {
	cmd := s.rdb.Del(ctx, sessCounterKey(tripID))
	return cmd.Err()
}

func (s *sessionStore) RefreshCounterTTL(ctx context.Context, tripID string) error {
	exprCmd := s.rdb.Expire(
		ctx,
		sessCounterKey(tripID),
		defaultSyncSessionConnTTL,
	)
	return exprCmd.Err()
}

// SubjBroadcastRequest is the NATS.io subj for client -> coordinator communication
// for control messages
func SubjBroadcastRequest(tripID string) string {
	return fmt.Sprintf("sync.broadcast.requests.%s", tripID)
}

// SubjBroadcastResponse is the NATS.io subj for client -> coordiator communication
// for control messages
func SubjBroadcastResponse(tripID string) string {
	return fmt.Sprintf("sync.broadcast.response.%s", tripID)
}

// SubjTOBRequest is the NATS.io subj for coordinator -> client communication
// for data messages
func SubjTOBRequest(tripID string) string {
	return fmt.Sprintf("sync.tob.requests.%s", tripID)
}

// SubjTOBResponse is the NATS.io subj for client -> coordinator communication
// for data messages
func SubjTOBResponse(tripID string) string {
	return fmt.Sprintf("sync.tob.response.%s", tripID)
}

type SyncMsgStore interface {
	PubBroadcastReq(tripID string, msg *SyncMsgBroadcast) error
	SubBroadcastReq(tripID string) (<-chan SyncMsgBroadcast, chan<- bool, error)
	PubBroadcastResp(tripID string, msg *SyncMsgBroadcast) error
	SubBroadcastResp(tripID string) (<-chan SyncMsgBroadcast, chan<- bool, error)

	PubTOBReq(tripID string, msg *SyncMsgTOB) error
	SubTOBReq(tripID string) (<-chan SyncMsgTOB, chan<- bool, error)
	SubTOBReqQueue(tripID, groupName string) (<-chan SyncMsgTOB, chan<- bool, error)
	PubTOBResp(tripID string, msg *SyncMsgTOB) error
	SubTOBResp(tripID string) (<-chan SyncMsgTOB, chan<- bool, error)
}

type syncMsgStore struct {
	nc     *nats.Conn
	logger *zap.Logger
}

func NewSyncMsgStore(nc *nats.Conn, logger *zap.Logger) SyncMsgStore {
	return &syncMsgStore{nc, logger.Named(syncMsgStoreLogger)}
}

func (s *syncMsgStore) publish(subj string, data []byte) error {
	s.logger.Debug("publish", zap.String("subj", subj))
	if err := s.nc.Publish(subj, data); err != nil {
		return err
	}
	return s.nc.Flush()
}

func (s *syncMsgStore) subBroadcast(subj, tripID string) (<-chan SyncMsgBroadcast, chan<- bool, error) {
	natsCh := make(chan *nats.Msg, common.DefaultChSize)
	msgCh := make(chan SyncMsgBroadcast, common.DefaultChSize)
	done := make(chan bool)

	sub, err := s.nc.ChanSubscribe(subj, natsCh)
	if err != nil {
		s.logger.Error("subBroadcast", zap.Error(err))
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
				var msg SyncMsgBroadcast
				if err := common.MsgpackUnmarshal(natsMsg.Data, &msg); err != nil {
					s.logger.Error("subBroadcast", zap.Error(err))
					continue
				}
				msgCh <- msg
			}
		}
	}()
	return msgCh, done, nil
}

func (s *syncMsgStore) subTOB(subj, tripID string) (<-chan SyncMsgTOB, chan<- bool, error) {
	natsCh := make(chan *nats.Msg, common.DefaultChSize)
	msgCh := make(chan SyncMsgTOB, common.DefaultChSize)
	done := make(chan bool)

	sub, err := s.nc.ChanSubscribe(subj, natsCh)
	if err != nil {
		s.logger.Error("subTOB", zap.Error(err))
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
				var msg SyncMsgTOB
				if err := common.MsgpackUnmarshal(natsMsg.Data, &msg); err != nil {
					s.logger.Error("subTOB", zap.Error(err))
					continue
				}
				msgCh <- msg
			}
		}
	}()
	return msgCh, done, nil
}

func (s *syncMsgStore) subTOBQueue(subj, tripID, groupName string) (<-chan SyncMsgTOB, chan<- bool, error) {
	natsCh := make(chan *nats.Msg, common.DefaultChSize)
	msgCh := make(chan SyncMsgTOB, common.DefaultChSize)

	done := make(chan bool)

	sub, err := s.nc.QueueSubscribeSyncWithChan(subj, groupName, natsCh)
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
				var msg SyncMsgTOB
				if err := common.MsgpackUnmarshal(natsMsg.Data, &msg); err != nil {
					s.logger.Error("subTOBQueue", zap.Error(err))
				}
				msgCh <- msg
			}
		}
	}()
	return msgCh, done, nil
}

func (s *syncMsgStore) PubBroadcastReq(tripID string, msg *SyncMsgBroadcast) error {
	data, err := common.MsgpackMarshal(msg)
	if err != nil {
		s.logger.Error("PubBroadcastReq", zap.Error(err))
	}
	return s.publish(SubjBroadcastRequest(tripID), data)
}

func (s *syncMsgStore) SubBroadcastReq(tripID string) (<-chan SyncMsgBroadcast, chan<- bool, error) {
	return s.subBroadcast(SubjBroadcastRequest(tripID), tripID)
}

func (s *syncMsgStore) PubBroadcastResp(tripID string, msg *SyncMsgBroadcast) error {
	data, err := common.MsgpackMarshal(msg)
	if err != nil {
		s.logger.Error("PubBroadcastRes", zap.Error(err))
	}
	return s.publish(SubjBroadcastResponse(tripID), data)
}

func (s *syncMsgStore) SubBroadcastResp(tripID string) (<-chan SyncMsgBroadcast, chan<- bool, error) {
	return s.subBroadcast(SubjBroadcastResponse(tripID), tripID)
}

func (s *syncMsgStore) PubTOBReq(tripID string, msg *SyncMsgTOB) error {
	data, err := common.MsgpackMarshal(msg)
	if err != nil {
		s.logger.Error("PubTOBReq", zap.Error(err))
	}
	return s.publish(SubjTOBRequest(tripID), data)
}

func (s *syncMsgStore) SubTOBReq(tripID string) (<-chan SyncMsgTOB, chan<- bool, error) {
	return s.subTOB(SubjTOBRequest(tripID), tripID)
}

func (s *syncMsgStore) SubTOBReqQueue(tripID, groupName string) (<-chan SyncMsgTOB, chan<- bool, error) {
	return s.subTOBQueue(SubjTOBRequest(tripID), tripID, groupName)
}

func (s *syncMsgStore) PubTOBResp(tripID string, msg *SyncMsgTOB) error {
	data, err := common.MsgpackMarshal(msg)
	if err != nil {
		s.logger.Error("PubTOBRes", zap.Error(err))
	}
	return s.publish(SubjTOBResponse(tripID), data)

}

func (s *syncMsgStore) SubTOBResp(tripID string) (<-chan SyncMsgTOB, chan<- bool, error) {
	return s.subTOB(SubjTOBResponse(tripID), tripID)
}
