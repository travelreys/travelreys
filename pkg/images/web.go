package images

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"go.uber.org/zap"
)

const (
	// https://unsplash.com/documentation
	webapiLoggerName = "images.webapi"
	unsplashUrl      = "https://api.unsplash.com"
)

var (
	ErrProviderUnsplashError = errors.New("images.webapi.unsplash.error")
)

type ImageSearchResponse struct {
	Total      uint64            `json:"total"`
	TotalPages uint64            `json:"total_pages"`
	Results    ImageMetadataList `json:"results"`
}

type WebImageAPI interface {
	Search(ctx context.Context, query string) (ImageMetadataList, error)
}

type unsplash struct {
	accesskey string
	logger    *zap.Logger
}

func NewWebImageAPI(accesskey string, logger *zap.Logger) WebImageAPI {
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

func (api unsplash) Search(ctx context.Context, query string) (ImageMetadataList, error) {
	if query == "" {
		return ImageMetadataList{}, ErrEmptySearchQuery
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
		api.logger.Error(
			"Search",
			zap.String("query", query),
			zap.Error(err),
		)
		return ImageMetadataList{}, ErrProviderUnsplashError
	}

	var res ImageSearchResponse
	err = json.Unmarshal(body, &res)
	if err != nil {
		api.logger.Error(
			"Search",
			zap.String("query", query),
			zap.Error(err),
		)
		return ImageMetadataList{}, ErrProviderUnsplashError
	}

	return res.Results, nil
}
