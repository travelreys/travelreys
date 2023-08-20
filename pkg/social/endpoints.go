package social

import (
	"context"
	"time"

	"github.com/go-kit/kit/endpoint"
	"github.com/travelreys/travelreys/pkg/common"
	"github.com/travelreys/travelreys/pkg/trips"
)

type GetProfileRequest struct {
	ID string `json:"id"`
}

type GetProfileResponse struct {
	Profile UserProfile `json:"profile"`
	Err     error
}

func (r GetProfileResponse) Error() error {
	return r.Err
}

func NewGetProfileRequestEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, epReq interface{}) (interface{}, error) {
		req, ok := epReq.(GetProfileRequest)
		if !ok {
			return GetProfileResponse{Err: common.ErrEndpointReqMismatch}, nil
		}
		profile, err := svc.GetProfile(ctx, req.ID)
		return GetProfileResponse{Profile: profile, Err: err}, nil
	}
}

type SendFollowRequestRequest struct {
	InitiatorID string `json:"initiatorID"`
	TargetID    string `json:"targetID"`
}
type SendFollowRequestResponse struct {
	Err error `json:"error,omitempty"`
}

func (r SendFollowRequestResponse) Error() error {
	return r.Err
}

func NewSendFollowRequestEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, epReq interface{}) (interface{}, error) {
		req, ok := epReq.(SendFollowRequestRequest)
		if !ok {
			return SendFollowRequestResponse{Err: common.ErrEndpointReqMismatch}, nil
		}
		err := svc.SendFollowRequest(ctx, req.InitiatorID, req.TargetID)
		return SendFollowRequestResponse{Err: err}, nil
	}
}

type AcceptFollowRequestRequest struct {
	UserID      string `json:"userID"`
	InitiatorID string `json:"initiatorID"`
	RequestID   string `json:"requestID"`
}
type AcceptFollowRequestResponse struct {
	Err error `json:"error,omitempty"`
}

func (r AcceptFollowRequestResponse) Error() error {
	return r.Err
}

func NewAcceptFollowRequestEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, epReq interface{}) (interface{}, error) {
		req, ok := epReq.(AcceptFollowRequestRequest)
		if !ok {
			return AcceptFollowRequestResponse{Err: common.ErrEndpointReqMismatch}, nil
		}
		err := svc.AcceptFollowRequest(
			ctx,
			req.UserID,
			req.InitiatorID,
			req.RequestID,
		)
		return AcceptFollowRequestResponse{Err: err}, nil
	}
}

type DeleteFollowRequestRequest struct {
	UserID    string `json:"userID"`
	RequestID string `json:"requestID"`
}
type DeleteFollowRequestResponse struct {
	Err error `json:"error,omitempty"`
}

func (r DeleteFollowRequestResponse) Error() error {
	return r.Err
}

func NewDeleteFollowRequestEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, epReq interface{}) (interface{}, error) {
		req, ok := epReq.(DeleteFollowRequestRequest)
		if !ok {
			return DeleteFollowRequestResponse{Err: common.ErrEndpointReqMismatch}, nil
		}
		err := svc.DeleteFollowRequest(ctx, req.UserID, req.RequestID)
		return DeleteFollowRequestResponse{Err: err}, nil
	}
}

type ListFollowRequestsRequest struct {
	ListFollowRequestsFilter
}
type ListFollowRequestsResponse struct {
	Requests FollowRequestList `json:"requests"`
	Err      error             `json:"error,omitempty"`
}

func (r ListFollowRequestsResponse) Error() error {
	return r.Err
}

func NewListFollowRequestsRequestEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, epReq interface{}) (interface{}, error) {
		req, ok := epReq.(ListFollowRequestsRequest)
		if !ok {
			return ListFollowRequestsResponse{Err: common.ErrEndpointReqMismatch}, nil
		}
		reqs, err := svc.ListFollowRequests(ctx, req.ListFollowRequestsFilter)
		return ListFollowRequestsResponse{Requests: reqs, Err: err}, nil
	}
}

type IsFollowingRequest struct {
	InitiatorID string `json:"initiatorID"`
	TargetID    string `json:"targetID"`
}

type IsFollowingResponse struct {
	OK  bool  `json:"ok"`
	Err error `json:"error,omitempty"`
}

func (r IsFollowingResponse) Error() error {
	return r.Err
}

func NewIsFollowingResponseEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, epReq interface{}) (interface{}, error) {
		req, ok := epReq.(IsFollowingRequest)
		if !ok {
			return ListFollowRequestsResponse{Err: common.ErrEndpointReqMismatch}, nil
		}
		ok, err := svc.IsFollowing(ctx, req.InitiatorID, req.TargetID)
		return IsFollowingResponse{OK: ok, Err: err}, nil
	}
}

