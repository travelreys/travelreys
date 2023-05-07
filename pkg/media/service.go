package media

import (
	"context"
	"fmt"
	"path/filepath"

	"github.com/travelreys/travelreys/pkg/storage"
	"go.uber.org/zap"
)

const (
	svcLoggerName = "media.svc"
)

type Service interface {
	GenerateMediaItems(ctx context.Context, userID string, params []NewMediaItemParams) (MediaItemList, MediaPresignedUrlList, error)
	SaveForUser(ctx context.Context, userID string, items MediaItemList) error
	List(ctx context.Context, ff ListMediaFilter, pg ListMediaPagination) (MediaItemList, string, error)
	ListWithSignedURLs(ctx context.Context, ff ListMediaFilter, pg ListMediaPagination) (MediaItemList, string, MediaPresignedUrlList, error)
	Delete(ctx context.Context, ff DeleteMediaFilter) error
	GenerateGetSignedURLsForItems(ctx context.Context, items MediaItemList) (MediaPresignedUrlList, error)
}

type service struct {
	store       Store
	cdnProvider CDNProvider
	storageSvc  storage.Service
	logger      *zap.Logger
}

func NewService(store Store, cdnProvider CDNProvider, storageSvc storage.Service, logger *zap.Logger) Service {
	return &service{store, cdnProvider, storageSvc, logger.Named(svcLoggerName)}
}

func (svc *service) GenerateMediaItems(ctx context.Context, userID string, params []NewMediaItemParams) (MediaItemList, MediaPresignedUrlList, error) {
	items := MediaItemList{}
	urls := MediaPresignedUrlList{}
	for _, param := range params {
		item := NewMediaItem(userID, param)
		items = append(items, item)
		signedURL, err := svc.storageSvc.PutPresignedURL(
			ctx, MediaItemBucket, item.Path, item.ID,
		)
		if err != nil {
			return nil, nil, err
		}

		if param.Type == MediaTypePicture {
			urls = append(urls, MediaPresignedUrl{
				ContentURL: signedURL,
				PreviewURL: "",
			})
			continue
		}
		previewURL, err := svc.storageSvc.PutPresignedURL(
			ctx, MediaItemBucket, item.PreviewPath(), item.ID,
		)
		if err != nil {
			return nil, nil, err
		}

		urls = append(urls, MediaPresignedUrl{
			ContentURL: signedURL,
			PreviewURL: previewURL,
		})
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

func (svc *service) ListWithSignedURLs(ctx context.Context, ff ListMediaFilter, pg ListMediaPagination) (MediaItemList, string, MediaPresignedUrlList, error) {
	items, lastId, err := svc.store.List(ctx, ff, pg)
	if err != nil {
		return nil, "", nil, err
	}

	urls, err := svc.GenerateGetSignedURLsForItems(ctx, items)
	return items, lastId, urls, err
}

func (svc *service) Delete(ctx context.Context, ff DeleteMediaFilter) error {
	listFF := ListMediaFilter{ff.UserID, ff.IDs}
	items, _, err := svc.store.List(ctx, listFF, ListMediaPagination{})
	if err != nil {
		return err
	}
	for _, item := range items {
		if err := svc.storageSvc.Remove(ctx, item.Object); err != nil {
			svc.logger.Error("delete", zap.Error(err))
		}
		item.Object.Path = item.PreviewPath()
		if err := svc.storageSvc.Remove(ctx, item.Object); err != nil {
			svc.logger.Error("delete preview", zap.Error(err))
		}
	}
	return svc.store.Delete(ctx, ff)
}

func (svc *service) GenerateGetSignedURLsForItems(ctx context.Context, items MediaItemList) (MediaPresignedUrlList, error) {
	urls := MediaPresignedUrlList{}
	for _, item := range items {
		signedURL, err := svc.cdnProvider.PresignedURL(ctx, fmt.Sprintf(
			"%s/%s",
			svc.cdnProvider.Domain(ctx, true),
			filepath.Join(item.Bucket, item.Path),
		))
		if err != nil {
			return nil, err
		}
		if item.Type == MediaTypePicture {
			urls = append(urls, MediaPresignedUrl{
				ContentURL: signedURL,
				PreviewURL: signedURL,
			})
			continue
		}

		previewURL, err := svc.cdnProvider.PresignedURL(ctx, fmt.Sprintf(
			"%s/%s",
			svc.cdnProvider.Domain(ctx, true),
			filepath.Join(item.Bucket, item.PreviewPath()),
		))
		if err != nil {
			return nil, err
		}
		urls = append(urls, MediaPresignedUrl{
			ContentURL: signedURL,
			PreviewURL: previewURL,
		})
	}
	return urls, nil
}
