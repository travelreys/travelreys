package trips

import (
	"context"

	"github.com/go-kit/kit/endpoint"
	"github.com/travelreys/travelreys/pkg/common"
)

type SendInviteRequest struct {
	UserID   string `json:"userID"`
	AuthorID string `json:"authorID"`
	TripID   string `json:"tripID"`
}

type SendInviteResponse struct {
	Err error `json:"error,omitempty"`
}

func (r SendInviteResponse) Error() error {
	return r.Err
}

func NewSendEndpoint(svc InviteService) endpoint.Endpoint {
	return func(ctx context.Context, epReq interface{}) (interface{}, error) {
		req, ok := epReq.(SendInviteRequest)
		if !ok {
			return SendInviteResponse{
				Err: common.ErrEndpointReqMismatch,
			}, nil
		}
		err := svc.Send(ctx, req.TripID, req.AuthorID, req.UserID)
		return SendInviteResponse{Err: err}, nil
	}
}

type AcceptInviteRequest struct {
	ID string `json:"id"`
}
type AcceptInviteResponse struct {
	Err error `json:"error,omitempty"`
}

func (r AcceptInviteResponse) Error() error {
	return r.Err
}

func NewAcceptEndpoint(svc InviteService) endpoint.Endpoint {
	return func(ctx context.Context, epReq interface{}) (interface{}, error) {
		req, ok := epReq.(AcceptInviteRequest)
		if !ok {
			return AcceptInviteResponse{
				Err: common.ErrEndpointReqMismatch,
			}, nil
		}
		err := svc.Accept(ctx, req.ID)
		return AcceptInviteResponse{Err: err}, nil
	}
}

type DeclineInviteRequest struct {
	ID string `json:"id"`
}

type DeclineInviteResponse struct {
	Err error `json:"error,omitempty"`
}

func (r DeclineInviteResponse) Error() error {
	return r.Err
}

func NewDeclineEndpoint(svc InviteService) endpoint.Endpoint {
	return func(ctx context.Context, epReq interface{}) (interface{}, error) {
		req, ok := epReq.(DeclineInviteRequest)
		if !ok {
			return DeclineInviteResponse{
				Err: common.ErrEndpointReqMismatch,
			}, nil
		}
		err := svc.Decline(ctx, req.ID)
		return DeclineInviteResponse{err}, nil
	}
}

type ListInvitesRequest struct {
	ListInvitesFilter
}

type ListInvitesResponse struct {
	Invites InviteList `json:"invites"`
	Err     error      `json:"error,omitempty"`
}

func (r ListInvitesResponse) Error() error {
	return r.Err
}

func NewListInvitesEndpoint(svc InviteService) endpoint.Endpoint {
	return func(ctx context.Context, epReq interface{}) (interface{}, error) {
		req, ok := epReq.(ListInvitesRequest)
		if !ok {
			return ListInvitesResponse{
				Err: common.ErrEndpointReqMismatch,
			}, nil
		}
		invites, err := svc.List(ctx, req.ListInvitesFilter)
		return ListInvitesResponse{invites, err}, nil
	}
}
