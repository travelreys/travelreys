package media

import (
	"context"

	"github.com/travelreys/travelreys/pkg/storage"
)

type Service interface {
	GenerateMediaItems(ctx context.Context, userID string, params []NewMediaItemParams) (MediaItemList, []string, error)
	SaveForUser(ctx context.Context, userID string, items MediaItemList) error
	List(ctx context.Context, ff ListMediaFilter) (MediaItemList, error)
	ListWithSignedURLs(ctx context.Context, ff ListMediaFilter) (MediaItemList, []string, error)
	Delete(ctx context.Context, ff DeleteMediaFilter) error
}

type service struct {
	store      Store
	storageSvc storage.Service
}

func NewService(store Store, storageSvc storage.Service) Service {
	return &service{store, storageSvc}
}

func (svc *service) GenerateMediaItems(ctx context.Context, userID string, params []NewMediaItemParams) (MediaItemList, []string, error) {
	items := MediaItemList{}
	urls := []string{}
	for _, param := range params {
		item := NewMediaItem(userID, param)
		items = append(items, item)
		signedURL, err := svc.storageSvc.PutPresignedURL(ctx, MediaItemBucket, item.Path, item.ID)
		if err != nil {
			return nil, nil, err
		}
		urls = append(urls, signedURL)
	}
	return items, urls, nil
}

func (svc *service) SaveForUser(ctx context.Context, userID string, items MediaItemList) error {
	for i := 0; i < len(items); i++ {
		items[i].UserID = userID
	}
	return svc.store.SaveForUser(ctx, userID, items)
}

func (svc *service) List(ctx context.Context, ff ListMediaFilter) (MediaItemList, error) {
	return svc.store.List(ctx, ff)
}

func (svc *service) ListWithSignedURLs(ctx context.Context, ff ListMediaFilter) (MediaItemList, []string, error) {
	items, err := svc.store.List(ctx, ff)
	if err != nil {
		return nil, nil, err
	}
	urls := []string{}
	for _, item := range items {
		signedURL, err := svc.storageSvc.GetPresignedURL(ctx, MediaItemBucket, item.Path, item.ID)
		if err != nil {
			return nil, nil, err
		}
		urls = append(urls, signedURL)
	}
	return items, urls, nil
}

func (svc *service) Delete(ctx context.Context, ff DeleteMediaFilter) error {
	return svc.store.Delete(ctx, ff)
}
