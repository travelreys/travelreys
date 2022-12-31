package trips

import (
	context "context"
	"time"

	"github.com/awhdesmond/tiinyplanet/pkg/common"
	"github.com/awhdesmond/tiinyplanet/pkg/reqctx"
	"github.com/go-kit/kit/endpoint"
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
			return CreateTripPlanResponse{Err: common.ErrInvalidEndpointRequestType}, nil
		}
		rctx, ok := ctx.(reqctx.Context)
		if !ok {
			return CreateTripPlanResponse{Err: common.ErrInvalidEndpointRequestContext}, nil
		}
		creator := TripMember{
			MemberID:    rctx.CallerInfo.UserID,
			MemberEmail: rctx.CallerInfo.UserEmail,
			Permission:  TripMemberPermCollaborator,
		}

		plan, err := svc.CreateTripPlan(rctx, creator, req.Name, req.StartDate, req.EndDate)
		return CreateTripPlanResponse{Plan: plan, Err: err}, nil
	}
}

type ReadTripPlanRequest struct {
	ID string `json:"id"`
}

type ReadTripPlanResponse struct {
	Plan TripPlan `json:"plan"`
	Err  error    `json:"error,omitempty"`
}

func (r ReadTripPlanResponse) Error() error {
	return r.Err
}

func NewReadTripPlanEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, epReq interface{}) (interface{}, error) {
		req, ok := epReq.(ReadTripPlanRequest)
		if !ok {
			return ReadTripPlanResponse{Err: common.ErrInvalidEndpointRequestType}, nil
		}
		rctx, ok := ctx.(reqctx.Context)
		if !ok {
			return CreateTripPlanResponse{Err: common.ErrInvalidEndpointRequestContext}, nil
		}

		plan, err := svc.ReadTripPlan(rctx, req.ID)
		return ReadTripPlanResponse{Plan: plan, Err: err}, nil
	}
}

type ListTripPlansRequest struct {
	FF ListTripPlansFilter
}
type ListTripPlansResponse struct {
	Plans TripPlansList `json:"plans"`
	Err   error         `json:"error"`
}

func (r ListTripPlansResponse) Error() error {
	return r.Err
}

func NewListTripPlansEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, epReq interface{}) (interface{}, error) {
		req, ok := epReq.(ListTripPlansRequest)
		if !ok {
			return ListTripPlansResponse{Err: common.ErrInvalidEndpointRequestType}, nil
		}
		rctx, ok := ctx.(reqctx.Context)
		if !ok {
			return CreateTripPlanResponse{Err: common.ErrInvalidEndpointRequestContext}, nil
		}

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
			return DeleteTripPlanResponse{Err: common.ErrInvalidEndpointRequestType}, nil
		}
		rctx, ok := ctx.(reqctx.Context)
		if !ok {
			return CreateTripPlanResponse{Err: common.ErrInvalidEndpointRequestContext}, nil
		}

		err := svc.DeleteTripPlan(rctx, req.ID)
		return DeleteTripPlanResponse{Err: err}, nil
	}
}
