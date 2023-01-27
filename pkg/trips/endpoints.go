package trips

import (
	context "context"
	"time"

	"github.com/go-kit/kit/endpoint"
	"github.com/tiinyplanet/tiinyplanet/pkg/common"
	"github.com/tiinyplanet/tiinyplanet/pkg/reqctx"
)

// Trips Endpoints

type CreateTripPlanRequest struct {
	Name      string    `json:"name"`
	StartDate time.Time `json:"startDate"`
	EndDate   time.Time `json:"endDate"`
}
type CreateTripPlanResponse struct {
	Plan TripPlan `json:"tripPlan"`
	Err  error    `json:"error,omitempty"`
}

func (r CreateTripPlanResponse) Error() error {
	return r.Err
}

func NewCreateTripPlanEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, epReq interface{}) (interface{}, error) {
		req, ok := epReq.(CreateTripPlanRequest)
		if !ok {
			return CreateTripPlanResponse{Err: common.ErrorInvalidEndpointRequestType}, nil
		}
		rctx := reqctx.Context{Context: ctx, CallerInfo: reqctx.CallerInfo{}}
		creator := TripMember{
			MemberID:    "1",
			MemberEmail: "awhdes@gmail.com",
		}

		plan, err := svc.CreateTripPlan(rctx, creator, req.Name, req.StartDate, req.EndDate)
		return CreateTripPlanResponse{Plan: plan, Err: err}, nil
	}
}

type ReadTripPlanRequest struct {
	ID string `json:"id"`
}

type ReadTripPlanResponse struct {
	Plan TripPlan `json:"tripPlan"`
	Err  error    `json:"error,omitempty"`
}

func (r ReadTripPlanResponse) Error() error {
	return r.Err
}

func NewReadTripPlanEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, epReq interface{}) (interface{}, error) {
		req, ok := epReq.(ReadTripPlanRequest)
		if !ok {
			return ReadTripPlanResponse{Err: common.ErrorInvalidEndpointRequestType}, nil
		}
		rctx := reqctx.Context{Context: ctx, CallerInfo: reqctx.CallerInfo{}}
		plan, err := svc.ReadTripPlan(rctx, req.ID)
		return ReadTripPlanResponse{Plan: plan, Err: err}, nil
	}
}

type ListTripPlansRequest struct {
	FF ListTripPlansFilter
}
type ListTripPlansResponse struct {
	Plans TripPlansList `json:"tripPlans"`
	Err   error         `json:"error,omitempty"`
}

func (r ListTripPlansResponse) Error() error {
	return r.Err
}

func NewListTripPlansEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, epReq interface{}) (interface{}, error) {
		req, ok := epReq.(ListTripPlansRequest)
		if !ok {
			return ListTripPlansResponse{Err: common.ErrorInvalidEndpointRequestType}, nil
		}
		rctx := reqctx.Context{Context: ctx, CallerInfo: reqctx.CallerInfo{}}
		plans, err := svc.ListTripPlans(rctx, req.FF)
		return ListTripPlansResponse{Plans: plans, Err: err}, nil
	}
}

type DeleteTripPlanRequest struct {
	ID string `json:"id"`
}

type DeleteTripPlanResponse struct {
	Err error `json:"error,omitempty"`
}

func (r DeleteTripPlanResponse) Error() error {
	return r.Err
}

func NewDeleteTripPlanEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, epReq interface{}) (interface{}, error) {
		req, ok := epReq.(DeleteTripPlanRequest)
		if !ok {
			return DeleteTripPlanResponse{Err: common.ErrorInvalidEndpointRequestType}, nil
		}
		rctx := reqctx.Context{Context: ctx, CallerInfo: reqctx.CallerInfo{}}
		err := svc.DeleteTripPlan(rctx, req.ID)
		return DeleteTripPlanResponse{Err: err}, nil
	}
}
