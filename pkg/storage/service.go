package storage

import (
	"context"
	"io"
	"os"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

// https://medium.com/google-cloud/using-google-cloud-storage-with-minio-object-storage-c994fe4aab6b

const (
	envHost      = "TRAVELREYS_STORAGE_HOST"
	envApiKey    = "TRAVELREYS_STORAGE_APIKEY"
	envSecretKey = "TRAVELREYS_STORAGE_SECRETKEY"
)

type Service interface {
	Read(ctx context.Context, bucket, path string) (Object, error)
	Upload(ctx context.Context, obj Object, file io.Reader) error
	Download(ctx context.Context, obj Object) (io.ReadCloser, error)
	Remove(ctx context.Context, obj Object) error
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

func (svc service) Read(ctx context.Context, bucket, path string) (Object, error) {
	info, err := svc.mc.StatObject(ctx, bucket, path, minio.StatObjectOptions{})
	if err != nil {
		return Object{}, err
	}
	obj := ObjectFromObjectInfo(info)
	obj.Bucket = bucket
	return obj, nil

}

func (svc service) Upload(ctx context.Context, obj Object, file io.Reader) error {
	_, err := svc.mc.PutObject(ctx, obj.Bucket, obj.Path, file, obj.Size, minio.PutObjectOptions{
		ContentType: obj.MIMEType,
	})
	return err
}

func (svc service) Download(ctx context.Context, obj Object) (io.ReadCloser, error) {
	return svc.mc.GetObject(ctx, obj.Bucket, obj.Path, minio.GetObjectOptions{})
}

func (svc service) Remove(ctx context.Context, obj Object) error {
	return svc.mc.RemoveObject(ctx, obj.Bucket, obj.Path, minio.RemoveObjectOptions{})
}
