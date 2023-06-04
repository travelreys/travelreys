package media

import (
	"context"
	"net/http"
	"os"
	"time"
)

const (
	cdnProviderMinio               = "minio"
	cdnProviderGcp                 = "gcpcloudcdn"
	defaultCDNProvider             = "minio"
	defaultPresignedURLDuration    = 60 * time.Minute
	defaultPresignedCookieDuration = 24 * time.Hour

	presignedCookieHeader = "_travelreysCookie"
)

type CDNProvider interface {
	Domain(ctx context.Context, withScheme bool) string
	PresignedURL(ctx context.Context, url string) (string, error)
	PresignedOptURL(ctx context.Context, url string) (string, error)
	PresignedCookie(ctx context.Context, domain, path string) (*http.Cookie, error)
}

func NewDefaultCDNProvider() (CDNProvider, error) {
	provider := os.Getenv("TRAVELREYS_CDN_PROVIDER")
	if provider == "" {
		provider = defaultCDNProvider
	}
	if provider == cdnProviderGcp {
		return NewDefaultGCPCloudCDNProvider()
	}
	return NewDefaultMinioProvider()
}
