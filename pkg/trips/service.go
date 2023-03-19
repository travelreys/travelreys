package trips

import (
	context "context"
	"io"
	"math/rand"
	"path/filepath"
	"time"

	"github.com/google/uuid"
	"github.com/travelreys/travelreys/pkg/auth"
	"github.com/travelreys/travelreys/pkg/images"
	"github.com/travelreys/travelreys/pkg/storage"
)

const (
	attachmentBucket = "trips"
)

type Service interface {
	Create(ctx context.Context, creator Member, name string, start, end time.Time) (Trip, error)
	Read(ctx context.Context, ID string) (Trip, error)
	ReadWithMembers(ctx context.Context, ID string) (Trip, auth.UsersMap, error)
	ReadMembers(ctx context.Context, ID string) (auth.UsersMap, error)
	List(ctx context.Context, ff ListFilter) (TripsList, error)
	Delete(ctx context.Context, ID string) error
	Upload(ctx context.Context, ID string, filename string, filesize int64, mimeType, attachmentType string, file io.Reader) error
	Download(ctx context.Context, ID string, obj storage.Object) (storage.Object, io.ReadCloser, error)
	DeleteFile(ctx context.Context, ID string, obj storage.Object) error
}

type service struct {
	store      Store
	authSvc    auth.Service
	imageSvc   images.Service
	storageSvc storage.Service
}

func NewService(store Store, authSvc auth.Service, imageSvc images.Service, storageSvc storage.Service) Service {
	return &service{store, authSvc, imageSvc, storageSvc}
}

func (svc service) Create(ctx context.Context, creator Member, name string, start, end time.Time) (Trip, error) {
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

func (svc service) Read(ctx context.Context, ID string) (Trip, error) {
	return svc.store.Read(ctx, ID)
}

func (svc service) ReadWithMembers(ctx context.Context, ID string) (Trip, auth.UsersMap, error) {
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

func (svc service) ReadMembers(ctx context.Context, ID string) (auth.UsersMap, error) {
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

func (svc service) List(ctx context.Context, ff ListFilter) (TripsList, error) {
	return svc.store.List(ctx, ff)
}

func (svc service) Delete(ctx context.Context, ID string) error {
	return svc.store.Delete(ctx, ID)
}

func (svc service) Upload(ctx context.Context, ID string, filename string, filesize int64, mimeType, attachmentType string, file io.Reader) error {
	if _, err := svc.Read(ctx, ID); err != nil {
		return err
	}

	obj := storage.Object{
		ID:           uuid.NewString(),
		Name:         filename,
		Bucket:       attachmentBucket,
		Size:         filesize,
		Path:         filepath.Join(ID, attachmentType, filename),
		MIMEType:     mimeType,
		LastModified: time.Now(),
	}
	return svc.storageSvc.Upload(ctx, obj, file)
}

func (svc service) Download(ctx context.Context, ID string, obj storage.Object) (storage.Object, io.ReadCloser, error) {
	obj.Bucket = attachmentBucket
	stat, err := svc.storageSvc.Read(ctx, obj.Bucket, obj.Path)
	if err != nil {
		return stat, nil, err
	}
	file, err := svc.storageSvc.Download(ctx, obj)
	if err != nil {
		return stat, nil, err
	}
	return stat, file, nil
}

func (svc service) DeleteFile(ctx context.Context, ID string, obj storage.Object) error {
	obj.Bucket = attachmentBucket
	return svc.storageSvc.Remove(ctx, obj)
}
