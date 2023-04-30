package trips

import (
	"context"
	"errors"
	"math/rand"
	"os"
	"path/filepath"
	"time"

	"github.com/travelreys/travelreys/pkg/auth"
	"github.com/travelreys/travelreys/pkg/images"
	"github.com/travelreys/travelreys/pkg/media"
	"github.com/travelreys/travelreys/pkg/storage"
)

var (
	attachmentBucket = os.Getenv("TRAVELREYS_TRIPS_BUCKET")
	mediaBucket      = os.Getenv("TRAVELREYS_MEDIA_BUCKET")

	ErrTripSharingNotEnabled = errors.New("trip.service.tripSharingNotEnabled")
)

type Service interface {
	Create(ctx context.Context, creator Member, name string, start, end time.Time) (Trip, error)
	Read(ctx context.Context, ID string) (Trip, error)
	ReadShare(ctx context.Context, ID string) (Trip, auth.UsersMap, error)
	ReadOGP(ctx context.Context, ID string) (TripOGP, error)
	ReadWithMembers(ctx context.Context, ID string) (Trip, auth.UsersMap, error)
	ReadMembers(ctx context.Context, ID string) (auth.UsersMap, error)
	List(ctx context.Context, ff ListFilter) (TripsList, error)
	Delete(ctx context.Context, ID string) error

	UploadAttachmentPresignedURL(ctx context.Context, ID, fileID string) (string, error)
	DownloadAttachmentPresignedURL(ctx context.Context, ID, path, fileID string) (string, error)
	DeleteAttachment(ctx context.Context, ID string, obj storage.Object) error

	GenerateMediaItems(ctx context.Context, userID string, params []media.NewMediaItemParams) (media.MediaItemList, []string, error)
	GenerateSignedURLs(ctx context.Context, ID string, items media.MediaItemList) ([]string, error)
}

type service struct {
	store      Store
	authSvc    auth.Service
	imageSvc   images.Service
	mediaSvc   media.Service
	storageSvc storage.Service
}

func NewService(
	store Store,
	authSvc auth.Service,
	imageSvc images.Service,
	mediaSvc media.Service,
	storageSvc storage.Service,
) Service {
	return &service{store, authSvc, imageSvc, mediaSvc, storageSvc}
}

func (svc *service) Create(ctx context.Context, creator Member, name string, start, end time.Time) (Trip, error) {
	trip := NewTripWithDates(creator, name, start, end)
	trip.CoverImage = CoverImage{
		Source:   CoverImageSourceWeb,
		WebImage: images.CoverStockImageList[rand.Intn(len(images.CoverStockImageList))],
	}

	// bootstrap itinerary dates
	numDays := trip.EndDate.Sub(trip.StartDate).Hours() / 24
	for i := 0; i <= int(numDays); i++ {
		dt := trip.StartDate.Add(time.Duration(i*24) * time.Hour)
		itin := NewItinerary(dt)
		trip.Itineraries[dt.Format("2006-01-02")] = itin
	}
	err := svc.store.Save(ctx, trip)
	return trip, err
}

func (svc *service) Read(ctx context.Context, ID string) (Trip, error) {
	return svc.store.Read(ctx, ID)
}

func (svc *service) ReadShare(ctx context.Context, ID string) (Trip, auth.UsersMap, error) {
	var (
		trip Trip
		err  error
	)

	ti, err := TripInfoFromCtx(ctx)
	if err == nil {
		trip = ti.Trip
	} else {
		trip, err = svc.store.Read(ctx, ID)
		if err != nil {
			return Trip{}, nil, err
		}
	}

	ff := auth.ListFilter{IDs: []string{trip.Creator.ID}}
	users, err := svc.authSvc.List(ctx, ff)
	if err != nil {
		return trip, nil, err
	}
	usersMap := auth.UsersMap{}
	for _, usr := range users {
		if usr.ID != trip.Creator.ID {
			continue
		}
		usersMap[usr.ID] = usr
	}
	return trip.PublicInfo(), usersMap, nil
}

func (svc *service) ReadOGP(ctx context.Context, ID string) (TripOGP, error) {
	trip, err := svc.store.Read(ctx, ID)
	if err != nil {
		return TripOGP{}, err
	}
	ff := auth.ListFilter{IDs: []string{trip.Creator.ID}}
	users, err := svc.authSvc.List(ctx, ff)
	if err != nil {
		return TripOGP{}, err
	}
	var creator auth.User
	for _, usr := range users {
		if usr.ID == trip.Creator.ID {
			creator = usr
			break
		}
	}

	// Select CoverImage URL
	coverImageURL := trip.CoverImage.WebImage.Urls.Regular
	if trip.CoverImage.Source == CoverImageSourceTrip {
		mediaItem := trip.MediaItems[MediaItemKeyCoverImage][trip.CoverImage.TripImage]
		urls, err := svc.GenerateSignedURLs(ctx, trip.ID, media.MediaItemList{mediaItem})
		if err != nil {
			return TripOGP{}, err
		}
		coverImageURL = urls[0]
	}

	return trip.OGP(creator, coverImageURL), nil
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

func (svc *service) UploadAttachmentPresignedURL(ctx context.Context, tripID, fileID string) (string, error) {
	return svc.storageSvc.PutPresignedURL(
		ctx,
		attachmentBucket,
		filepath.Join(tripID, fileID),
		fileID)
}

func (svc *service) DownloadAttachmentPresignedURL(ctx context.Context, tripID, path, fileID string) (string, error) {
	return svc.storageSvc.GetPresignedURL(ctx, attachmentBucket, path, fileID)
}

func (svc *service) DeleteAttachment(ctx context.Context, tripID string, obj storage.Object) error {
	obj.Bucket = attachmentBucket
	return svc.storageSvc.Remove(ctx, obj)
}

func (svc *service) GenerateMediaItems(ctx context.Context, userID string, params []media.NewMediaItemParams) (media.MediaItemList, []string, error) {
	return svc.mediaSvc.GenerateMediaItems(ctx, userID, params)
}

func (svc *service) GenerateSignedURLs(ctx context.Context, ID string, items media.MediaItemList) ([]string, error) {
	return svc.mediaSvc.GenerateGetSignedURLsForItems(ctx, items)
}
