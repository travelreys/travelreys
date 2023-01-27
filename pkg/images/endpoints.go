package images

import (
	"context"

	"github.com/go-kit/kit/endpoint"
	"github.com/tiinyplanet/tiinyplanet/pkg/common"
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

func NewSearchEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, epReq interface{}) (interface{}, error) {
		req, ok := epReq.(SearchRequest)
		if !ok {
			return SearchResponse{Err: common.ErrorInvalidEndpointRequestType}, nil
		}
		images, err := svc.Search(ctx, req.Query)
		return SearchResponse{Images: images, Err: err}, nil
	}
}
