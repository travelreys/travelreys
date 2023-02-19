package trips

import (
	context "context"
	"time"

	"github.com/go-kit/kit/endpoint"
	"github.com/tiinyplanet/tiinyplanet/pkg/auth"
	"github.com/tiinyplanet/tiinyplanet/pkg/common"
)

// Trips Endpoints

type CreateTripRequest struct {
	Name      string    `json:"name"`
	StartDate time.Time `json:"startDate"`
	EndDate   time.Time `json:"endDate"`
}
type CreateTripResponse struct {
	Plan TripPlan `json:"tripPlan"`
	Err  error    `json:"error,omitempty"`
}

func (r CreateTripResponse) Error() error {
	return r.Err
}

func NewCreateTripEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, epReq interface{}) (interface{}, error) {
		req, ok := epReq.(CreateTripRequest)
		if !ok {
			return CreateTripResponse{
				Err: common.ErrorInvalidEndpointRequestType}, nil
		}

		ci, err := common.ReadClientInfoFromCtx(ctx)
		if err != nil {
			return CreateTripResponse{
				Plan: TripPlan{},
				Err:  ErrRBACMissing,
			}, nil
		}

		creator := NewMember(ci.UserID, MemberRoleCreator)
		plan, err := svc.CreateTrip(ctx, creator, req.Name, req.StartDate, req.EndDate)
		return CreateTripResponse{Plan: plan, Err: err}, nil
	}
}

type ReadTripRequest struct {
	ID        string `json:"id"`
	WithUsers bool   `json:"withUsers"`
}

type ReadTripResponse struct {
	Plan TripPlan `json:"tripPlan"`
	Err  error    `json:"error,omitempty"`
}

func (r ReadTripResponse) Error() error {
	return r.Err
}

type ReadTripWithUsersResponse struct {
	Plan  TripPlan      `json:"tripPlan"`
	Users auth.UsersMap `json:"users"`
	Err   error         `json:"error,omitempty"`
}

func (r ReadTripWithUsersResponse) Error() error {
	return r.Err
}

func NewReadTripEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, epReq interface{}) (interface{}, error) {
		req, ok := epReq.(ReadTripRequest)
		if !ok {
			return ReadTripResponse{Err: common.ErrorInvalidEndpointRequestType}, nil
		}

		if req.WithUsers {
			plan, users, err := svc.ReadTripWithUsers(ctx, req.ID)
			return ReadTripWithUsersResponse{Plan: plan, Users: users, Err: err}, nil
		}

		plan, err := svc.ReadTrip(ctx, req.ID)
		return ReadTripResponse{Plan: plan, Err: err}, nil
	}
}

type ListTripsRequest struct {
	FF ListTripsFilter
}
type ListTripsResponse struct {
	Plans TripPlansList `json:"tripPlans"`
	Err   error         `json:"error,omitempty"`
}

func (r ListTripsResponse) Error() error {
	return r.Err
}

func NewListTripsEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, epReq interface{}) (interface{}, error) {
		req, ok := epReq.(ListTripsRequest)
		if !ok {
			return ListTripsResponse{Err: common.ErrorInvalidEndpointRequestType}, nil
		}
		plans, err := svc.ListTrips(ctx, req.FF)
		return ListTripsResponse{Plans: plans, Err: err}, nil
	}
}

type DeleteTripRequest struct {
	ID string `json:"id"`
}

type DeleteTripResponse struct {
	Err error `json:"error,omitempty"`
}

func (r DeleteTripResponse) Error() error {
	return r.Err
}

func NewDeleteTripEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, epReq interface{}) (interface{}, error) {
		req, ok := epReq.(DeleteTripRequest)
		if !ok {
			return DeleteTripResponse{Err: common.ErrorInvalidEndpointRequestType}, nil
		}
		err := svc.DeleteTrip(ctx, req.ID)
		return DeleteTripResponse{Err: err}, nil
	}
}
