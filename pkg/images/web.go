package images

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"

	"go.uber.org/zap"
)

const (
	// https://unsplash.com/documentation
	webapiLoggerName     = "images.webapi"
	unsplashUrl          = "https://api.unsplash.com"
	envUnsplashAccessKey = "TIINYPLANET_UNSPLASH_ACCESSKEY"
)

var (
	ErrProviderUnsplashError = errors.New("images.webapi.unsplash.error")
)

type WebAPI interface {
	Search(context.Context, string) (MetadataList, error)
}

type ImageSearchResponse struct {
	Total      uint64       `json:"total"`
	TotalPages uint64       `json:"total_pages"`
	Results    MetadataList `json:"results"`
}

type unsplash struct {
	accesskey string
	logger    *zap.Logger
}

func GetApiAccessKey() string {
	return os.Getenv(envUnsplashAccessKey)
}

func NewDefaultWebAPI(logger *zap.Logger) WebAPI {
	return NewWebAPI(GetApiAccessKey(), logger)
}

func NewWebAPI(accesskey string, logger *zap.Logger) WebAPI {
	return unsplash{accesskey, logger.Named(webapiLoggerName)}
}

func (api unsplash) fullUrl(path string) string {
	return fmt.Sprintf("%s/%s", unsplashUrl, path)
}

func (api unsplash) Get(getURL *url.URL) ([]byte, error) {
	req, err := http.NewRequest(http.MethodGet, getURL.String(), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Accept-Version", "v1")

	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	body, err := io.ReadAll(resp.Body)
	return body, err
}

func (api unsplash) Search(ctx context.Context, query string) (MetadataList, error) {
	if query == "" {
		return MetadataList{}, ErrEmptySearchQuery
	}

	queryURL, _ := url.Parse(api.fullUrl("/search/photos"))
	queryParams := queryURL.Query()
	queryParams.Set("query", query)
	queryParams.Set("per_page", "30")
	queryParams.Set("orientation", "landscape")
	queryParams.Set("client_id", api.accesskey)
	queryURL.RawQuery = queryParams.Encode()

	body, err := api.Get(queryURL)
	if err != nil {
		api.logger.Error("Search", zap.String("query", query), zap.Error(err))
		return MetadataList{}, ErrProviderUnsplashError
	}

	var res ImageSearchResponse
	err = json.Unmarshal(body, &res)
	if err != nil {
		api.logger.Error("Search", zap.String("query", query), zap.Error(err))
		return MetadataList{}, ErrProviderUnsplashError
	}

	return res.Results, nil
}
