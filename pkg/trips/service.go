package trips

import (
	context "context"
	"time"
)

// Trips Service

type Service interface {
	CreateTripPlan(ctx context.Context, creatorID, name string, start, end time.Time) (TripPlan, error)
	ListTripPlans(ctx context.Context) ([]TripPlan, error)
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
func (svc *service) ListTripPlans(ctx context.Context) ([]TripPlan, error) {
	return svc.store.ListTripPlans(ctx)
}

func (svc *service) DeleteTripPlan(ctx context.Context, ID string) error {
	return svc.store.DeleteTripPlan(ctx, ID)
}

// Collaboration Service

type CollabSessionService interface {
	JoinSession(planID string) error
	LeaveSession(planID string) error
}

type CollabService interface {
	// Flights and Transits
	CreateFlight() error
	DeleteFlight() error
	UpdateFlight() error

	CreateTransit() error
	DeleteTransit() error
	UpdateTransit() error

	// Lodging
	CreateLodging() error
	DeleteLodging() error
	UpdateLodging() error

	// Content List
	CreateContentList() error
	DeleteContentList() error

	CreateContentBlock() error
	UpdateContentBlock() error
	DeleteContentBlock() error

	// Itinerary
	CreateItineraryContentBlock() error
	UpdateItineraryContentBlock() error
	DeleteItineraryContentBlock() error
}
