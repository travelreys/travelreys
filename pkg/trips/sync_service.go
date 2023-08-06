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
	Join(ctx context.Context, msg SyncMsgControl) error
	Leave(ctx context.Context, msg SyncMsgControl) error
	Ping(ctx context.Context, msg SyncMsgControl) error

	UpdateTrip(ctx context.Context, msg SyncMsgData) error
	SubSyncMsgDataResp(ctx context.Context, tripID string) (<-chan SyncMsgData, chan<- bool, error)
}

type syncService struct {
	store        Store
	sessStore    SessionStore
	ctrlMsgStore SyncMsgControlStore
	dataMsgStore SyncMsgDataStore
}

func NewSyncService(
	store Store,
	sessStore SessionStore,
	ctrlMsgStore SyncMsgControlStore,
	dataMsgStore SyncMsgDataStore,
) SyncService {
	return &syncService{store, sessStore, ctrlMsgStore, dataMsgStore}
}

func (p *syncService) Join(ctx context.Context, msg SyncMsgControl) error {
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

	return p.ctrlMsgStore.PubReq(msg.TripID, msg)
}

func (p *syncService) Leave(ctx context.Context, msg SyncMsgControl) error {
	sessCtx := SessionContext{
		ConnID: msg.ConnID,
		TripID: msg.TripID,
	}
	p.sessStore.RemoveSessCtx(ctx, sessCtx)
	return p.ctrlMsgStore.PubReq(msg.TripID, msg)
}

func (p *syncService) Ping(ctx context.Context, msg SyncMsgControl) error {
	sessCtx := SessionContext{
		ConnID:   msg.ConnID,
		TripID:   msg.TripID,
		MemberID: msg.MemberID,
	}
	return p.sessStore.AddSessCtx(ctx, sessCtx)
}

func (p *syncService) UpdateTrip(ctx context.Context, msg SyncMsgData) error {
	return p.dataMsgStore.PubReq(msg.TripID, msg)
}

func (p *syncService) SubSyncMsgDataResp(
	ctx context.Context,
	tripID string,
) (<-chan SyncMsgData, chan<- bool, error) {
	return p.dataMsgStore.SubRes(tripID)
}
