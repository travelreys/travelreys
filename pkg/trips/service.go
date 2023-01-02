package trips

import (
	"time"

	"github.com/awhdesmond/tiinyplanet/pkg/reqctx"
)

type Service interface {
	CreateTripPlan(ctx reqctx.Context, creator TripMember, name string, start, end time.Time) (TripPlan, error)
	ReadTripPlan(ctx reqctx.Context, ID string) (TripPlan, error)
	ListTripPlans(ctx reqctx.Context, ff ListTripPlansFilter) ([]TripPlan, error)
	DeleteTripPlan(ctx reqctx.Context, ID string) error
}

type service struct {
	store Store
}

func NewService(store Store) Service {
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
