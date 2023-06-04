package media

import (
	"context"
	"net/http"
	"os"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type minioProvider struct {
	host      string
	apikey    string
	secretkey string

	mc *minio.Client
}

func NewDefaultMinioProvider() (CDNProvider, error) {
	return NewMinioProvider(
		os.Getenv("TRAVELREYS_MINIO_HOST"),
		os.Getenv("TRAVELREYS_MINIO_APIKEY"),
		os.Getenv("TRAVELREYS_MINIO_SECRETKEY"),
	)
}

func NewMinioProvider(host, apikey, secretkey string) (CDNProvider, error) {
	mc, err := minio.New(host, &minio.Options{
		Creds: credentials.NewStaticV4(apikey, secretkey, ""),
	})
	if err != nil {
		return nil, err
	}
	return &minioProvider{host, apikey, secretkey, mc}, nil
}

func (prv minioProvider) Domain(ctx context.Context, withScheme bool) string {
	domain := os.Getenv("TRAVELREYS_MEDIA_DOMAIN")
	if withScheme {
		return "http://" + domain // cdn.travelreys.com
	}
	return domain
}

func (prv minioProvider) PresignedURL(ctx context.Context, url string) (string, error) {
	return url, nil
}

func (prv minioProvider) PresignedOptURL(ctx context.Context, url string) (string, error) {
	return url, nil
}

func (prv minioProvider) PresignedCookie(ctx context.Context, domain, path string) (*http.Cookie, error) {
	return &http.Cookie{}, nil
}
