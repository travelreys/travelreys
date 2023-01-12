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

/*********
 * Proxy *
 *********/

// Proxy proxies clients' updates to the backend sync infrastructure.
type Proxy interface {
	JoinSession(ctx context.Context, planID string, msg SyncMessage) (SyncSession, error)
	LeaveSession(ctx context.Context, planID string, msg SyncMessage) error
	ReadTripPlan(ctx context.Context, planID string, msg SyncMessage) (trips.TripPlan, error)
	UpdateTripPlan(ctx context.Context, planID string, msg SyncMessage) error

	SubscribeTOBUpdates(ctx context.Context, planID string) (chan SyncMessage, error)
}

type proxy struct {
	store        Store
	tobUpdatesCh chan SyncMessage
}

func NewProxy(Store Store) (Proxy, error) {
	return &proxy{Store, nil}, nil
}

func (p *proxy) SubscribeTOBUpdates(ctx context.Context, planID string) (chan SyncMessage, error) {
	return p.store.SubscribeTOBUpdates(ctx, planID)
}

// Session

func (p *proxy) JoinSession(ctx context.Context, planID string, msg SyncMessage) (SyncSession, error) {
	conn := SyncConnection{
		PlanID:       planID,
		ConnectionID: msg.ID,
		Member:       msg.SyncDataJoinSession.TripMember,
	}
	err := p.store.AddConnToSession(ctx, conn)
	if err != nil {
		return SyncSession{}, err
	}

	p.store.PublishSessRequest(ctx, planID, msg)
	return p.store.ReadSyncSession(ctx, planID)
}

func (p *proxy) LeaveSession(ctx context.Context, planID string, msg SyncMessage) error {
	conn := SyncConnection{
		PlanID:       planID,
		ConnectionID: msg.ID,
	}
	p.store.RemoveConnFromSession(ctx, conn)
	p.store.PublishSessRequest(ctx, planID, msg)
	return nil
}

// Plans

func (p *proxy) ReadTripPlan(ctx context.Context, planID string, msg SyncMessage) (trips.TripPlan, error) {
	return p.store.ReadTripPlan(ctx, planID)
}

func (p *proxy) UpdateTripPlan(ctx context.Context, planID string, msg SyncMessage) error {
	return p.store.PublishSessRequest(ctx, planID, msg)
}

/*********
 * Store *
 *********/

type Store interface {
	ReadSyncSession(ctx context.Context, planID string) (SyncSession, error)
	AddConnToSession(ctx context.Context, conn SyncConnection) error
	RemoveConnFromSession(ctx context.Context, conn SyncConnection) error

	SubscribeTOBUpdates(ctx context.Context, planID string) (chan SyncMessage, error)
	PublishSessRequest(ctx context.Context, planID string, msg SyncMessage) error

	ReadTripPlan(ctx context.Context, ID string) (trips.TripPlan, error)
}

type store struct {
	tripStore trips.Store
	nc        *nats.Conn
	rdb       redis.UniversalClient

	done chan bool
}

func NewStore(tripStore trips.Store, nc *nats.Conn, rdb redis.UniversalClient) Store {
	doneCh := make(chan bool)
	return &store{tripStore, nc, rdb, doneCh}
}

// Subscribe coordinator -> client

func (s *store) SubscribeTOBUpdates(ctx context.Context, planID string) (chan SyncMessage, error) {
	subj := syncSessTOBSubj(planID)
	natsCh := make(chan *nats.Msg, common.DefaultChSize)
	msgCh := make(chan SyncMessage, common.DefaultChSize)

	sub, err := s.nc.ChanSubscribe(subj, natsCh)
	if err != nil {
		return nil, err
	}

	go func() {
		for {
			select {
			case <-s.done:
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
	return msgCh, nil
}

// Publish client -> coordinator

func (s *store) PublishSessRequest(ctx context.Context, planID string, msg SyncMessage) error {
	subj := syncSessRequestSubj(planID)
	data, _ := json.Marshal(msg)
	fmt.Println("publishing", string(data))
	return s.nc.Publish(subj, data)
}

// Session

func (s *store) ReadSyncSession(ctx context.Context, planID string) (SyncSession, error) {
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

func (s *store) AddConnToSession(ctx context.Context, conn SyncConnection) error {
	key := syncSessConnectionsKey(conn.PlanID)
	value, _ := json.Marshal(conn.Member)
	cmd := s.rdb.HSet(ctx, key, conn.ConnectionID, value)
	return cmd.Err()
}

func (s *store) RemoveConnFromSession(ctx context.Context, conn SyncConnection) error {
	key := syncSessConnectionsKey(conn.PlanID)
	cmd := s.rdb.HDel(ctx, key, conn.ConnectionID)
	fmt.Println("leaving", conn.ConnectionID)
	fmt.Println(cmd.Err())
	return cmd.Err()
}

// Plans

func (s *store) ReadTripPlan(ctx context.Context, ID string) (trips.TripPlan, error) {
	return s.tripStore.ReadTripPlan(ctx, ID)
}
