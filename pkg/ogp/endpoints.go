package ogp

import (
	"context"

	"github.com/go-kit/kit/endpoint"
	"github.com/travelreys/travelreys/pkg/common"
)

type FetchRequest struct {
	Url string `json:"url"`
}
type FetchResponse struct {
	Graph Opengraph `json:"ogp"`
	Err   error     `json:"error,omitempty"`
}

func (r FetchResponse) Error() error {
	return r.Err
}

func NewFetchEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, epReq interface{}) (interface{}, error) {
		req, ok := epReq.(FetchRequest)
		if !ok {
			return FetchResponse{
				Err: common.ErrEndpointReqMismatch}, nil
		}
		g, err := svc.Fetch(ctx, req.Url)
		return FetchResponse{Graph: g, Err: err}, nil
	}
}
