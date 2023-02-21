package trips

import (
	context "context"
	"time"

	"github.com/go-kit/kit/endpoint"
	"github.com/tiinyplanet/tiinyplanet/pkg/auth"
	"github.com/tiinyplanet/tiinyplanet/pkg/common"
	"github.com/tiinyplanet/tiinyplanet/pkg/reqctx"
)

// Trips Endpoints

type CreateRequest struct {
	Name      string    `json:"name"`
	StartDate time.Time `json:"startDate"`
	EndDate   time.Time `json:"endDate"`
}
type CreateResponse struct {
	Plan Trip  `json:"trip"`
	Err  error `json:"error,omitempty"`
}

func (r CreateResponse) Error() error {
	return r.Err
}

func NewCreateEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, epReq interface{}) (interface{}, error) {
		req, ok := epReq.(CreateRequest)
		if !ok {
			return CreateResponse{
				Err: common.ErrorMismatchEndpointReq}, nil
		}

		ci, err := reqctx.ClientInfoFromCtx(ctx)
		if err != nil {
			return CreateResponse{Plan: Trip{}, Err: ErrRBAC}, nil
		}

		creator := NewMember(ci.UserID, MemberRoleCreator)
		plan, err := svc.Create(ctx, creator, req.Name, req.StartDate, req.EndDate)
		return CreateResponse{Plan: plan, Err: err}, nil
	}
}

type ReadRequest struct {
	ID        string `json:"id"`
	WithUsers bool   `json:"withUsers"`
}

type ReadResponse struct {
	Plan Trip  `json:"trip"`
	Err  error `json:"error,omitempty"`
}

func (r ReadResponse) Error() error {
	return r.Err
}

type ReadWithUsersResponse struct {
	Plan  Trip          `json:"trip"`
	Users auth.UsersMap `json:"users"`
	Err   error         `json:"error,omitempty"`
}

func (r ReadWithUsersResponse) Error() error {
	return r.Err
}

func NewReadEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, epReq interface{}) (interface{}, error) {
		req, ok := epReq.(ReadRequest)
		if !ok {
			return ReadResponse{Err: common.ErrorMismatchEndpointReq}, nil
		}

		if req.WithUsers {
			plan, users, err := svc.ReadWithUsers(ctx, req.ID)
			return ReadWithUsersResponse{Plan: plan, Users: users, Err: err}, nil
		}

		plan, err := svc.Read(ctx, req.ID)
		return ReadResponse{Plan: plan, Err: err}, nil
	}
}

type ListRequest struct {
	FF ListFilter
}
type ListResponse struct {
	Plans TripsList `json:"trips"`
	Err   error     `json:"error,omitempty"`
}

func (r ListResponse) Error() error {
	return r.Err
}

func NewListEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, epReq interface{}) (interface{}, error) {
		req, ok := epReq.(ListRequest)
		if !ok {
			return ListResponse{Err: common.ErrorMismatchEndpointReq}, nil
		}
		plans, err := svc.List(ctx, req.FF)
		return ListResponse{Plans: plans, Err: err}, nil
	}
}

type DeleteRequest struct {
	ID string `json:"id"`
}

type DeleteResponse struct {
	Err error `json:"error,omitempty"`
}

func (r DeleteResponse) Error() error {
	return r.Err
}

func NewDeleteEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, epReq interface{}) (interface{}, error) {
		req, ok := epReq.(DeleteRequest)
		if !ok {
			return DeleteResponse{Err: common.ErrorMismatchEndpointReq}, nil
		}
		err := svc.Delete(ctx, req.ID)
		return DeleteResponse{Err: err}, nil
	}
}
