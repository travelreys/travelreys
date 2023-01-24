package flights

import (
	"context"
	"time"
)

type Service interface {
	Search(ctx context.Context, origIATA, destIATA string, numAdults uint64, departDate time.Time, opts FlightsSearchOptions) (ItinerariesList, error)
}

type service struct {
	flightAPI WebFlightsAPI
}

func NewService(flightAPI WebFlightsAPI) Service {
	return &service{flightAPI}
}

func (svc *service) Search(
	ctx context.Context,
	origIATA,
	destIATA string,
	numAdults uint64,
	departDate time.Time,
	opts FlightsSearchOptions,
) (ItinerariesList, error) {
	return svc.flightAPI.Search(ctx, origIATA, destIATA, numAdults, departDate, opts)
}
