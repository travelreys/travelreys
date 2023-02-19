package trips

import (
	context "context"
	"math/rand"
	"time"

	"github.com/tiinyplanet/tiinyplanet/pkg/auth"
	"github.com/tiinyplanet/tiinyplanet/pkg/images"
)

type Service interface {
	CreateTrip(ctx context.Context, creator Member, name string, start, end time.Time) (TripPlan, error)
	ReadTrip(ctx context.Context, ID string) (TripPlan, error)
	ReadTripWithUsers(ctx context.Context, ID string) (TripPlan, auth.UsersMap, error)
	ListTrips(ctx context.Context, ff ListTripsFilter) (TripPlansList, error)
	DeleteTrip(ctx context.Context, ID string) error
}

type service struct {
	store    Store
	authSvc  auth.Service
	imageSvc images.Service
}

func NewService(store Store, authSvc auth.Service, imageSvc images.Service) Service {
	return &service{store, authSvc, imageSvc}
}

func (svc *service) CreateTrip(ctx context.Context, creator Member, name string, start, end time.Time) (TripPlan, error) {
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

func (svc *service) ReadTrip(ctx context.Context, ID string) (TripPlan, error) {
	return svc.store.ReadTrip(ctx, ID)
}

func (svc *service) ReadTripWithUsers(ctx context.Context, ID string) (TripPlan, auth.UsersMap, error) {
	trip, err := svc.ReadTrip(ctx, ID)
	if err != nil {
		return trip, nil, err
	}
	membersIDs := []string{trip.Creator.ID}
	for id := range trip.Members {
		membersIDs = append(membersIDs, id)
	}
	ff := auth.ListUsersFilter{IDs: membersIDs}
	users, err := svc.authSvc.ListUsers(ctx, ff)
	if err != nil {
		return trip, nil, err
	}
	usersMap := auth.UsersMap{}
	for _, usr := range users {
		usersMap[usr.ID] = usr
	}
	return trip, usersMap, nil

}

func (svc *service) ListTrips(ctx context.Context, ff ListTripsFilter) (TripPlansList, error) {
	return svc.store.ListTrips(ctx, ff)
}

func (svc *service) DeleteTrip(ctx context.Context, ID string) error {
	return svc.store.DeleteTrip(ctx, ID)
}
