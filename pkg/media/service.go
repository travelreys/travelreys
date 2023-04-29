package media

import (
	"context"
	"fmt"
	"path/filepath"

	"github.com/travelreys/travelreys/pkg/storage"
)

type Service interface {
	GenerateMediaItems(ctx context.Context, userID string, params []NewMediaItemParams) (MediaItemList, []string, error)
	SaveForUser(ctx context.Context, userID string, items MediaItemList) error
	List(ctx context.Context, ff ListMediaFilter, pg ListMediaPagination) (MediaItemList, string, error)
	ListWithSignedURLs(ctx context.Context, ff ListMediaFilter, pg ListMediaPagination) (MediaItemList, string, []string, error)
	Delete(ctx context.Context, ff DeleteMediaFilter) error
}

type service struct {
	store       Store
	cdnProvider CDNProvider
	storageSvc  storage.Service
}

func NewService(store Store, cdnProvider CDNProvider, storageSvc storage.Service) Service {
	return &service{store, cdnProvider, storageSvc}
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

func (svc *service) List(ctx context.Context, ff ListMediaFilter, pg ListMediaPagination) (MediaItemList, string, error) {
	return svc.store.List(ctx, ff, pg)
}

func (svc *service) ListWithSignedURLs(ctx context.Context, ff ListMediaFilter, pg ListMediaPagination) (MediaItemList, string, []string, error) {
	items, lastId, err := svc.store.List(ctx, ff, pg)
	if err != nil {
		return nil, "", nil, err
	}
	urls, err := svc.GenerateGetSignedURLsForItems(ctx, items)
	return items, lastId, urls, err
}

func (svc *service) Delete(ctx context.Context, ff DeleteMediaFilter) error {
	fmt.Println(ff)
	listFF := ListMediaFilter{ff.UserID, ff.IDs}
	items, _, err := svc.store.List(ctx, listFF, ListMediaPagination{})
	if err != nil {
		return err
	}
	fmt.Println(items)
	for _, item := range items {
		go func(obj storage.Object) {
			svc.storageSvc.Remove(ctx, obj)
		}(item.Object)
	}
	return svc.store.Delete(ctx, ff)
}

func (svc *service) GenerateGetSignedURLsForItems(ctx context.Context, items MediaItemList) ([]string, error) {
	urls := []string{}
	for _, item := range items {
		cdnDomain := svc.cdnProvider.Domain(ctx, true)
		urlToSign := fmt.Sprintf("%s/%s", cdnDomain, filepath.Join(item.Bucket, item.Path))
		signedURL, err := svc.cdnProvider.PresignedURL(ctx, urlToSign)
		if err != nil {
			return nil, err
		}
		urls = append(urls, signedURL)
	}
	return urls, nil
}
