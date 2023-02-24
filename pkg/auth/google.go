package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"io"

	"net/url"
	"os"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

const (
	envOAuthGoolgeSecretFile = "TIINYPLANET_OAUTH_GOOGLE_SECRET_FILE"
	userInfoURL              = "https://www.googleapis.com/oauth2/v2/userinfo?access_token"
)

// GetOAuthGoogleSecretFile retrieves the file path to Google OAuth2 secrets
func GetOAuthGoogleSecretFile() string {
	return os.Getenv(envOAuthGoolgeSecretFile)
}

// GoogleOAuth2Config contains OAuth2 secrets for Google SSO
type GoogleOAuth2Config struct {
	Web struct {
		ClientID                string   `json:"client_id"`
		ProjectID               string   `json:"project_id"`
		AuthURI                 string   `json:"auth_uri"`
		TokenURI                string   `json:"token_uri"`
		AuthProviderX509CertURL string   `json:"auth_provider_x509_cert_url"`
		ClientSecret            string   `json:"client_secret"`
		JavascriptOrigins       []string `json:"javascript_origins"`
	} `json:"web"`
}

// NewGoogleOAuth2ConfigFromFile parses the secret file
// and generates the GoogleOAuth2Config
func NewGoogleOAuth2ConfigFromFile(file string) (GoogleOAuth2Config, error) {
	cfg := GoogleOAuth2Config{}
	jsonFile, err := os.Open(file)
	if err != nil {
		return cfg, err
	}
	defer jsonFile.Close()

	data, _ := io.ReadAll(jsonFile)
	json.Unmarshal(data, &cfg)
	return cfg, nil
}

type GoogleProvider struct {
	cfg *oauth2.Config
}

func NewDefaultGoogleProvider() (GoogleProvider, error) {
	return NewGoogleProvider(GetOAuthGoogleSecretFile())
}

// NewGoogleProvider returns a new Google OAuth2 provider
func NewGoogleProvider(cfgFile string) (GoogleProvider, error) {
	cfg, err := NewGoogleOAuth2ConfigFromFile(cfgFile)
	if err != nil {
		return GoogleProvider{}, err
	}
	return GoogleProvider{
		cfg: &oauth2.Config{
			ClientID:     cfg.Web.ClientID,
			ClientSecret: cfg.Web.ClientSecret,
			Scopes: []string{
				"profile",
				"email",
				"openid",
				"https://www.googleapis.com/auth/userinfo.profile",
				"https://www.googleapis.com/auth/userinfo.email",
			},
			// Note: https://github.com/MomenSherif/react-oauth/issues/12
			RedirectURL: "postmessage",
			Endpoint:    google.Endpoint,
		},
	}, nil
}

func (gp *GoogleProvider) TokenToUserInfo(ctx context.Context, code string) (GoogleUser, error) {
	token, err := gp.cfg.Exchange(ctx, code)
	if err != nil {
		return GoogleUser{}, err
	}

	client := gp.cfg.Client(ctx, token)
	resp, err := client.Get(fmt.Sprintf("%s=%s", userInfoURL, url.QueryEscape(token.AccessToken)))
	if err != nil {
		return GoogleUser{}, err
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return GoogleUser{}, err
	}
	gusr := GoogleUser{}
	err = json.Unmarshal(data, &gusr)
	return gusr, err
}
