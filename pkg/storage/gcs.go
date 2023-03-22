package storage

import (
	"context"
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"cloud.google.com/go/storage"
	"google.golang.org/api/option"
)

var (
	projectID     = os.Getenv("TRAVELREYS_GCS_PROJECT")
	cookieKeyName = os.Getenv("TRAVELREYS_GCS_COOKIE_KEY_NAME")
	cookieKeyPath = os.Getenv("TRAVELREYS_GCS_COOKIE_KEY_PATH")
)

const (
	cookieHeader = "Cloud-CDN-Cookie"
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

// signCookie creates a signed cookie for an endpoint served by Cloud CDN.
// - urlPrefix must start with "https://" and should include the path prefix
// for which the cookie will authorize access to.
// - key should be in raw form (not base64url-encoded) which is
// 16-bytes long.
// - keyName must match a key added to the backend service or bucket.
func signCookie(urlPrefix, keyName string, key []byte, expiration time.Time) (string, error) {
	encodedURLPrefix := base64.URLEncoding.EncodeToString([]byte(urlPrefix))
	input := fmt.Sprintf(
		"URLPrefix=%s:Expires=%d:KeyName=%s",
		encodedURLPrefix,
		expiration.Unix(),
		keyName)

	mac := hmac.New(sha1.New, key)
	mac.Write([]byte(input))
	sig := base64.URLEncoding.EncodeToString(mac.Sum(nil))

	signedValue := fmt.Sprintf("%s:Signature=%s", input, sig)
	return signedValue, nil
}

// readKeyFile reads the base64url-encoded key file and decodes it.
func readKeyFile(path string) ([]byte, error) {
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read key file: %+v", err)
	}
	d := make([]byte, base64.URLEncoding.DecodedLen(len(b)))
	n, err := base64.URLEncoding.Decode(d, b)
	if err != nil {
		return nil, fmt.Errorf("failed to base64url decode: %+v", err)
	}
	return d[:n], nil
}

func (svc gcsService) GeneratePresignedCookie(ctx context.Context, domain, path string) (*http.Cookie, error) {
	// Note: consider using the GCP Secret Manager for managing access to your
	// signing key(s).
	key, err := readKeyFile(cookieKeyPath)
	if err != nil {
		return nil, err
	}

	expiration := defaultPresignedCookieDuration
	signedValue, err := signCookie(
		fmt.Sprintf("https://%s%s", domain, path),
		cookieKeyName,
		key,
		time.Now().Add(expiration),
	)
	if err != nil {
		return nil, err
	}

	// Use Go's http.Cookie type to construct a cookie.
	// domain and path should match the user-facing URL for accessing content.
	cookie := &http.Cookie{
		Name:   cookieHeader,
		Value:  signedValue,
		Path:   path, // Best practice: only send the cookie for paths it is valid for
		Domain: domain,
		MaxAge: int(expiration.Seconds()),
	}

	// In a real application, use the SetCookie method on a http.ResponseWriter
	// to write the cookie to the user.
	return cookie, nil
}
