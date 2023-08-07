package trips

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/travelreys/travelreys/pkg/common"
	"github.com/vmihailenco/msgpack/v5"
	clientv3 "go.etcd.io/etcd/client/v3"
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
	return fmt.Sprintf("sync-session/%s/connections/%s", tripID, connID)
}

// sessCounterKey is the Redis for TOB counter
func sessCounterKey(tripID string) string {
	return fmt.Sprintf("sync-session/%s/counter", tripID)
}

type SessionStore interface {
	AddSessCtx(ctx context.Context, sessCtx SessionContext) error
	RemoveSessCtx(ctx context.Context, sessCtx SessionContext) error
	ReadTripSessCtx(ctx context.Context, tripID string) (SessionContextList, error)

	GetCounter(ctx context.Context, tripID string) (uint64, error)
	IncrCounter(ctx context.Context, tripID string, leaseID int64) (int64, error)
	DeleteCounter(ctx context.Context, tripID string, leaseID int64) error
	RefreshCounterTTL(ctx context.Context, tripID string, leaseID int64) error
}

type sessionStore struct {
	cli    *clientv3.Client
	logger *zap.Logger
}

func NewSessionStore(cli *clientv3.Client, logger *zap.Logger) SessionStore {
	return &sessionStore{cli, logger.Named(sessStoreLogger)}
}

func (s *sessionStore) AddSessCtx(ctx context.Context, sessCtx SessionContext) error {
	value, _ := msgpack.Marshal(sessCtx)
	_, err := s.cli.Put(
		ctx,
		sessConnKey(sessCtx.TripID, sessCtx.ConnID),
		string(value),
	)
	return err
}

func (s *sessionStore) RemoveSessCtx(ctx context.Context, sessCtx SessionContext) error {
	_, err := s.cli.Delete(ctx, sessConnKey(sessCtx.TripID, sessCtx.ConnID))
	return err
}

func (s *sessionStore) ReadTripSessCtx(
	ctx context.Context,
	tripID string,
) (SessionContextList, error) {
	resp, err := s.cli.Get(ctx, sessConnKey(tripID, ""), clientv3.WithPrefix())
	if err != nil {
		return nil, err
	}

	var l SessionContextList
	for _, item := range resp.Kvs {
		var sessCtx SessionContext
		err := msgpack.Unmarshal(item.Value, &sessCtx)
		if err != nil {
			continue
		}
		l = append(l, sessCtx)
	}
	return l, nil
}

func (s *sessionStore) GetCounter(ctx context.Context, tripID string) (uint64, error) {
	resp, err := s.cli.Get(ctx, sessCounterKey(tripID))
	if err != nil {
		return 0, err
	}

	if len(resp.Kvs) == 0 {
		return 0, ErrCounterNotFound
	}

	counter, _ := strconv.Atoi(string(resp.Kvs[0].Value))
	return uint64(counter), err
}

func (s *sessionStore) IncrCounter(ctx context.Context, tripID string, leaseID int64) (int64, error) {
	key := sessCounterKey(tripID)

	ctr, err := s.GetCounter(ctx, tripID)
	if err == ErrCounterNotFound {
		leaseResp, err := s.cli.Lease.Grant(ctx, int64(defaultSyncSessionConnTTL.Seconds()))
		if err != nil {
			return 0, nil
		}

		if _, err = s.cli.Put(
			ctx, key, "1", clientv3.WithLease(leaseResp.ID),
		); err != nil {
			return 0, nil
		}
		return int64(leaseResp.ID), nil
	}

	ctr += 1
	if _, err = s.cli.Put(
		ctx, key, string(ctr), clientv3.WithLease(clientv3.LeaseID(leaseID)),
	); err != nil {
		return 0, nil
	}
	return leaseID, s.RefreshCounterTTL(ctx, tripID, leaseID)
}

func (s *sessionStore) DeleteCounter(ctx context.Context, tripID string, leaseID int64) error {
	s.cli.Revoke(ctx, clientv3.LeaseID(leaseID))

	_, err := s.cli.Delete(ctx, sessCounterKey(tripID))
	return err
}

func (s *sessionStore) RefreshCounterTTL(ctx context.Context, tripID string, leaseID int64) error {
	_, err := s.cli.KeepAliveOnce(ctx, clientv3.LeaseID(leaseID))
	return err
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
				if err := msgpack.Unmarshal(natsMsg.Data, &msg); err != nil {
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
				if err := msgpack.Unmarshal(natsMsg.Data, &msg); err != nil {
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
				if err := msgpack.Unmarshal(natsMsg.Data, &msg); err != nil {
					s.logger.Error("subTOBQueue", zap.Error(err))
				}
				msgCh <- msg
			}
		}
	}()
	return msgCh, done, nil
}

func (s *syncMsgStore) PubBroadcastReq(tripID string, msg *SyncMsgBroadcast) error {
	data, err := msgpack.Marshal(msg)
	if err != nil {
		s.logger.Error("PubBroadcastReq", zap.Error(err))
	}
	return s.publish(SubjBroadcastRequest(tripID), data)
}

func (s *syncMsgStore) SubBroadcastReq(tripID string) (<-chan SyncMsgBroadcast, chan<- bool, error) {
	return s.subBroadcast(SubjBroadcastRequest(tripID), tripID)
}

func (s *syncMsgStore) PubBroadcastResp(tripID string, msg *SyncMsgBroadcast) error {
	data, err := msgpack.Marshal(msg)
	if err != nil {
		s.logger.Error("PubBroadcastRes", zap.Error(err))
	}
	return s.publish(SubjBroadcastResponse(tripID), data)
}

func (s *syncMsgStore) SubBroadcastResp(tripID string) (<-chan SyncMsgBroadcast, chan<- bool, error) {
	return s.subBroadcast(SubjBroadcastResponse(tripID), tripID)
}

func (s *syncMsgStore) PubTOBReq(tripID string, msg *SyncMsgTOB) error {
	data, err := msgpack.Marshal(msg)
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
	data, err := msgpack.Marshal(msg)
	if err != nil {
		s.logger.Error("PubTOBRes", zap.Error(err))
	}
	return s.publish(SubjTOBResponse(tripID), data)

}

func (s *syncMsgStore) SubTOBResp(tripID string) (<-chan SyncMsgTOB, chan<- bool, error) {
	return s.subTOB(SubjTOBResponse(tripID), tripID)
}
