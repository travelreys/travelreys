package trips

import (
	"context"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/travelreys/travelreys/pkg/auth"
	"github.com/travelreys/travelreys/pkg/images"
	"github.com/travelreys/travelreys/pkg/storage"
)

const ()

var (
	attachmentBucket = os.Getenv("TRAVELREYS_TRIPS_BUCKET") // "trips"
	mediaBucket      = os.Getenv("TRAVELREYS_MEDIA_BUCKET") // "media"
	mediaCDNDomain   = os.Getenv("TRAVELREYS_MEDIA_DOMAIN")
)

type Service interface {
	Create(ctx context.Context, creator Member, name string, start, end time.Time) (Trip, error)
	Read(ctx context.Context, ID string) (Trip, error)
	ReadWithMembers(ctx context.Context, ID string) (Trip, auth.UsersMap, error)
	ReadMembers(ctx context.Context, ID string) (auth.UsersMap, error)
	List(ctx context.Context, ff ListFilter) (TripsList, error)
	Delete(ctx context.Context, ID string) error

	UploadAttachmentPresignedURL(ctx context.Context, ID, fileID string) (string, error)
	DownloadAttachmentPresignedURL(ctx context.Context, ID, path, fileID string) (string, error)
	DeleteAttachment(ctx context.Context, ID string, obj storage.Object) error

	UploadMediaPresignedURL(ctx context.Context, ID, fileID string) (string, error)
	GenerateMediaPresignedCookie(ctx context.Context, ID, domain string) (*http.Cookie, error)
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

func (svc service) UploadAttachmentPresignedURL(ctx context.Context, tripID, fileID string) (string, error) {
	return svc.storageSvc.PutPresignedURL(
		ctx,
		attachmentBucket,
		filepath.Join(tripID, fileID),
		fileID)
}

func (svc service) DownloadAttachmentPresignedURL(ctx context.Context, tripID, path, fileID string) (string, error) {
	return svc.storageSvc.GetPresignedURL(ctx, attachmentBucket, path, fileID)
}

func (svc service) DeleteAttachment(ctx context.Context, tripID string, obj storage.Object) error {
	obj.Bucket = attachmentBucket
	return svc.storageSvc.Remove(ctx, obj)
}

func (svc service) UploadMediaPresignedURL(ctx context.Context, tripID, fileID string) (string, error) {
	return svc.storageSvc.PutPresignedURL(ctx, mediaBucket, filepath.Join(tripID, fileID), fileID)
}

func (svc service) GenerateMediaPresignedCookie(ctx context.Context, tripID, domain string) (*http.Cookie, error) {
	ck, err := svc.storageSvc.GeneratePresignedCookie(ctx, domain, fmt.Sprintf("/%s", tripID))
	return ck, err
}
