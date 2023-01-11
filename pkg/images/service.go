package images

import (
	"context"
	"errors"
)

// Service

var (
	ErrEmptyQuery = errors.New("empty-image-search-query")
)

type Service interface {
	Search(ctx context.Context, query string) (ImageMetadataList, error)
}

type service struct {
	api WebImageAPI
}

func NewService(webAPI WebImageAPI) Service {
	return &service{webAPI}
}

func (svc *service) Search(ctx context.Context, query string) (ImageMetadataList, error) {
	if query == "" {
		return ImageMetadataList{}, ErrEmptyQuery
	}
	return svc.api.Search(ctx, query)
}
