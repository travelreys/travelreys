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
			return GetProfileResponse{Err: common.ErrorEndpointReqMismatch}, nil
		}
		profile, err := svc.GetProfile(ctx, req.ID)
		return GetProfileResponse{Profile: profile, Err: err}, nil
	}
}

type SendFriendRequestRequest struct {
	InitiatorID string `json:"initiatorID"`
	TargetID    string `json:"targetID"`
}
type SendFriendRequestResponse struct {
	Err error `json:"error,omitempty"`
}

func (r SendFriendRequestResponse) Error() error {
	return r.Err
}

func NewSendFriendRequestEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, epReq interface{}) (interface{}, error) {
		req, ok := epReq.(SendFriendRequestRequest)
		if !ok {
			return SendFriendRequestResponse{Err: common.ErrorEndpointReqMismatch}, nil
		}
		err := svc.SendFriendRequest(ctx, req.InitiatorID, req.TargetID)
		return SendFriendRequestResponse{Err: err}, nil
	}
}

type AcceptFriendRequestRequest struct {
	UserID    string `json:"uid"`
	RequestID string `json:"rid"`
}
type AcceptFriendRequestResponse struct {
	Err error `json:"error,omitempty"`
}

func (r AcceptFriendRequestResponse) Error() error {
	return r.Err
}

func NewAcceptFriendRequestEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, epReq interface{}) (interface{}, error) {
		req, ok := epReq.(AcceptFriendRequestRequest)
		if !ok {
			return AcceptFriendRequestResponse{Err: common.ErrorEndpointReqMismatch}, nil
		}
		err := svc.AcceptFriendRequest(ctx, req.UserID, req.RequestID)
		return AcceptFriendRequestResponse{Err: err}, nil
	}
}

type DeleteFriendRequestRequest struct {
	UserID    string `json:"uid"`
	RequestID string `json:"rid"`
}
type DeleteFriendRequestResponse struct {
	Err error `json:"error,omitempty"`
}

func (r DeleteFriendRequestResponse) Error() error {
	return r.Err
}

func NewDeleteFriendRequestEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, epReq interface{}) (interface{}, error) {
		req, ok := epReq.(DeleteFriendRequestRequest)
		if !ok {
			return DeleteFriendRequestResponse{Err: common.ErrorEndpointReqMismatch}, nil
		}
		err := svc.DeleteFriendRequest(ctx, req.UserID, req.RequestID)
		return DeleteFriendRequestResponse{Err: err}, nil
	}
}

type ListFriendRequestsRequest struct {
	ListFriendRequestsFilter
}
type ListFriendRequestsResponse struct {
	Requests FriendRequestList `json:"requests"`
	Err      error             `json:"error,omitempty"`
}

func (r ListFriendRequestsResponse) Error() error {
	return r.Err
}

func NewListFriendRequestsRequestEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, epReq interface{}) (interface{}, error) {
		req, ok := epReq.(ListFriendRequestsRequest)
		if !ok {
			return ListFriendRequestsResponse{Err: common.ErrorEndpointReqMismatch}, nil
		}
		reqs, err := svc.ListFriendRequests(ctx, req.ListFriendRequestsFilter)
		return ListFriendRequestsResponse{Requests: reqs, Err: err}, nil
	}
}

type AreTheyFriendsRequest struct {
	InitiatorID string `json:"initiatorID"`
	TargetID    string `json:"targetID"`
}

type AreTheyFriendsResponse struct {
	OK  bool  `json:"ok"`
	Err error `json:"error,omitempty"`
}

func (r AreTheyFriendsResponse) Error() error {
	return r.Err
}

func NewAreTheyFriendsResponseEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, epReq interface{}) (interface{}, error) {
		req, ok := epReq.(AreTheyFriendsRequest)
		if !ok {
			return ListFriendRequestsResponse{Err: common.ErrorEndpointReqMismatch}, nil
		}
		ok, err := svc.AreTheyFriends(ctx, req.InitiatorID, req.TargetID)
		return AreTheyFriendsResponse{OK: ok, Err: err}, nil
	}
}

type ListFriendsRequest struct {
	UserID string `json:"userID"`
}
type ListFriendsResponse struct {
	Friends FriendsList `json:"friends"`
	Err     error       `json:"error,omitempty"`
}

func (r ListFriendsResponse) Error() error {
	return r.Err
}

func NewListFollowersRequestEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, epReq interface{}) (interface{}, error) {
		req, ok := epReq.(ListFriendsRequest)
		if !ok {
			return ListFriendsResponse{Err: common.ErrorEndpointReqMismatch}, nil
		}
		friends, err := svc.ListFollowers(ctx, req.UserID)
		return ListFriendsResponse{Friends: friends, Err: err}, nil
	}
}

func NewListFollowingRequestEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, epReq interface{}) (interface{}, error) {
		req, ok := epReq.(ListFriendsRequest)
		if !ok {
			return ListFriendsResponse{Err: common.ErrorEndpointReqMismatch}, nil
		}
		friends, err := svc.ListFollowing(ctx, req.UserID)
		return ListFriendsResponse{Friends: friends, Err: err}, nil
	}
}

type DeleteFriendRequest struct {
	UserID     string `json:"userID"`
	BindingKey string `json:"bindingKey"`
}
type DeleteFriendResponse struct {
	Err error `json:"error,omitempty"`
}

func (r DeleteFriendResponse) Error() error {
	return r.Err
}

func NewDeleteFriendEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, epReq interface{}) (interface{}, error) {
		req, ok := epReq.(DeleteFriendRequest)
		if !ok {
			return DeleteFriendResponse{Err: common.ErrorEndpointReqMismatch}, nil
		}
		err := svc.DeleteFriend(ctx, req.UserID, req.BindingKey)
		return DeleteFriendResponse{Err: err}, nil
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
			return ReadResponse{Err: common.ErrorEndpointReqMismatch}, nil
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
			return ListResponse{Err: common.ErrorEndpointReqMismatch}, nil
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
			return ListFollowingTripsResponse{Err: common.ErrorEndpointReqMismatch}, nil
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
			return DuplicateResponse{Err: common.ErrorEndpointReqMismatch}, nil
		}
		id, err := svc.DuplicateTrip(
			ctx, "", req.ReferrerID, req.TripID, req.Name, req.StartDate,
		)
		return DuplicateResponse{ID: id, Err: err}, nil
	}
}
