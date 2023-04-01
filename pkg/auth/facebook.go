package auth

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
)

const (
	fbGraphUserURL = "https://graph.facebook.com/me"
)

type FacebookProvider struct{}

func NewFacebookProvider() FacebookProvider {
	return FacebookProvider{}
}

func (fb FacebookProvider) Get(getURL *url.URL) ([]byte, error) {
	req, err := http.NewRequest(http.MethodGet, getURL.String(), nil)
	if err != nil {
		return nil, err
	}
	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	body, err := io.ReadAll(resp.Body)
	return body, err
}

func (fb FacebookProvider) TokenToUserInfo(ctx context.Context, code string) (FacebookUser, error) {
	queryURL, _ := url.Parse(fbGraphUserURL)
	queryParams := queryURL.Query()
	queryParams.Set("fields", "id,name,email,picture")
	queryParams.Set("access_token", code)
	queryURL.RawQuery = queryParams.Encode()

	body, err := fb.Get(queryURL)
	if err != nil {
		return FacebookUser{}, err
	}

	fbUsr := FacebookUser{}
	err = json.Unmarshal(body, &fbUsr)
	return fbUsr, err
}
