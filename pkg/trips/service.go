package trips

import (
	context "context"
	"math/rand"
	"time"

	"github.com/tiinyplanet/tiinyplanet/pkg/auth"
	"github.com/tiinyplanet/tiinyplanet/pkg/images"
)

type Service interface {
	Create(ctx context.Context, creator Member, name string, start, end time.Time) (Trip, error)
	Read(ctx context.Context, ID string) (Trip, error)
	ReadWithUsers(ctx context.Context, ID string) (Trip, auth.UsersMap, error)
	List(ctx context.Context, ff ListFilter) (TripsList, error)
	Delete(ctx context.Context, ID string) error
}

type service struct {
	store    Store
	authSvc  auth.Service
	imageSvc images.Service
}

func NewService(store Store, authSvc auth.Service, imageSvc images.Service) Service {
	return &service{store, authSvc, imageSvc}
}

func (svc *service) Create(ctx context.Context, creator Member, name string, start, end time.Time) (Trip, error) {
	plan := NewTripWithDates(creator, name, start, end)
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

	err := svc.store.Save(ctx, plan)
	return plan, err
}

func (svc *service) Read(ctx context.Context, ID string) (Trip, error) {
	return svc.store.Read(ctx, ID)
}

func (svc *service) ReadWithUsers(ctx context.Context, ID string) (Trip, auth.UsersMap, error) {
	trip, err := svc.Read(ctx, ID)
	if err != nil {
		return trip, nil, err
	}
	membersIDs := []string{trip.Creator.ID}
	for id := range trip.Members {
		membersIDs = append(membersIDs, id)
	}
	ff := auth.ListFilter{IDs: membersIDs}
	users, err := svc.authSvc.List(ctx, ff)
	if err != nil {
		return trip, nil, err
	}
	usersMap := auth.UsersMap{}
	for _, usr := range users {
		usersMap[usr.ID] = usr
	}
	return trip, usersMap, nil

}

func (svc *service) List(ctx context.Context, ff ListFilter) (TripsList, error) {
	return svc.store.List(ctx, ff)
}

func (svc *service) Delete(ctx context.Context, ID string) error {
	return svc.store.Delete(ctx, ID)
}
