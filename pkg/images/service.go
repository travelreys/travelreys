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
	Search(ctx context.Context, query string) (MetadataList, error)
}

type service struct {
	api WebAPI
}

func NewService(webAPI WebAPI) Service {
	return &service{webAPI}
}

func (svc *service) Search(ctx context.Context, query string) (MetadataList, error) {
	if query == "" {
		return MetadataList{}, ErrEmptySearchQuery
	}
	return svc.api.Search(ctx, query)
}
