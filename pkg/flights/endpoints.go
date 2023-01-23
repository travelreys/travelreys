package flights

import (
	"context"
	"time"

	"github.com/go-kit/kit/endpoint"
	"github.com/tiinyplanet/tiinyplanet/pkg/common"
)

type SearchRequest struct {
	origIATA   string
	destIATA   string
	numAdults  uint64
	departDate time.Time
	opts       FlightsSearchOptions
}

type SearchResponse struct {
	Itineraries ItinerariesList `json:"itineraries"`
	Err         error           `json:"error,omitempty"`
}

func (r SearchResponse) Error() error {
	return r.Err
}

func NewSearchEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, epReq interface{}) (interface{}, error) {
		req, ok := epReq.(SearchRequest)
		if !ok {
			return SearchResponse{Err: common.ErrInvalidEndpointRequestType}, nil
		}
		itins, err := svc.Search(ctx, req.origIATA, req.destIATA, req.numAdults, req.departDate, req.opts)
		return SearchResponse{Itineraries: itins, Err: err}, nil
	}
}
