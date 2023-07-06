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
	"go.uber.org/zap"
)

var (
	attachmentBucket = os.Getenv("TRAVELREYS_TRIPS_BUCKET")
	mediaBucket      = os.Getenv("TRAVELREYS_MEDIA_BUCKET")

	ErrDeleteAnotherTripMedia = errors.New("trips.ErrDeleteAnotherTripMedia")
)

type Service interface {
	Create(ctx context.Context, creator Member, name string, start, end time.Time) (Trip, error)
	Read(ctx context.Context, ID string) (Trip, error)
	ReadOGP(ctx context.Context, ID string) (TripOGP, error)
	ReadWithMembers(ctx context.Context, ID string) (Trip, auth.UsersMap, error)
	ReadMembers(ctx context.Context, ID string) (auth.UsersMap, error)
	List(ctx context.Context, ff ListFilter) (TripsList, error)
	ListWithMembers(ctx context.Context, ff ListFilter) (TripsList, auth.UsersMap, error)
	Delete(ctx context.Context, ID string) error

	// Attachments
	UploadAttachmentPresignedURL(ctx context.Context, ID, fileID string) (string, error)
	DownloadAttachmentPresignedURL(ctx context.Context, ID, path, fileID string) (string, error)
	DeleteAttachment(ctx context.Context, ID string, obj storage.Object) error

	// Media Items
	GenerateMediaItems(ctx context.Context, id, userID string, params []media.NewMediaItemParams) (media.MediaItemList, media.MediaPresignedUrlList, error)
	SaveMediaItems(ctx context.Context, id string, items media.MediaItemList) error
	DeleteMediaItems(ctx context.Context, id string, items media.MediaItemList) error
	GenerateGetSignedURLs(ctx context.Context, ID string, items media.MediaItemList) (media.MediaPresignedUrlList, error)
}

type service struct {
	store      Store
	authSvc    auth.Service
	imageSvc   images.Service
	mediaSvc   media.Service
	storageSvc storage.Service

	logger *zap.Logger
}

func NewService(
	store Store,
	authSvc auth.Service,
	imageSvc images.Service,
	mediaSvc media.Service,
	storageSvc storage.Service,
	logger *zap.Logger,
) Service {
	return &service{store, authSvc, imageSvc, mediaSvc, storageSvc, logger}
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
	trip, err := svc.store.Read(ctx, ID)
	if err != nil {
		return Trip{}, err
	}
	svc.augmentMediaItemURLs(ctx, &trip)
	return trip, nil
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
	contentURL, _ := svc.augmentCoverImageURL(ctx, &trip)
	return trip.OGP(creator, contentURL), nil
}

func (svc *service) ReadWithMembers(ctx context.Context, ID string) (Trip, auth.UsersMap, error) {
	trip, err := svc.readTripWithContext(ctx, ID)
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
	usersMap.Scrub()

	svc.augmentMediaItemURLs(ctx, &trip)
	return trip, usersMap, nil
}

func (svc *service) ReadMembers(ctx context.Context, ID string) (auth.UsersMap, error) {
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
			return nil, err
		}
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
	trips, err := svc.store.List(ctx, ff)
	if err != nil {
		return nil, err
	}
	for i := 0; i < len(trips); i++ {
		svc.augmentCoverImageURL(ctx, &trips[i])
	}
	return trips, nil
}

func (svc *service) ListWithMembers(ctx context.Context, ff ListFilter) (TripsList, auth.UsersMap, error) {
	trips, err := svc.List(ctx, ff)
	if err != nil {
		return nil, nil, err
	}

	usersID := []string{}
	for _, t := range trips {
		usersID = append(usersID, t.GetAllMembersID()...)
	}

	authff := auth.ListFilter{IDs: usersID}
	users, err := svc.authSvc.List(ctx, authff)
	if err != nil {
		return nil, nil, err
	}
	usersMap := auth.UsersMap{}
	for _, usr := range users {
		usersMap[usr.ID] = usr
	}
	usersMap.Scrub()

	return trips, usersMap, nil
}

