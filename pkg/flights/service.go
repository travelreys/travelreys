package flights

import (
	"context"
)

type Service interface {
	Search(ctx context.Context, opts SearchOptions) (Itineraries, error)
}

type service struct {
	flightAPI WebAPI
}

func NewService(flightAPI WebAPI) Service {
	return &service{flightAPI}
}

func (svc *service) Search(ctx context.Context, opts SearchOptions) (Itineraries, error) {
	return svc.flightAPI.Search(ctx, opts)
}
