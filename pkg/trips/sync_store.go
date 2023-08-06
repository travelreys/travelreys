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

	sessStoreLogger           = "coordinator.sessStore"
	syncMsgControlStoreLogger = "coordinator.syncMsgControlStore"
	syncMsgDataStoreLogger    = "coordinator.tobStore"
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

// SubjControlRequest is the NATS.io subj for client -> coordinator communication
// for control messages
func SubjControlRequest(tripID string) string {
	return fmt.Sprintf("sync.control.requests.%s", tripID)
}

// SubjControlResponse is the NATS.io subj for client -> coordiator communication
// for control messages
func SubjControlResponse(tripID string) string {
	return fmt.Sprintf("sync.control.response.%s", tripID)
}

// SubjDataRequest is the NATS.io subj for coordinator -> client communication
// for data messages
func SubjDataRequest(tripID string) string {
	return fmt.Sprintf("sync.data.requests.%s", tripID)
}

// SubjDataResponse is the NATS.io subj for client -> coordinator communication
// for data messages
func SubjDataResponse(tripID string) string {
	return fmt.Sprintf("sync.data.response.%s", tripID)
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
		ctx, sessConnKey(sessCtx.TripID, sessCtx.ConnID),
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

type SyncMsgControlStore interface {
	PubReq(tripID string, msg SyncMsgControl) error
	SubReq(tripID string) (<-chan SyncMsgControl, chan<- bool, error)
	PubRes(tripID string, msg SyncMsgControl) error
	SubRes(tripID string) (<-chan SyncMsgControl, chan<- bool, error)
}

type syncMsgControlStore struct {
	nc     *nats.Conn
	logger *zap.Logger
}

func NewSyncMsgControlStore(nc *nats.Conn, logger *zap.Logger) SyncMsgControlStore {
	return &syncMsgControlStore{nc, logger.Named(syncMsgControlStoreLogger)}
}

func (s *syncMsgControlStore) publish(subj, tripID string, msg SyncMsgControl) error {
	data, err := msgpack.Marshal(msg)
	if err != nil {
		s.logger.Error("Publish", zap.Error(err))
	}
	s.logger.Debug("publish", zap.String("subj", subj))
	if err = s.nc.Publish(subj, data); err != nil {
		return err
	}
	return s.nc.Flush()
}

func (s *syncMsgControlStore) subscribe(subj, tripID string) (<-chan SyncMsgControl, chan<- bool, error) {
	natsCh := make(chan *nats.Msg, common.DefaultChSize)
	msgCh := make(chan SyncMsgControl, common.DefaultChSize)
	done := make(chan bool)

	sub, err := s.nc.ChanSubscribe(subj, natsCh)
	if err != nil {
		s.logger.Error("Subscribe", zap.Error(err))
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
				var msg SyncMsgControl
				if err := msgpack.Unmarshal(natsMsg.Data, &msg); err != nil {
					s.logger.Error("Subscribe", zap.Error(err))
					continue
				}
				msgCh <- msg
			}
		}
	}()
	return msgCh, done, nil
}

func (s *syncMsgControlStore) PubReq(tripID string, msg SyncMsgControl) error {
	return s.publish(SubjControlRequest(tripID), tripID, msg)
}

func (s *syncMsgControlStore) SubReq(tripID string) (<-chan SyncMsgControl, chan<- bool, error) {
	return s.subscribe(SubjControlRequest(tripID), tripID)
}

func (s *syncMsgControlStore) PubRes(tripID string, msg SyncMsgControl) error {
	return s.publish(SubjControlResponse(tripID), tripID, msg)
}

func (s *syncMsgControlStore) SubRes(tripID string) (<-chan SyncMsgControl, chan<- bool, error) {
	return s.subscribe(SubjControlResponse(tripID), tripID)
}

type SyncMsgDataStore interface {
	PubReq(tripID string, msg SyncMsgData) error
	SubReq(tripID string) (<-chan SyncMsgData, chan<- bool, error)
	SubReqQueue(tripID, groupName string) (<-chan SyncMsgData, chan<- bool, error)

	PubRes(tripID string, msg SyncMsgData) error
	SubRes(tripID string) (<-chan SyncMsgData, chan<- bool, error)
}

type syncMsgDataStore struct {
	nc *nats.Conn

	logger *zap.Logger
}

func NewSyncMsgDataStore(nc *nats.Conn, logger *zap.Logger) SyncMsgDataStore {
	return &syncMsgDataStore{nc, logger.Named(syncMsgDataStoreLogger)}
}

func (s *syncMsgDataStore) publish(subj, tripID string, msg SyncMsgData) error {
	data, _ := msgpack.Marshal(msg)

	if err := s.nc.Publish(subj, data); err != nil {
		s.logger.Error("Publish", zap.Error(err))
		return err
	}
	return s.nc.Flush()
}

func (s *syncMsgDataStore) subscribe(subj, tripID string) (<-chan SyncMsgData, chan<- bool, error) {
	natsCh := make(chan *nats.Msg, common.DefaultChSize)
	msgCh := make(chan SyncMsgData, common.DefaultChSize)
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
				var msg SyncMsgData
				if err := msgpack.Unmarshal(natsMsg.Data, &msg); err != nil {
					s.logger.Error("Subscribe", zap.Error(err))
				}
				msgCh <- msg
			}
		}
	}()
	return msgCh, done, nil
}

func (s *syncMsgDataStore) subscribeQueue(subj, tripID, groupName string) (<-chan SyncMsgData, chan<- bool, error) {
	natsCh := make(chan *nats.Msg, common.DefaultChSize)
	msgCh := make(chan SyncMsgData, common.DefaultChSize)

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
				var msg SyncMsgData
				if err := msgpack.Unmarshal(natsMsg.Data, &msg); err != nil {
					s.logger.Error("SubscribeQueue", zap.Error(err))
				}
				msgCh <- msg
			}
		}
	}()
	return msgCh, done, nil
}

func (s *syncMsgDataStore) PubReq(tripID string, msg SyncMsgData) error {
	return s.publish(SubjDataRequest(tripID), tripID, msg)
}

func (s *syncMsgDataStore) SubReq(tripID string) (<-chan SyncMsgData, chan<- bool, error) {
	return s.subscribe(SubjDataRequest(tripID), tripID)
}

func (s *syncMsgDataStore) SubReqQueue(tripID, groupName string) (<-chan SyncMsgData, chan<- bool, error) {
	return s.subscribeQueue(SubjDataRequest(tripID), tripID, groupName)
}

func (s *syncMsgDataStore) PubRes(tripID string, msg SyncMsgData) error {
	return s.publish(SubjDataResponse(tripID), tripID, msg)

}

func (s *syncMsgDataStore) SubRes(tripID string) (<-chan SyncMsgData, chan<- bool, error) {
	return s.subscribe(SubjDataResponse(tripID), tripID)
}
