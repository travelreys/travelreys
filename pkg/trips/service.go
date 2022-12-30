package trips

import (
	context "context"
	"time"
)

// Trips Service

type Service interface {
	CreateTripPlan(ctx context.Context, creatorID, name string, start, end time.Time) (TripPlan, error)
	ReadTripPlan(ctx context.Context, ID string) (TripPlan, error)
	ListTripPlans(ctx context.Context, ff ListTripPlansFilter) ([]TripPlan, error)
	DeleteTripPlan(ctx context.Context, ID string) error
}

type service struct {
	store Store
}

func NewService(store Store) Service {
	return &service{store}
}

func (svc *service) CreateTripPlan(ctx context.Context, creatorID, name string, start, end time.Time) (TripPlan, error) {
	plan := NewTripPlanWithDates(creatorID, name, start, end)
	err := svc.store.SaveTripPlan(ctx, plan)
	return plan, err
}

func (svc *service) ReadTripPlan(ctx context.Context, ID string) (TripPlan, error) {
	return svc.store.ReadTripPlan(ctx, ID)
}

func (svc *service) ListTripPlans(ctx context.Context, ff ListTripPlansFilter) ([]TripPlan, error) {
	return svc.store.ListTripPlans(ctx, ff)
}

func (svc *service) DeleteTripPlan(ctx context.Context, ID string) error {
	return svc.store.DeleteTripPlan(ctx, ID)
}

// Collaboration Service

type CollabService interface {
	JoinSession(planID string) error
	LeaveSession(planID string) error
	FetchTripPlan(planID string) (TripPlan, error)
	UpdateTripPlan(planID string, patch CollabOpUpdateTripPatch)
}

type collabService struct {
	store CollabStore
}

func NewCollabService(store collabStore) CollabService {
	return &collabService{store}
}

func (svc *collabService) JoinSession(planID string) error {

}

func (svc *collabService) LeaveSession(planID string) error {

}

func (svc *collabService) FetchTripPlan(planID string) (TripPlan, error) {

}

func (svc *collabService) UpdateTripPlan(planID string, patch CollabOpUpdateTreePatch) {

}
