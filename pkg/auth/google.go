package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/url"
	"os"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

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

func NewGoogleOAuth2ConfigFromFile(file string) (GoogleOAuth2Config, error) {
	cfg := GoogleOAuth2Config{}
	jsonFile, err := os.Open(file)
	if err != nil {
		return cfg, err
	}
	defer jsonFile.Close()

	data, _ := ioutil.ReadAll(jsonFile)
	json.Unmarshal(data, &cfg)
	return cfg, nil
}

type GoogleProvider struct {
	cfg *oauth2.Config
}

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

func (gp *GoogleProvider) AuthCodeToToken(ctx context.Context, code string) (*oauth2.Token, error) {
	return gp.cfg.Exchange(ctx, code)
}

func (gp *GoogleProvider) TokenToUserInfo(ctx context.Context, token *oauth2.Token) (GoogleUser, error) {
	gusr := GoogleUser{}

	client := gp.cfg.Client(ctx, token)
	userInfoURLPrefix := "https://www.googleapis.com/oauth2/v2/userinfo?access_token=%s"
	resp, err := client.Get(fmt.Sprintf(userInfoURLPrefix, url.QueryEscape(token.AccessToken)))
	if err != nil {
		fmt.Println(err.Error())
		return gusr, err
	}
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return gusr, err
	}

	err = json.Unmarshal(data, &gusr)
	return gusr, err
}
