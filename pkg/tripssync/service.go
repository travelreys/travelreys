package tripssync

import (
	context "context"
	"errors"

	"github.com/tiinyplanet/tiinyplanet/pkg/trips"
)

var (
	ErrInvalidOp     = errors.New("service.invalid-sync-op")
	ErrInvalidOpData = errors.New("service.invalid-sync-op-data")
)

// Service handles the control & data updates made by users in the collaboration session.
type Service interface {
	JoinSession(context.Context, string, Message) (Session, error)
	LeaveSession(context.Context, string, Message) error
	UpdateTrip(context.Context, string, Message) error
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

func (p *service) JoinSession(ctx context.Context, tripID string, msg Message) (Session, error) {
	conn := Connection{
		ID:     msg.ConnID,
		TripID: tripID,
		Member: msg.Data.JoinSession.Member,
	}
	if err := p.store.AddConn(ctx, conn); err != nil {
		return Session{}, err
	}
	p.msgStore.Publish(tripID, msg)
	return p.store.Read(ctx, tripID)
}

func (p *service) LeaveSession(ctx context.Context, tripID string, msg Message) error {
	conn := Connection{
		ID:     msg.ConnID,
		TripID: tripID,
		Member: msg.Data.JoinSession.Member,
	}
	p.store.RemoveConn(ctx, conn)
	p.msgStore.Publish(tripID, msg)
	return nil
}

func (p *service) UpdateTrip(ctx context.Context, tripID string, msg Message) error {
	return p.msgStore.Publish(tripID, msg)
}

func (p *service) SubscribeTOBUpdates(ctx context.Context, tripID string) (<-chan Message, chan<- bool, error) {
	return p.tobStore.Subscribe(tripID)
}
