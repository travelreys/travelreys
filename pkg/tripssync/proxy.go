package tripssync

import (
	context "context"
	"encoding/json"

	"github.com/go-redis/redis/v9"
	"github.com/nats-io/nats.go"
	"github.com/tiinyplanet/tiinyplanet/pkg/common"
	"github.com/tiinyplanet/tiinyplanet/pkg/trips"
)

/*********
 * Proxy *
 *********/

// Proxy proxies clients' updates to the backend collab infrastructure.
type Proxy interface {
	JoinSession(ctx context.Context, planID string, msg CollabOpMessage) (CollabSession, error)
	LeaveSession(ctx context.Context, planID string, msg CollabOpMessage) error
	ReadTripPlan(ctx context.Context, planID string, msg CollabOpMessage) (trips.TripPlan, error)
	UpdateTripPlan(ctx context.Context, planID string, msg CollabOpMessage) error
}

type proxy struct {
	store Store
}

func NewProxy(Store Store) (Proxy, error) {
	return &proxy{Store}, nil
}

func (p *proxy) Run() {}

// Session

func (p *proxy) JoinSession(ctx context.Context, planID string, msg CollabOpMessage) (CollabSession, error) {
	err := p.store.AddMemberToCollabSession(ctx, planID, msg.JoinSessionReq.TripMember)
	if err != nil {
		return CollabSession{}, err
	}

	p.store.PublishSessUpdates(ctx, planID, msg)
	return p.store.ReadCollabSession(ctx, planID)
}

func (p *proxy) LeaveSession(ctx context.Context, planID string, msg CollabOpMessage) error {
	p.store.RemoveMemberFromCollabSession(ctx, planID, msg.LeaveSessionReq.TripMember)
	p.store.PublishSessUpdates(ctx, planID, msg)
	return nil
}

// Plans

func (p *proxy) ReadTripPlan(ctx context.Context, planID string, msg CollabOpMessage) (trips.TripPlan, error) {
	return p.store.ReadTripPlan(ctx, planID)
}

func (p *proxy) UpdateTripPlan(ctx context.Context, planID string, msg CollabOpMessage) error {
	return p.store.PublishSessUpdates(ctx, planID, msg)
}

/*********
 * Store *
 *********/

type Store interface {
	ReadCollabSession(ctx context.Context, planID string) (CollabSession, error)
	AddMemberToCollabSession(ctx context.Context, planID string, member trips.TripMember) error
	RemoveMemberFromCollabSession(ctx context.Context, planID string, member trips.TripMember) error

	SubscribeTOBUpdates(ctx context.Context, planID string) (chan<- CollabOpMessage, error)
	PublishSessUpdates(ctx context.Context, planID string, msg CollabOpMessage) error

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

// Subscribe

func (s *store) SubscribeTOBUpdates(ctx context.Context, planID string) (chan<- CollabOpMessage, error) {
	subj := collabSessTOBKey(planID)
	natsCh := make(chan *nats.Msg, common.DefaultChSize)
	msgCh := make(chan CollabOpMessage, common.DefaultChSize)

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
				var msg CollabOpMessage
				err := json.Unmarshal(natsMsg.Data, &msg)
				if err == nil {
					msgCh <- msg
				}
			}
		}
	}()
	return msgCh, nil
}

// Publish

func (s *store) PublishSessUpdates(ctx context.Context, planID string, msg CollabOpMessage) error {
	subj := collabSessUpdatesKey(planID)
	data, _ := json.Marshal(msg)
	return s.nc.Publish(subj, data)
}

// Session

func (s *store) ReadCollabSession(ctx context.Context, planID string) (CollabSession, error) {
	var members trips.TripMembersList
	key := collabSessMembersKey(planID)
	data := s.rdb.SMembers(ctx, key)
	err := data.ScanSlice(&members)
	return CollabSession{members}, err
}

func (s *store) AddMemberToCollabSession(ctx context.Context, planID string, member trips.TripMember) error {
	key := collabSessMembersKey(planID)
	value, _ := json.Marshal(member)
	cmd := s.rdb.HSet(ctx, key, member.MemberID, string(value))
	return cmd.Err()
}

func (s *store) RemoveMemberFromCollabSession(ctx context.Context, planID string, member trips.TripMember) error {
	key := collabSessMembersKey(planID)
	cmd := s.rdb.HDel(ctx, key, member.MemberID)
	return cmd.Err()
}

// Plans

func (s *store) ReadTripPlan(ctx context.Context, ID string) (trips.TripPlan, error) {
	return s.tripStore.ReadTripPlan(ctx, ID)
}
