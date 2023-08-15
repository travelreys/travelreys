package trips

import (
	context "context"
	"errors"

	"github.com/travelreys/travelreys/pkg/common"
)

var (
	ErrInvalidOp     = errors.New("trips.ErrInvalidOp")
	ErrInvalidOpData = errors.New("trips.ErrInvalidOpData")
)

// Service handles the control & data updates made by users in the collaboration session.
type SyncService interface {
	Ping(ctx context.Context, msg *SyncMsgBroadcast) error

	Join(ctx context.Context, msg *SyncMsgTOB) error
	Leave(ctx context.Context, msg *SyncMsgTOB) error
	Update(ctx context.Context, msg *SyncMsgTOB) error

	SubSyncMsgBroadcastResp(ctx context.Context, tripID string) (<-chan SyncMsgBroadcast, chan<- bool, error)
	SubSyncMsgTOBResp(ctx context.Context, tripID string) (<-chan SyncMsgTOB, chan<- bool, error)
}

type syncService struct {
	store     Store
	sessStore SessionStore
	msgStore  SyncMsgStore
}

func NewSyncService(
	store Store,
	sessStore SessionStore,
	msgStore SyncMsgStore,
) SyncService {
	return &syncService{store, sessStore, msgStore}
}

// Broadcast

func (p *syncService) Ping(ctx context.Context, msg *SyncMsgBroadcast) error {
	sessCtx := SessionContext{
		ConnID:   msg.ConnID,
		TripID:   msg.TripID,
		MemberID: msg.MemberID,
	}
	return p.sessStore.AddSessCtx(ctx, sessCtx)
}

// TOB

func (p *syncService) Join(ctx context.Context, msg *SyncMsgTOB) error {
	trip, err := p.store.Read(ctx, msg.TripID)
	if err != nil {
		return err
	}

	if !common.StringContains(trip.GetMemberIDs(), msg.MemberID) {
		return ErrRBAC
	}
	sessCtx := SessionContext{
		ConnID:   msg.ConnID,
		TripID:   msg.TripID,
		MemberID: msg.MemberID,
	}
	if err := p.sessStore.AddSessCtx(ctx, sessCtx); err != nil {
		return err
	}
	return p.msgStore.PubTOBReq(msg.TripID, msg)
}

func (p *syncService) Leave(ctx context.Context, msg *SyncMsgTOB) error {
	sessCtx := SessionContext{
		ConnID: msg.ConnID,
		TripID: msg.TripID,
	}
	p.sessStore.RemoveSessCtx(ctx, sessCtx)
	return p.msgStore.PubTOBReq(msg.TripID, msg)
}

func (p *syncService) Update(ctx context.Context, msg *SyncMsgTOB) error {
	return p.msgStore.PubTOBReq(msg.TripID, msg)
}

func (p *syncService) SubSyncMsgBroadcastResp(
	ctx context.Context,
	tripID string,
) (<-chan SyncMsgBroadcast, chan<- bool, error) {
	return p.msgStore.SubBroadcastResp(tripID)
}

func (p *syncService) SubSyncMsgTOBResp(
	ctx context.Context,
	tripID string,
) (<-chan SyncMsgTOB, chan<- bool, error) {
	return p.msgStore.SubTOBResp(tripID)
}
