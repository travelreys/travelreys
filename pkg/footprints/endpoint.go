package footprints

import (
	"context"

	"github.com/go-kit/kit/endpoint"
	"github.com/travelreys/travelreys/pkg/common"
	"github.com/travelreys/travelreys/pkg/trips"
)

type CheckInRequest struct {
	UserID   string         `json:"userID"`
	TripID   string         `json:"tripID"`
	Activity trips.Activity `json:"activity"`
}

type CheckInResponse struct {
	Err error `json:"error,omitepty"`
}

func (r CheckInResponse) Error() error {
	return r.Err
}

func NewCheckInEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, epReq interface{}) (interface{}, error) {
		req, ok := epReq.(CheckInRequest)
		if !ok {
			return CheckInResponse{Err: common.ErrorEndpointReqMismatch}, nil
		}
		err := svc.CheckIn(ctx, req.UserID, req.TripID, req.Activity)
		return CheckInResponse{Err: err}, nil
	}
}

type ListRequest struct {
	ListFootprintsFilter
}

type ListResponse struct {
	Footprints FootprintList `json:"footprints"`
	Err        error         `json:"error,omitepty"`
}

func (r ListResponse) Error() error {
	return r.Err
}

func NewListEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, epReq interface{}) (interface{}, error) {
		req, ok := epReq.(ListRequest)
		if !ok {
			return ListResponse{Err: common.ErrorEndpointReqMismatch}, nil
		}
		items, err := svc.List(ctx, req.ListFootprintsFilter)
		return ListResponse{items, err}, nil
	}
}
