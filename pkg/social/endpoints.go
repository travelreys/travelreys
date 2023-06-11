package social

import (
	"context"

	"github.com/go-kit/kit/endpoint"
	"github.com/travelreys/travelreys/pkg/common"
)

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
	ID string `json:"id"`
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
		err := svc.AcceptFriendRequest(ctx, req.ID)
		return AcceptFriendRequestResponse{Err: err}, nil
	}
}

type DeleteFriendRequestRequest struct {
	ID string `json:"id"`
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
		err := svc.DeleteFriendRequest(ctx, req.ID)
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

func NewListFriendsRequestEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, epReq interface{}) (interface{}, error) {
		req, ok := epReq.(ListFriendsRequest)
		if !ok {
			return ListFriendsResponse{Err: common.ErrorEndpointReqMismatch}, nil
		}
		friends, err := svc.ListFriends(ctx, req.UserID)
		return ListFriendsResponse{Friends: friends, Err: err}, nil
	}
}

type DeleteFriendRequest struct {
	UserID   string `json:"userID"`
	FriendID string `json:"friendID"`
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
		err := svc.DeleteFriend(ctx, req.UserID, req.FriendID)
		return DeleteFriendResponse{Err: err}, nil
	}
}
