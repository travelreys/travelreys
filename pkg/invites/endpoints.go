package invites

import (
	"context"

	"github.com/go-kit/kit/endpoint"
	"github.com/travelreys/travelreys/pkg/common"
)

type SendTripInviteRequest struct {
	UserID   string `json:"userID"`
	AuthorID string `json:"authorID"`
	TripID   string `json:"tripID"`
}

type SendTripInviteResponse struct {
	Err error `json:"error,omitempty"`
}

func (r SendTripInviteResponse) Error() error {
	return r.Err
}

func NewSendTripInviteEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, epReq interface{}) (interface{}, error) {
		req, ok := epReq.(SendTripInviteRequest)
		if !ok {
			return SendTripInviteResponse{
				Err: common.ErrEndpointReqMismatch,
			}, nil
		}
		err := svc.SendTripInvite(ctx, req.TripID, req.AuthorID, req.UserID)
		return SendTripInviteResponse{Err: err}, nil
	}
}

type AcceptTripInviteRequest struct {
	ID string `json:"id"`
}
type AcceptTripInviteResponse struct {
	Err error `json:"error,omitempty"`
}

func (r AcceptTripInviteResponse) Error() error {
	return r.Err
}

func NewAcceptTripInviteEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, epReq interface{}) (interface{}, error) {
		req, ok := epReq.(AcceptTripInviteRequest)
		if !ok {
			return AcceptTripInviteResponse{
				Err: common.ErrEndpointReqMismatch,
			}, nil
		}
		err := svc.AcceptTripInvite(ctx, req.ID)
		return AcceptTripInviteResponse{Err: err}, nil
	}
}

type DeclineTripInviteRequest struct {
	ID string `json:"id"`
}

type DeclineTripInviteResponse struct {
	Err error `json:"error,omitempty"`
}

func (r DeclineTripInviteResponse) Error() error {
	return r.Err
}

func NewDeclineTripInviteEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, epReq interface{}) (interface{}, error) {
		req, ok := epReq.(DeclineTripInviteRequest)
		if !ok {
			return DeclineTripInviteResponse{
				Err: common.ErrEndpointReqMismatch,
			}, nil
		}
		err := svc.DeclineTripInvite(ctx, req.ID)
		return DeclineTripInviteResponse{err}, nil
	}
}

type ListTripInvitesRequest struct {
	ListTripInvitesFilter
}

type ListTripInvitesResponse struct {
	TripInvites TripInviteList `json:"invites"`
	Err         error          `json:"error,omitempty"`
}

func (r ListTripInvitesResponse) Error() error {
	return r.Err
}

func NewListTripInvitesEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, epReq interface{}) (interface{}, error) {
		req, ok := epReq.(ListTripInvitesRequest)
		if !ok {
			return ListTripInvitesResponse{
				Err: common.ErrEndpointReqMismatch,
			}, nil
		}
		invites, err := svc.ListTripInvites(ctx, req.ListTripInvitesFilter)
		return ListTripInvitesResponse{invites, err}, nil
	}
}
