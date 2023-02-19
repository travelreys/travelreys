package auth

import (
	"context"

	"github.com/go-kit/kit/endpoint"
	"github.com/tiinyplanet/tiinyplanet/pkg/common"
)

type LoginRequest struct {
	Code     string `json:"code"`
	Provider string `json:"provider"`
}

type LoginResponse struct {
	JWTToken string `json:"jwtToken"`
	Err      error  `json:"error,omitempty"`
}

func (r LoginResponse) Error() error {
	return r.Err
}

func NewLoginEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, epReq interface{}) (interface{}, error) {
		req, ok := epReq.(LoginRequest)
		if !ok {
			return LoginResponse{Err: common.ErrorInvalidEndpointRequestType}, nil
		}
		jwtTkn, err := svc.Login(ctx, req.Code, OIDCProviderGoogle)
		return LoginResponse{JWTToken: jwtTkn, Err: err}, nil
	}
}

type ReadUserRequest struct {
	ID string `json:"id"`
}

type ReadUserResponse struct {
	User User  `json:"user"`
	Err  error `json:"error,omitempty"`
}

func (r ReadUserResponse) Error() error {
	return r.Err
}

func NewReadUserEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, epReq interface{}) (interface{}, error) {
		req, ok := epReq.(ReadUserRequest)
		if !ok {
			return ReadUserResponse{
				Err: common.ErrorInvalidEndpointRequestType,
			}, nil
		}
		usr, err := svc.ReadUser(ctx, req.ID)
		return ReadUserResponse{usr, err}, nil
	}
}

type UpdateUserRequest struct {
	ID string           `json:"id"`
	FF UpdateUserFilter `json:"ff"`
}

type UpdateUserResponse struct {
	Err error `json:"error,omitempty"`
}

func (r UpdateUserResponse) Error() error {
	return r.Err
}

func NewUpdateUserEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, epReq interface{}) (interface{}, error) {
		req, ok := epReq.(UpdateUserRequest)
		if !ok {
			return UpdateUserResponse{
				Err: common.ErrorInvalidEndpointRequestType,
			}, nil
		}
		err := svc.UpdateUser(ctx, req.ID, req.FF)
		return UpdateUserResponse{err}, nil
	}
}

type ListUsersRequest struct {
	FF ListUsersFilter
}
type ListUsersResponse struct {
	Users UsersList `json:"users"`
	Err   error     `json:"error,omitempty"`
}

func (r ListUsersResponse) Error() error {
	return r.Err
}

func NewListUsersEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, epReq interface{}) (interface{}, error) {
		req, ok := epReq.(ListUsersRequest)
		if !ok {
			return ListUsersResponse{Err: common.ErrorInvalidEndpointRequestType}, nil
		}
		users, err := svc.ListUsers(ctx, req.FF)
		return ListUsersResponse{Users: users, Err: err}, nil
	}
}
