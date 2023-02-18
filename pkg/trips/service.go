package trips

import (
	"math/rand"
	"time"

	"github.com/tiinyplanet/tiinyplanet/pkg/images"
	"github.com/tiinyplanet/tiinyplanet/pkg/reqctx"
)

type Service interface {
	CreateTripPlan(ctx reqctx.Context, creator TripMember, name string, start, end time.Time) (TripPlan, error)
	ReadTripPlan(ctx reqctx.Context, ID string) (TripPlan, error)
	ListTripPlans(ctx reqctx.Context, ff ListTripPlansFilter) ([]TripPlan, error)
	DeleteTripPlan(ctx reqctx.Context, ID string) error
}

type service struct {
	store    Store
	imageSvc images.Service
}

func NewService(store Store, imageSvc images.Service) Service {
	return &service{store, imageSvc}
}

func (svc *service) CreateTripPlan(ctx reqctx.Context, creator TripMember, name string, start, end time.Time) (TripPlan, error) {
	plan := NewTripPlanWithDates(creator, name, start, end)
	plan.CoverImage = images.CoverStockImageList[rand.Intn(len(images.CoverStockImageList))]

	// bootstrap 1 content list
	contentList := NewTripContentList("")
	plan.Contents[contentList.ID] = contentList

	// bootstrap itinerary dates
	numDays := plan.EndDate.Sub(plan.StartDate).Hours() / 24
	for i := 0; i <= int(numDays); i++ {
		dt := plan.StartDate.Add(time.Duration(i*24) * time.Hour)
		itinList := NewItineraryList(dt)
		plan.Itinerary = append(plan.Itinerary, itinList)
	}

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
