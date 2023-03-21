package tripssync

import (
	context "context"
	"errors"

	"github.com/travelreys/travelreys/pkg/trips"
)

var (
	ErrInvalidOp     = errors.New("service.invalid-sync-op")
	ErrInvalidOpData = errors.New("service.invalid-sync-op-data")
)

// Service handles the control & data updates made by users in the collaboration session.
type Service interface {
	JoinSession(context.Context, Message) (Session, error)
	LeaveSession(context.Context, Message) error
	PingSession(context.Context, Message) error
	UpdateTrip(context.Context, Message) error
	SubscribeTOBUpdates(context.Context, string) (<-chan Message, chan<- bool, error)
}

type service struct {
	store     Store
	msgStore  MessageStore
	tobStore  TOBMessageStore
	tripStore trips.Store
}

func NewService(
	store Store,
	msgStore MessageStore,
	tobStore TOBMessageStore,
	tripStore trips.Store,
) Service {
	return &service{store, msgStore, tobStore, tripStore}
}

func (p *service) JoinSession(ctx context.Context, msg Message) (Session, error) {
	sessCtx := SessionContext{
		ID:     msg.ConnID,
		TripID: msg.TripID,
		Member: msg.Data.JoinSession.Member,
	}
	if err := p.store.AddSessCtx(ctx, sessCtx); err != nil {
		return Session{}, err
	}
	p.msgStore.Publish(msg.TripID, msg)
	return p.store.Read(ctx, msg.TripID)
}

func (p *service) LeaveSession(ctx context.Context, msg Message) error {
	sessCtx := SessionContext{
		ID:     msg.ConnID,
		TripID: msg.TripID,
	}
	p.store.RemoveSessCtx(ctx, sessCtx)
	p.msgStore.Publish(msg.TripID, msg)
	return nil
}

func (p *service) PingSession(ctx context.Context, msg Message) error {
	sessCtx := SessionContext{
		ID:     msg.ConnID,
		TripID: msg.TripID,
		Member: msg.Data.JoinSession.Member,
	}
	return p.store.AddSessCtx(ctx, sessCtx)
}

func (p *service) UpdateTrip(ctx context.Context, msg Message) error {
	return p.msgStore.Publish(msg.TripID, msg)
}

func (p *service) SubscribeTOBUpdates(ctx context.Context, tripID string) (<-chan Message, chan<- bool, error) {
	return p.tobStore.Subscribe(tripID)
}
