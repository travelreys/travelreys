package storage

import (
	"context"
	"fmt"
	"net/url"
	"os"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

// https://medium.com/google-cloud/using-google-cloud-storage-with-minio-object-storage-c994fe4aab6b

const (
	envHost      = "TRAVELREYS_STORAGE_HOST"
	envApiKey    = "TRAVELREYS_STORAGE_APIKEY"
	envSecretKey = "TRAVELREYS_STORAGE_SECRETKEY"

	defaultPresignedURLDuration = 30 * time.Minute
)

type Service interface {
	Stat(ctx context.Context, bucket, path string) (Object, error)
	Remove(ctx context.Context, obj Object) error
	GetPresignedURL(ctx context.Context, bucket, path, filename string) (string, error)
	PutPresignedURL(ctx context.Context, bucket, path, filename string) (string, error)
}

type service struct {
	host      string
	apikey    string
	secretkey string

	mc *minio.Client
}

func NewDefaultMinioService() (Service, error) {
	return NewMinioService(
		os.Getenv(envHost),
		os.Getenv(envApiKey),
		os.Getenv(envSecretKey),
	)
}

func NewMinioService(host, apikey, secretkey string) (Service, error) {
	mc, err := minio.New(host, &minio.Options{
		Creds: credentials.NewStaticV4(apikey, secretkey, ""),
	})
	if err != nil {
		return nil, err
	}

	return &service{host, apikey, secretkey, mc}, nil
}

func (svc service) Stat(ctx context.Context, bucket, path string) (Object, error) {
	info, err := svc.mc.StatObject(ctx, bucket, path, minio.StatObjectOptions{})
	if err != nil {
		return Object{}, err
	}
	obj := ObjectFromObjectInfo(info)
	obj.Bucket = bucket
	return obj, nil
}

func (svc service) Remove(ctx context.Context, obj Object) error {
	return svc.mc.RemoveObject(ctx, obj.Bucket, obj.Path, minio.RemoveObjectOptions{})
}

func (svc service) GetPresignedURL(ctx context.Context, bucket, path, filename string) (string, error) {
	// Set request parameters for content-disposition.
	reqParams := make(url.Values)
	reqParams.Set("response-content-disposition", fmt.Sprintf("attachment; filename=\"%s\"", filename))

	// Generates a presigned url which expires in a day.
	presignedURL, err := svc.mc.PresignedGetObject(ctx, bucket, path, defaultPresignedURLDuration, reqParams)
	return presignedURL.String(), err
}

func (svc service) PutPresignedURL(ctx context.Context, bucket, path, filename string) (string, error) {
	presignedURL, err := svc.mc.PresignedPutObject(
		ctx,
		bucket,
		path,
		defaultPresignedURLDuration)
	return presignedURL.String(), err
}
