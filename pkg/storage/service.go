package storage

import (
	"context"
	"net/http"
	"os"
	"time"
)

// https://medium.com/google-cloud/using-google-cloud-storage-with-minio-object-storage-c994fe4aab6b

const (
	storageProviderMinio           = "minio"
	storageProviderGcs             = "gcs"
	defaultStorageProvider         = "minio"
	defaultPresignedURLDuration    = 30 * time.Minute
	defaultPresignedCookieDuration = 24 * time.Hour
)

type Service interface {
	Stat(ctx context.Context, bucket, path string) (Object, error)
	Remove(ctx context.Context, obj Object) error
	GetPresignedURL(ctx context.Context, bucket, path, filename string) (string, error)
	PutPresignedURL(ctx context.Context, bucket, path, filename string) (string, error)
	GeneratePresignedCookie(ctx context.Context, domain, path string) (*http.Cookie, error)
}

func NewDefaultStorageService(ctx context.Context) (Service, error) {
	provider := os.Getenv("TRAVELREYS_STORAGE_PROVIDER")
	if provider == "" {
		provider = defaultStorageProvider
	}
	if provider == storageProviderGcs {
		return NewDefaultGCSService(ctx)
	}
	return NewDefaultMinioService()
}