// Delete performs a logical delete on the trip by updating
// the delete flag
func (svc *service) Delete(ctx context.Context, ID string) error {
	trip, err := svc.store.Read(ctx, ID)
	if err == ErrTripNotFound {
		return nil
	}
	if err != nil {
		return err
	}

	trip.Deleted = true
	return svc.store.Save(ctx, trip)
}

// Attachments

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

// MediaItems

func (svc *service) GenerateMediaItems(ctx context.Context, tripID, userID string, params []media.NewMediaItemParams) (media.MediaItemList, media.MediaPresignedUrlList, error) {
	items := media.MediaItemList{}
	for _, param := range params {
		path := filepath.Join("trips", tripID, param.Hash)
		item := media.NewMediaItem(tripID, userID, path, param)
		items = append(items, item)
	}
	urls, err := svc.mediaSvc.GeneratePutSignedURLs(ctx, items)
	return items, urls, err
}

func (svc *service) SaveMediaItems(ctx context.Context, id string, items media.MediaItemList) error {
	if _, err := svc.Read(ctx, id); err != nil {
		return err
	}
	for i := 0; i < len(items); i++ {
		items[i].TripID = id
	}
	return svc.mediaSvc.Save(ctx, items)
}

func (svc *service) DeleteMediaItems(ctx context.Context, id string, items media.MediaItemList) error {
	if _, err := svc.Read(ctx, id); err != nil {
		return err
	}

	for _, item := range items {
		if item.TripID != id {
			return ErrDeleteAnotherTripMedia
		}
	}
	return svc.mediaSvc.Delete(ctx, items)
}

func (svc *service) GenerateGetSignedURLs(ctx context.Context, ID string, items media.MediaItemList) (media.MediaPresignedUrlList, error) {
	if _, err := svc.store.Read(ctx, ID); err != nil {
		return nil, err
	}
	return svc.mediaSvc.GenerateGetSignedURLs(ctx, items)
}

func (svc *service) augmentMediaItemURLs(ctx context.Context, trip *Trip) {
	for key := range trip.MediaItems {
		urls, _ := svc.mediaSvc.GenerateGetSignedURLs(ctx, trip.MediaItems[key])
		for i := 0; i < len(trip.MediaItems[key]); i++ {
			trip.MediaItems[key][i].URLs = urls[i]
		}
	}
}

func (svc *service) augmentCoverImageURL(ctx context.Context, trip *Trip) (string, error) {
	if trip.CoverImage.Source == CoverImageSourceWeb {
		return trip.CoverImage.WebImage.Urls.Regular, nil
	}

	key, id, err := trip.CoverImage.SplitTripImageKey()
	if err != nil {
		svc.logger.Error("augmentCoverImageURL", zap.Error(err))
		return "", err
	}

	if _, ok := trip.MediaItems[key]; !ok {
		return "", nil
	}

	mediaItemIdx := -1
	for idx, item := range trip.MediaItems[key] {
		if item.ID == id {
			mediaItemIdx = idx
		}
	}
	if mediaItemIdx < 0 {
		return "", nil
	}

	urls, err := svc.GenerateGetSignedURLs(ctx, trip.ID, media.MediaItemList{
		trip.MediaItems[key][mediaItemIdx],
	})
	if err != nil {
		svc.logger.Error("augmentCoverImageURL", zap.Error(err))
		return "", err
	}
	trip.MediaItems[key][mediaItemIdx].URLs = urls[0]
	return urls[0].Image.OptimizedURL, nil
}

func (svc *service) readTripWithContext(ctx context.Context, tripID string) (Trip, error) {
	var (
		trip Trip
		err  error
	)
	ti, err := TripInfoFromCtx(ctx)
	if err == nil {
		trip = ti.Trip
	} else {
		trip, err = svc.store.Read(ctx, tripID)
		if err != nil {
			return Trip{}, err
		}
	}
	return trip, nil
}
