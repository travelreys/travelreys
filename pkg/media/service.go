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
	GenerateGetSignedURLs(ctx context.Context, items MediaItemList) (MediaPresignedUrlList, error)
	GeneratePutSignedURLs(ctx context.Context, items MediaItemList) (MediaPresignedUrlList, error)
	Save(ctx context.Context, items MediaItemList) error
	Delete(ctx context.Context, items MediaItemList) error
}

type service struct {
	store       Store
	cdnProvider CDNProvider
	storageSvc  storage.Service

	logger *zap.Logger
}

func NewService(
	store Store,
	cdnProvider CDNProvider,
	storageSvc storage.Service,
	logger *zap.Logger,
) Service {
	return &service{store, cdnProvider, storageSvc, logger.Named(svcLoggerName)}
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
			Video: VideoPresignedUrls{
				PreviewURL: previewURL,
			},
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
				ContentURL: contentURL,
				Image: ImagePresignedUrls{
					OptimizedURL: optimizedURL,
				},
			})
			continue
		}

		previewURL, err := svc.makeCDNPreviewURL(ctx, item)
		if err != nil {
			return nil, err
		}

		vidH264 := MediaItem{Object: storage.Object{Path: item.VideoH264Path()}}
		sourceH264, err := svc.makeCDNOptimizedURL(ctx, vidH264)
		if err != nil {
			return nil, err
		}
		vidH265 := MediaItem{Object: storage.Object{Path: item.VideoH265Path()}}
		sourceH265, err := svc.makeCDNOptimizedURL(ctx, vidH265)
		if err != nil {
			return nil, err
		}

		urls = append(urls, MediaPresignedUrl{
			ContentURL: contentURL,
			Video: VideoPresignedUrls{
				PreviewURL: previewURL,
				Sources: []VideoSource{
					{Source: sourceH265, Codecs: "hvc1"},
					{Source: sourceH264, Codecs: "avc1"},
				},
			},
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
		filepath.Join(MediaItemOptimizedBucket, path),
	))
}

func (svc *service) makeCDNOptimizedURL(ctx context.Context, item MediaItem) (string, error) {
	return svc.cdnProvider.PresignedURL(ctx, fmt.Sprintf(
		"%s/%s",
		svc.cdnProvider.Domain(ctx, true),
		filepath.Join(MediaItemOptimizedBucket, item.Path),
	))
}

func (svc *service) Save(ctx context.Context, items MediaItemList) error {
	return svc.store.Save(ctx, items)
}

func (svc *service) Delete(ctx context.Context, items MediaItemList) error {
	ids := []string{}
	for _, item := range items {
		ids = append(ids, item.ID)
		// content
		if err := svc.storageSvc.Remove(ctx, item.Object); err != nil {
			svc.logger.Error("Delete", zap.Error(err))
		}

		if item.Type == MediaTypePicture {
			// optimised
			item.Object.Bucket = MediaItemOptimizedBucket
			if err := svc.storageSvc.Remove(ctx, item.Object); err != nil {
				svc.logger.Error("Delete", zap.Error(err))
			}
		} else {
			// preview
			item.Object.Path = item.PreviewPath()
			if err := svc.storageSvc.Remove(ctx, item.Object); err != nil {
				svc.logger.Error("Delete", zap.Error(err))
			}
			// preview in optimised
			item.Object.Bucket = MediaItemOptimizedBucket
			if err := svc.storageSvc.Remove(ctx, item.Object); err != nil {
				svc.logger.Error("Delete", zap.Error(err))
			}
			// codecs in optimised
			item.Object.Path = item.VideoH264Path()
			if err := svc.storageSvc.Remove(ctx, item.Object); err != nil {
				svc.logger.Error("Delete", zap.Error(err))
			}
			item.Object.Path = item.VideoH265Path()
			if err := svc.storageSvc.Remove(ctx, item.Object); err != nil {
				svc.logger.Error("Delete", zap.Error(err))
			}
		}
	}

	return svc.store.Delete(ctx, DeleteMediaFilter{IDs: ids})
}
