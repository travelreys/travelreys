package images

import (
	"context"

	"github.com/awhdesmond/tiinyplanet/pkg/common"
	"github.com/go-kit/kit/endpoint"
)

type SearchRequest struct {
	Query string
}

type SearchResponse struct {
	Images ImageMetadataList `json:"images"`
	Err    error             `json:"error,omitempty"`
}

func (r SearchResponse) Error() error {
	return r.Err
}

func NewCreateTripPlanEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, epReq interface{}) (interface{}, error) {
		req, ok := epReq.(SearchRequest)
		if !ok {
			return SearchResponse{Err: common.ErrInvalidEndpointRequestType}, nil
		}
		images, err := svc.Search(ctx, req.Query)
		return SearchResponse{Images: images, Err: err}, nil
	}
}
