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
	attachmentBucket          = os.Getenv("TRAVELREYS_TRIPS_BUCKET")
	ErrDeleteAnotherTripMedia = errors.New("trips.ErrDeleteAnotherTripMedia")
)

type Service interface {
	Create(ctx context.Context, creatorID, name string, start, end time.Time) (*Trip, error)
	Save(ctx context.Context, trip *Trip) error

	Read(ctx context.Context, ID string) (*Trip, error)
	ReadOGP(ctx context.Context, ID string) (TripOGP, error)
	ReadMembers(ctx context.Context, ID string) (MembersMap, error)
	List(ctx context.Context, ff ListFilter) (TripsList, error)
	ListWithMembers(ctx context.Context, ff ListFilter) (TripsList, auth.UsersMap, error)

	Delete(ctx context.Context, ID string) error

	// Attachments
	UploadAttachmentPresignedURL(ctx context.Context, ID, fileID, fileType string) (string, error)
	DownloadAttachmentPresignedURL(ctx context.Context, ID, path, fileID string) (string, error)
	DeleteAttachment(ctx context.Context, ID string, obj storage.Object) error

	// Media Items
	GenerateMediaItems(ctx context.Context, ID, userID string, params []media.NewMediaItemParams) (media.MediaItemList, media.MediaPresignedUrlList, error)
	SaveMediaItems(ctx context.Context, ID string, items media.MediaItemList) error
	DeleteMediaItems(ctx context.Context, ID string, items media.MediaItemList) error
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

func (svc *service) tripFromContext(ctx context.Context, ID string) (*Trip, error) {
	var (
		trip *Trip
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
	return trip, err
}

func (svc *service) Create(
	ctx context.Context,
	creatorID, name string,
	start, end time.Time,
) (*Trip, error) {
	trip := NewTripWithDates(NewCreator(creatorID), name, start, end)

	// pick random cover image
	trip.CoverImage = &CoverImage{
		Source:   CoverImageSourceWeb,
		WebImage: images.CoverStockImageList[rand.Intn(len(images.CoverStockImageList))],
	}

	// bootstrap itinerary dates
	numDays := trip.EndDate.Sub(trip.StartDate).Hours() / 24
	for i := 0; i <= int(numDays); i++ {
		dt := trip.StartDate.Add(time.Duration(i*24) * time.Hour)
		itin := NewItinerary(dt)
		trip.Itineraries[dt.Format(ItineraryDtKeyFormat)] = itin
	}
	err := svc.Save(ctx, trip)
	return trip, err
}

func (svc *service) Save(ctx context.Context, trip *Trip) error {
	return svc.store.Save(ctx, trip)
}

func (svc *service) Read(ctx context.Context, ID string) (*Trip, error) {
	trip, err := svc.store.Read(ctx, ID)
	if err != nil {
		return nil, err
	}
	svc.augmentMediaItemURLs(ctx, trip)
	return trip, nil
}

func (svc *service) ReadOGP(ctx context.Context, ID string) (TripOGP, error) {
	trip, err := svc.store.Read(ctx, ID)
	if err != nil {
		return TripOGP{}, err
	}

	creator, err := svc.authSvc.Read(ctx, trip.Creator.ID)
	if err != nil {
		return TripOGP{}, err
	}

	contentURL, _ := svc.augmentCoverImageURL(ctx, trip)
	return trip.ToOGP(creator.Username, contentURL), nil
}

func (svc *service) ReadMembers(
	ctx context.Context,
	ID string,
) (MembersMap, error) {
	trip, err := svc.tripFromContext(ctx, ID)
	if err != nil {
		return nil, err
	}

	ff := auth.ListFilter{IDs: trip.GetMemberIDs()}
	users, err := svc.authSvc.List(ctx, ff)
	if err != nil {
		return nil, err
	}

	for _, usr := range users {
		if usr.ID == trip.Creator.ID {
			trip.Creator.augmentMemberWithUser(usr)
			continue
		}
		trip.Members[usr.ID].augmentMemberWithUser(usr)
	}

	trip.Members[trip.Creator.ID] = &trip.Creator
	return trip.Members, nil
}

func (svc *service) List(ctx context.Context, ff ListFilter) (TripsList, error) {
	trips, err := svc.store.List(ctx, ff)
	if err != nil {
		return nil, err
	}
	for i := 0; i < len(trips); i++ {
		svc.augmentCoverImageURL(ctx, trips[i])
	}
	return trips, nil
}

func (svc *service) ListWithMembers(
	ctx context.Context,
	ff ListFilter,
) (TripsList, auth.UsersMap, error) {
	trips, err := svc.List(ctx, ff)
	if err != nil {
		return nil, nil, err
	}

	usersID := []string{}
	for _, t := range trips {
		usersID = append(usersID, t.GetMemberIDs()...)
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

// Delete performs a logical delete on the trip
// by updating the delete flag
func (svc *service) Delete(ctx context.Context, ID string) error {
	trip, err := svc.tripFromContext(ctx, ID)
	if err != nil {
		return err
	}
	trip.Delete()
	return svc.store.Save(ctx, trip)
}

// Attachments

func (svc *service) UploadAttachmentPresignedURL(ctx context.Context, tripID, fileID, fileType string) (string, error) {
	return svc.storageSvc.PutPresignedURL(
		ctx, attachmentBucket, filepath.Join(tripID, fileID), fileID, fileType,
	)
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

func (svc *service) SaveMediaItems(ctx context.Context, ID string, items media.MediaItemList) error {
	if _, err := svc.Read(ctx, ID); err != nil {
		return err
	}
	for i := 0; i < len(items); i++ {
		items[i].TripID = ID
	}
	return svc.mediaSvc.Save(ctx, items)
}

func (svc *service) DeleteMediaItems(ctx context.Context, ID string, items media.MediaItemList) error {
	if _, err := svc.Read(ctx, ID); err != nil {
		return err
	}

	for _, item := range items {
		if item.TripID != ID {
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
