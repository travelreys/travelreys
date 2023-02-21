package flights

import (
	"context"

	"github.com/go-kit/kit/endpoint"
	"github.com/tiinyplanet/tiinyplanet/pkg/common"
)

type SearchRequest struct {
	opts SearchOptions
}

type SearchResponse struct {
	Itineraries Itineraries `json:"itineraries"`
	Err         error       `json:"error,omitempty"`
}

func (r SearchResponse) Error() error {
	return r.Err
}

func NewSearchEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, epReq interface{}) (interface{}, error) {
		req, ok := epReq.(SearchRequest)
		if !ok {
			return SearchResponse{Err: common.ErrorMismatchEndpointReq}, nil
		}
		itins, err := svc.Search(ctx, req.opts)
		return SearchResponse{Itineraries: itins, Err: err}, nil
	}
}