type ListFriendsRequest struct {
	UserID string `json:"userID"`
}
type ListFriendsResponse struct {
	Friends FollowingsList `json:"friends"`
	Err     error          `json:"error,omitempty"`
}

func (r ListFriendsResponse) Error() error {
	return r.Err
}

func NewListFollowersRequestEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, epReq interface{}) (interface{}, error) {
		req, ok := epReq.(ListFriendsRequest)
		if !ok {
			return ListFriendsResponse{Err: common.ErrEndpointReqMismatch}, nil
		}
		friends, err := svc.ListFollowers(ctx, req.UserID)
		return ListFriendsResponse{Friends: friends, Err: err}, nil
	}
}

func NewListFollowingRequestEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, epReq interface{}) (interface{}, error) {
		req, ok := epReq.(ListFriendsRequest)
		if !ok {
			return ListFriendsResponse{Err: common.ErrEndpointReqMismatch}, nil
		}
		friends, err := svc.ListFollowing(ctx, req.UserID)
		return ListFriendsResponse{Friends: friends, Err: err}, nil
	}
}

type DeleteFollowRequest struct {
	UserID     string `json:"userID"`
	BindingKey string `json:"bindingKey"`
}
type DeleteFollowResponse struct {
	Err error `json:"error,omitempty"`
}

func (r DeleteFollowResponse) Error() error {
	return r.Err
}

func NewDeleteFollowEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, epReq interface{}) (interface{}, error) {
		req, ok := epReq.(DeleteFollowRequest)
		if !ok {
			return DeleteFollowResponse{Err: common.ErrEndpointReqMismatch}, nil
		}
		err := svc.DeleteFollowing(ctx, req.UserID, req.BindingKey)
		return DeleteFollowResponse{Err: err}, nil
	}
}

type ReadRequest struct {
	TripID     string `json:"id"`
	ReferrerID string `json:"referrerID"`
}
type ReadResponse struct {
	Trip        *trips.Trip `json:"trip"`
	UserProfile UserProfile `json:"profile"`
	Err         error       `json:"error,omitempty"`
}

func (r ReadResponse) Error() error {
	return r.Err
}

func NewReadTripPublicInfoEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, epReq interface{}) (interface{}, error) {
		req, ok := epReq.(ReadRequest)
		if !ok {
			return ReadResponse{Err: common.ErrEndpointReqMismatch}, nil
		}
		trip, profile, err := svc.ReadTripPublicInfo(ctx, req.TripID, req.ReferrerID)
		return ReadResponse{Trip: trip, UserProfile: profile, Err: err}, nil
	}
}

type ListRequest struct {
	trips.ListFilter
}
type ListResponse struct {
	Trips trips.TripsList `json:"trips"`
	Err   error           `json:"error,omitempty"`
}

func (r ListResponse) Error() error {
	return r.Err
}

func NewListTripPublicInfoEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, epReq interface{}) (interface{}, error) {
		req, ok := epReq.(ListRequest)
		if !ok {
			return ListResponse{Err: common.ErrEndpointReqMismatch}, nil
		}
		tripslist, err := svc.ListTripPublicInfo(ctx, req.ListFilter)
		return ListResponse{Trips: tripslist, Err: err}, nil
	}
}

type ListFollowingTripsRequest struct {
	InitiatorID string `json:"initiatorID"`
}

type ListFollowingTripsResponse struct {
	Trips          trips.TripsList `json:"trips"`
	UserProfileMap UserProfileMap  `json:"profiles"`
	Err            error           `json:"error,omitempty"`
}

func (r ListFollowingTripsResponse) Error() error {
	return r.Err
}

func NewListFollowingTripsEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, epReq interface{}) (interface{}, error) {
		req, ok := epReq.(ListFollowingTripsRequest)
		if !ok {
			return ListFollowingTripsResponse{Err: common.ErrEndpointReqMismatch}, nil
		}
		tripslist, profiles, err := svc.ListFollowingTrips(ctx, req.InitiatorID)
		return ListFollowingTripsResponse{Trips: tripslist, UserProfileMap: profiles, Err: err}, nil
	}
}

type DuplicateRequest struct {
	TripID     string    `json:"id"`
	Name       string    `json:"name"`
	StartDate  time.Time `json:"startDate"`
	ReferrerID string    `json:"referrerID"`
}
type DuplicateResponse struct {
	ID  string `json:"id"`
	Err error  `json:"error,omitempty"`
}

func (r DuplicateResponse) Error() error {
	return r.Err
}

func NewDuplicateTripEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, epReq interface{}) (interface{}, error) {
		req, ok := epReq.(DuplicateRequest)
		if !ok {
			return DuplicateResponse{Err: common.ErrEndpointReqMismatch}, nil
		}
		id, err := svc.DuplicateTrip(
			ctx, "", req.ReferrerID, req.TripID, req.Name, req.StartDate,
		)
		return DuplicateResponse{ID: id, Err: err}, nil
	}
}
