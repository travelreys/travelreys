package storage

import (
	"context"
	"fmt"
	"os"
	"time"

	"cloud.google.com/go/storage"
	"google.golang.org/api/option"
)

var (
	projectID = os.Getenv("TRAVELREYS_GCS_PROJECT")
)

type gcsService struct {
	cl        *storage.Client
	projectID string
	credsPath string
}

func NewDefaultGCSService(ctx context.Context) (Service, error) {
	return NewGCSService(ctx, os.Getenv("TRAVELREYS_GCS_AUTH_JSON"))
}

func NewGCSService(ctx context.Context, credsPath string) (Service, error) {
	var (
		client *storage.Client
		err    error
	)
	if credsPath != "" {
		client, err = storage.NewClient(ctx, option.WithCredentialsFile(credsPath))
	} else {
		client, err = storage.NewClient(ctx)
	}
	if err != nil {
		return nil, err
	}

	return &gcsService{
		client,
		projectID,
		credsPath,
	}, nil
}
func (svc gcsService) Stat(ctx context.Context, bucket, path string) (Object, error) {
	attrs, err := svc.cl.Bucket(bucket).Object(path).Attrs(ctx)
	if err != nil {
		return Object{}, err
	}
	obj := ObjectFromAttrs(attrs)
	obj.Bucket = bucket
	return obj, nil
}

func (svc gcsService) Remove(ctx context.Context, obj Object) error {
	return svc.cl.Bucket(obj.Bucket).Object(obj.Path).Delete(ctx)
}

func (svc gcsService) GetPresignedURL(ctx context.Context, bucket, path, filename string) (string, error) {
	// https://cloud.google.com/storage/docs/access-control/signing-urls-with-helpers
	// Signing a URL requires credentials authorized to sign a URL. You can pass
	// these in through SignedURLOptions with one of the following options:
	//    a. a Google service account private key, obtainable from the Google Developers Console
	//    b. a Google Access ID with iam.serviceAccounts.signBlob permissions
	//    c. a SignBytes function implementing custom signing.
	// In this example, none of these options are used, which means the SignedURL
	// function attempts to use the same authentication that was used to instantiate
	// the Storage client. This authentication must include a private key or have
	// iam.serviceAccounts.signBlob permissions.
	opts := &storage.SignedURLOptions{
		Scheme:  storage.SigningSchemeV4,
		Method:  "GET",
		Headers: []string{fmt.Sprintf("response-content-disposition:attachment; filename=\"%s\"", filename)},
		Expires: time.Now().Add(15 * time.Minute),
	}

	u, err := svc.cl.Bucket(bucket).SignedURL(path, opts)
	if err != nil {
		return "", fmt.Errorf("Bucket(%q).SignedURL: %v", bucket, err)
	}
	return u, nil
}

func (svc gcsService) PutPresignedURL(ctx context.Context, bucket, path, filename string) (string, error) {
	// https://cloud.google.com/storage/docs/access-control/signing-urls-with-helpers
	// Signing a URL requires credentials authorized to sign a URL. You can pass
	// these in through SignedURLOptions with one of the following options:
	//    a. a Google service account private key, obtainable from the Google Developers Console
	//    b. a Google Access ID with iam.serviceAccounts.signBlob permissions
	//    c. a SignBytes function implementing custom signing.
	// In this example, none of these options are used, which means the SignedURL
	// function attempts to use the same authentication that was used to instantiate
	// the Storage client. This authentication must include a private key or have
	// iam.serviceAccounts.signBlob permissions.
	opts := &storage.SignedURLOptions{
		Scheme:  storage.SigningSchemeV4,
		Method:  "PUT",
		Expires: time.Now().Add(15 * time.Minute),
	}

	u, err := svc.cl.Bucket(bucket).SignedURL(path, opts)
	if err != nil {
		return "", fmt.Errorf("Bucket(%q).SignedURL: %v", bucket, err)
	}
	return u, nil
}
