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
	ReadWithMembers(ctx context.Context, ID string) (Trip, auth.UsersMap, error)
	ReadMembers(ctx context.Context, ID string) (auth.UsersMap, error)
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
	trip := NewTripWithDates(creator, name, start, end)
	trip.CoverImage = images.CoverStockImageList[rand.Intn(len(images.CoverStockImageList))]

	// bootstrap 1 activity list
	activityList := NewActivityList("")
	trip.Activities[activityList.ID] = activityList

	// bootstrap itinerary dates
	numDays := trip.EndDate.Sub(trip.StartDate).Hours() / 24
	for i := 0; i <= int(numDays); i++ {
		dt := trip.StartDate.Add(time.Duration(i*24) * time.Hour)
		itin := NewItinerary(dt)
		trip.Itinerary = append(trip.Itinerary, itin)
	}
	err := svc.store.Save(ctx, trip)
	return trip, err
}

func (svc *service) Read(ctx context.Context, ID string) (Trip, error) {
	return svc.store.Read(ctx, ID)
}

func (svc *service) ReadWithMembers(ctx context.Context, ID string) (Trip, auth.UsersMap, error) {
	trip, err := svc.Read(ctx, ID)
	if err != nil {
		return trip, nil, err
	}
	ff := auth.ListFilter{IDs: trip.GetAllMembersID()}
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

func (svc *service) ReadMembers(ctx context.Context, ID string) (auth.UsersMap, error) {
	trip, err := svc.Read(ctx, ID)
	if err != nil {
		return nil, err
	}
	ff := auth.ListFilter{IDs: trip.GetAllMembersID()}
	users, err := svc.authSvc.List(ctx, ff)
	if err != nil {
		return nil, err
	}
	usersMap := auth.UsersMap{}
	for _, usr := range users {
		usersMap[usr.ID] = usr
	}
	return usersMap, nil
}

func (svc *service) List(ctx context.Context, ff ListFilter) (TripsList, error) {
	return svc.store.List(ctx, ff)
}

func (svc *service) Delete(ctx context.Context, ID string) error {
	return svc.store.Delete(ctx, ID)
}
