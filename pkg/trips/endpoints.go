package trips

import (
	context "context"
	"time"

	"github.com/go-kit/kit/endpoint"
	"github.com/tiinyplanet/tiinyplanet/pkg/common"
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
		creator := TripMember{
			MemberID:    "1",
			MemberEmail: "awhdes@gmail.com",
		}

		plan, err := svc.CreateTripPlan(ctx, creator, req.Name, req.StartDate, req.EndDate)
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
		plan, err := svc.ReadTripPlan(ctx, req.ID)
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
		plans, err := svc.ListTripPlans(ctx, req.FF)
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
		err := svc.DeleteTripPlan(ctx, req.ID)
		return DeleteTripPlanResponse{Err: err}, nil
	}
}
