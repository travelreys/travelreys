package storage

import (
	"context"
	"fmt"
	"net/url"
	"os"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type minioService struct {
	host      string
	apikey    string
	secretkey string

	mc *minio.Client
}

func NewDefaultMinioService() (Service, error) {
	return NewMinioService(
		os.Getenv("TRAVELREYS_MINIO_HOST"),
		os.Getenv("TRAVELREYS_MINIO_APIKEY"),
		os.Getenv("TRAVELREYS_MINIO_SECRETKEY"),
	)
}

func NewMinioService(host, apikey, secretkey string) (Service, error) {
	mc, err := minio.New(host, &minio.Options{
		Creds: credentials.NewStaticV4(apikey, secretkey, ""),
	})
	if err != nil {
		return nil, err
	}
	return &minioService{host, apikey, secretkey, mc}, nil
}

func (svc minioService) Stat(ctx context.Context, bucket, path string) (Object, error) {
	info, err := svc.mc.StatObject(ctx, bucket, path, minio.StatObjectOptions{})
	if err != nil {
		return Object{}, err
	}
	obj := ObjectFromObjectInfo(info)
	obj.Bucket = bucket
	return obj, nil
}

func (svc minioService) Remove(ctx context.Context, obj Object) error {
	return svc.mc.RemoveObject(ctx, obj.Bucket, obj.Path, minio.RemoveObjectOptions{})
}

func (svc minioService) GetPresignedURL(ctx context.Context, bucket, path, filename string) (string, error) {
	// Set request parameters for content-disposition.
	reqParams := make(url.Values)
	reqParams.Set("response-content-disposition", fmt.Sprintf("attachment; filename=\"%s\"", filename))

	// Generates a presigned url which expires in a day.
	presignedURL, err := svc.mc.PresignedGetObject(ctx, bucket, path, defaultPresignedURLDuration, reqParams)
	if err != nil {
		return "", err
	}
	return presignedURL.String(), nil
}

func (svc minioService) PutPresignedURL(ctx context.Context, bucket, path, filename string) (string, error) {
	presignedURL, err := svc.mc.PresignedPutObject(
		ctx,
		bucket,
		path,
		defaultPresignedURLDuration)
	if err != nil {
		return "", err
	}
	return presignedURL.String(), err
}
