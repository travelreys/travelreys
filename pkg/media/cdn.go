package media

import (
	"context"
	"net/http"
	"os"
	"time"

	"github.com/travelreys/travelreys/pkg/storage"
)

const (
	cdnProviderMinio                      = "minio"
	cdnProviderGcp                        = "gcpcloudcdn"
	defaultCDNProvider                    = "minio"
	defaultPresignedURLDuration           = 15 * time.Minute
	defaultPresignedCookieDuration        = 24 * time.Hour
	DefaultPresignedCookieRefreshDuration = 23 * time.Hour
)

type CDNProvider interface {
	PresignedURL(ctx context.Context, bucket, path, filename string) (string, error)
	PresignedCookie(ctx context.Context, domain, path string) (*http.Cookie, error)
}

func NewDefaultCDNProvider(ctx context.Context) (CDNProvider, error) {
	provider := os.Getenv("TRAVELREYS_CDN_PROVIDER")
	if provider == "" {
		provider = defaultCDNProvider
	}
	if provider == cdnProviderGcp {
		return NewDefaultGCPCloudCDNProvider(ctx)
	}
	return storage.NewDefaultMinioService()
}
