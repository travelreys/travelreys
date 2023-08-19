package finance

import (
	"context"

	"github.com/go-kit/kit/endpoint"
	"github.com/travelreys/travelreys/pkg/common"
)

type GetFxRatesRequest struct {
	Base string
}

type GetFxRatesResponse struct {
	Rates ExchangeRates `json:"rates"`
	Err   error         `json:"error,omitempty"`
}

func (r GetFxRatesResponse) Error() error {
	return r.Err
}

func NewGetFxRatesEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, epReq interface{}) (interface{}, error) {
		req, ok := epReq.(GetFxRatesRequest)
		if !ok {
			return GetFxRatesResponse{Err: common.ErrEndpointReqMismatch}, nil
		}
		rates, err := svc.GetFxRates(ctx, req.Base)
		return GetFxRatesResponse{Rates: rates, Err: err}, nil
	}
}
