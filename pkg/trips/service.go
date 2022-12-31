package trips

import (
	"time"

	"github.com/awhdesmond/tiinyplanet/pkg/reqctx"
)

// Trips Service

type Service interface {
	CreateTripPlan(ctx reqctx.Context, creator TripMember, name string, start, end time.Time) (TripPlan, error)
	ReadTripPlan(ctx reqctx.Context, ID string) (TripPlan, error)
	ListTripPlans(ctx reqctx.Context, ff ListTripPlansFilter) ([]TripPlan, error)
	DeleteTripPlan(ctx reqctx.Context, ID string) error
}

type service struct {
	store TripStore
}

func NewService(store TripStore) Service {
	return &service{store}
}

func (svc *service) CreateTripPlan(ctx reqctx.Context, creator TripMember, name string, start, end time.Time) (TripPlan, error) {
	plan := NewTripPlanWithDates(creator, name, start, end)
	err := svc.store.SaveTripPlan(ctx, plan)
	return plan, err
}

func (svc *service) ReadTripPlan(ctx reqctx.Context, ID string) (TripPlan, error) {
	return svc.store.ReadTripPlan(ctx, ID)
}

func (svc *service) ListTripPlans(ctx reqctx.Context, ff ListTripPlansFilter) ([]TripPlan, error) {
	return svc.store.ListTripPlans(ctx, ff)
}

func (svc *service) DeleteTripPlan(ctx reqctx.Context, ID string) error {
	return svc.store.DeleteTripPlan(ctx, ID)
}

// Collaboration Service

type CollabService interface {
	FetchTripPlan(ctx reqctx.Context, planID string, msg CollabOpMessage) (TripPlan, error)

	JoinSession(ctx reqctx.Context, planID string, msg CollabOpMessage) (CollabSession, error)
	LeaveSession(ctx reqctx.Context, planID string, msg CollabOpMessage) error
	UpdateTripPlan(ctx reqctx.Context, planID string, msg CollabOpMessage) error
}

type collabService struct {
	collabStore CollabStore
	tripStore   TripStore
}

func NewCollabService(collabStore CollabStore, tripStore TripStore) (CollabService, error) {
	return &collabService{collabStore, tripStore}, nil
}

func (svc *collabService) FetchTripPlan(ctx reqctx.Context, planID string, msg CollabOpMessage) (TripPlan, error) {
	return svc.tripStore.ReadTripPlan(ctx, planID)
}

func (svc *collabService) JoinSession(ctx reqctx.Context, planID string, msg CollabOpMessage) (CollabSession, error) {
	err := svc.collabStore.AddMemberToCollabSession(ctx, planID, msg.JoinSessionReq.TripMember)
	if err != nil {
		return CollabSession{}, err
	}

	svc.collabStore.PublishCollabOpMessages(ctx, planID, msg)
	return svc.collabStore.ReadCollabSession(ctx, planID)
}

func (svc *collabService) LeaveSession(ctx reqctx.Context, planID string, msg CollabOpMessage) error {
	svc.collabStore.RemoveMemberFromCollabSession(ctx, planID, msg.LeaveSessionReq.TripMember)
	svc.collabStore.PublishCollabOpMessages(ctx, planID, msg)
	return nil
}

func (svc *collabService) UpdateTripPlan(ctx reqctx.Context, planID string, msg CollabOpMessage) error {
	return svc.collabStore.PublishCollabOpMessages(ctx, planID, msg)
}
