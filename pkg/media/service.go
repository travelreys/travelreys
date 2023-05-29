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
	GenerateMediaItems(ctx context.Context, userID string, params []NewMediaItemParams) (MediaItemList, error)
	GenerateGetSignedURLs(ctx context.Context, items MediaItemList) (MediaPresignedUrlList, error)
	GeneratePutSignedURLs(ctx context.Context, items MediaItemList) (MediaPresignedUrlList, error)

	// TO DEPRECATE
	Delete(ctx context.Context, ff DeleteMediaFilter) error
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

func (svc *service) GenerateMediaItems(ctx context.Context, userID string, params []NewMediaItemParams) (MediaItemList, error) {
	items := MediaItemList{}
	for _, param := range params {
		item := NewMediaItem(userID, param)
		items = append(items, item)

	}
	return items, nil
}

func (svc *service) GeneratePutSignedURLs(ctx context.Context, items MediaItemList) (MediaPresignedUrlList, error) {
	urls := MediaPresignedUrlList{}
	for _, item := range items {
		contentURL, err := svc.storageSvc.PutPresignedURL(
			ctx, MediaItemBucket, item.Path, item.ID,
		)
		if err != nil {
			return nil, err
		}
		if item.Type == MediaTypePicture {
			urls = append(urls, MediaPresignedUrl{
				ContentURL: contentURL,
			})
			continue
		}

		previewURL, err := svc.storageSvc.PutPresignedURL(
			ctx, MediaItemOptimizedBucket, item.PreviewPath(), item.ID,
		)
		if err != nil {
			return nil, err
		}
		urls = append(urls, MediaPresignedUrl{
			ContentURL: contentURL,
			PreviewURL: previewURL,
		})
	}
	return urls, nil
}

func (svc *service) GenerateGetSignedURLs(ctx context.Context, items MediaItemList) (MediaPresignedUrlList, error) {
	urls := MediaPresignedUrlList{}
	for _, item := range items {
		contentURL, err := svc.makeCDNContentURL(ctx, item)
		if err != nil {
			return nil, err
		}
		optimizedURL, err := svc.makeCDNOptimizedURL(ctx, item)

		if item.Type == MediaTypePicture {
			urls = append(urls, MediaPresignedUrl{
				ContentURL:   contentURL,
				PreviewURL:   optimizedURL,
				OptimizedURL: optimizedURL,
			})
			continue
		}

		previewURL, err := svc.makeCDNPreviewURL(ctx, item)
		if err != nil {
			return nil, err
		}
		urls = append(urls, MediaPresignedUrl{
			ContentURL:   contentURL,
			PreviewURL:   previewURL,
			OptimizedURL: optimizedURL,
		})
	}
	return urls, nil
}

func (svc *service) makeCDNContentURL(ctx context.Context, item MediaItem) (string, error) {
	return svc.cdnProvider.PresignedURL(ctx, fmt.Sprintf(
		"%s/%s",
		svc.cdnProvider.Domain(ctx, true),
		filepath.Join(item.Bucket, item.Path),
	))
}

func (svc *service) makeCDNPreviewURL(ctx context.Context, item MediaItem) (string, error) {
	path := item.Path
	if item.Type == MediaTypeVideo {
		path = item.PreviewPath()
	}
	return svc.cdnProvider.PresignedURL(ctx, fmt.Sprintf(
		"%s/%s",
		svc.cdnProvider.Domain(ctx, true),
		filepath.Join(item.Bucket, path),
	))
}

func (svc *service) makeCDNOptimizedURL(ctx context.Context, item MediaItem) (string, error) {
	return svc.cdnProvider.PresignedURL(ctx, fmt.Sprintf(
		"%s/%s",
		svc.cdnProvider.Domain(ctx, true),
		filepath.Join(MediaItemOptimizedBucket, item.Path),
	))
}

// TO DEPRECATE

func (svc *service) Delete(ctx context.Context, ff DeleteMediaFilter) error {
	listFF := ListMediaFilter{ff.UserID, ff.IDs}
	items, _, err := svc.store.List(ctx, listFF, ListMediaPagination{})
	if err != nil {
		return err
	}
	for _, item := range items {
		if err := svc.storageSvc.Remove(ctx, item.Object); err != nil {
			svc.logger.Error("Delete", zap.Error(err))
		}
		item.Object.Path = item.PreviewPath()
		if err := svc.storageSvc.Remove(ctx, item.Object); err != nil {
			svc.logger.Error("Delete", zap.Error(err))
		}
	}
	return svc.store.Delete(ctx, ff)
}
