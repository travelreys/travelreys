package images

import (
	"context"
	"errors"
)

// Service

var (
	ErrEmptySearchQuery = errors.New("images.service.search.emptyquery")
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
		return ImageMetadataList{}, ErrEmptySearchQuery
	}
	return svc.api.Search(ctx, query)
}
