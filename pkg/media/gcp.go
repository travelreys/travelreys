package media

import (
	"context"
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"
)

var (
	projectID          = os.Getenv("TRAVELREYS_GCP_PROJECT")
	gcpCloudCDNKeyName = os.Getenv("TRAVELREYS_GCP_CLOUD_CDN_KEY_NAME")
	gcpCloudCDNKeyPath = os.Getenv("TRAVELREYS_GCP_CLOUD_CDN_KEY_PATH")
)

type gcpProvider struct {
	projectID string
	keyName   string
	keyFile   string
}

func NewDefaultGCPCloudCDNProvider() CDNProvider {
	return gcpProvider{projectID, gcpCloudCDNKeyName, gcpCloudCDNKeyPath}
}

func (gcp gcpProvider) Domain(ctx context.Context, withScheme bool) string {
	domain := os.Getenv("TRAVELREYS_MEDIA_DOMAIN")
	if withScheme {
		return "https://" + domain // cdn.travelreys.com
	}
	return domain
}

// readKeyFile reads the base64url-encoded key file and decodes it.
func (gcp gcpProvider) readKeyFile() ([]byte, error) {
	b, err := ioutil.ReadFile(gcp.keyFile)
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

func (gcp gcpProvider) PresignedURL(ctx context.Context, url string) (string, error) {
	sep := "?"
	if strings.Contains(url, "?") {
		sep = "&"
	}
	url += sep
	url += fmt.Sprintf("Expires=%d", time.Now().Add(defaultPresignedURLDuration).Unix())
	url += fmt.Sprintf("&KeyName=%s", gcp.keyName)

	key, err := gcp.readKeyFile()
	if err != nil {
		return "", nil
	}
	mac := hmac.New(sha1.New, key)
	mac.Write([]byte(url))
	sig := base64.URLEncoding.EncodeToString(mac.Sum(nil))
	url += fmt.Sprintf("&Signature=%s", sig)
	return url, nil
}

// signCookie creates a signed cookie for an endpoint served by Cloud CDN.
// - urlPrefix must start with "https://" and should include the path prefix
// for which the cookie will authorize access to.
// - key should be in raw form (not base64url-encoded) which is
// 16-bytes long.
// - keyName must match a key added to the backend service or bucket.
func (gcp gcpProvider) signCookie(urlPrefix, keyName string, key []byte, expiration time.Time) (string, error) {
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

func (gcp gcpProvider) PresignedCookie(ctx context.Context, domain, path string) (*http.Cookie, error) {
	// Note: consider using the GCP Secret Manager for managing access to your
	// signing key(s).
	key, err := gcp.readKeyFile()
	if err != nil {
		return nil, err
	}

	expiration := defaultPresignedCookieDuration
	signedValue, err := gcp.signCookie(
		fmt.Sprintf("https://%s/%s/", domain, path),
		gcpCloudCDNKeyName,
		key,
		time.Now().Add(expiration),
	)
	if err != nil {
		return nil, err
	}

	// Use Go's http.Cookie type to construct a cookie.
	// domain and path should match the user-facing URL for accessing content.
	// Best practice: only send the cookie for paths it is valid for
	cookie := &http.Cookie{
		Name:     presignedCookieHeader,
		Value:    signedValue,
		Path:     "/",
		Domain:   os.Getenv("TRAVELREYS_MEDIA_COOKIE_DOMAIN"),
		MaxAge:   int(expiration.Seconds()),
		SameSite: http.SameSiteNoneMode,
		Secure:   true,
		HttpOnly: true,
	}
	return cookie, nil
}