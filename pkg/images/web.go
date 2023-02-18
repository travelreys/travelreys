package images

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

const (
	// https://unsplash.com/documentation
	UNSPLASH_URL = "https://api.unsplash.com"
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
}

func NewWebImageAPI(accesskey string) WebImageAPI {
	return unsplash{accesskey}
}

func (api unsplash) fullUrlPath(path string) string {
	return fmt.Sprintf("%s/%s", UNSPLASH_URL, path)
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
		return ImageMetadataList{}, ErrEmptyQuery
	}

	queryURL, _ := url.Parse(api.fullUrlPath("/search/photos"))
	queryParams := queryURL.Query()
	queryParams.Set("query", query)
	queryParams.Set("per_page", "30")
	queryParams.Set("orientation", "landscape")
	queryParams.Set("client_id", api.accesskey)
	queryURL.RawQuery = queryParams.Encode()

	body, err := api.Get(queryURL)
	if err != nil {
		return ImageMetadataList{}, err
	}

	var res ImageSearchResponse
	err = json.Unmarshal(body, &res)
	if err != nil {
		return ImageMetadataList{}, err
	}

	return res.Results, nil
}
