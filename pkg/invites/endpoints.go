package invites

import (
	"context"
	"fmt"
	"net/http"

	"github.com/go-kit/kit/endpoint"
	"github.com/travelreys/travelreys/pkg/auth"
	"github.com/travelreys/travelreys/pkg/common"
)

// Trip Invites

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

// Email Trip Invites

type SendEmailTripInviteRequest struct {
	UserEmail string `json:"userEmail"`
	AuthorID  string `json:"authorID"`
	TripID    string `json:"tripID"`
}

type SendEmailTripInviteResponse struct {
	Err error `json:"error,omitempty"`
}

func (r SendEmailTripInviteResponse) Error() error {
	return r.Err
}

func NewSendEmailTripInviteEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, epReq interface{}) (interface{}, error) {
		req, ok := epReq.(SendEmailTripInviteRequest)
		if !ok {
			return SendEmailTripInviteResponse{
				Err: common.ErrEndpointReqMismatch,
			}, nil
		}
		err := svc.SendEmailTripInvite(
			ctx, req.TripID, req.AuthorID, req.UserEmail,
		)
		return SendEmailTripInviteResponse{Err: err}, nil
	}
}

type AcceptEmailTripInviteRequest struct {
	ID         string `json:"id"`
	Code       string `json:"code"`
	Sig        string `json:"sig"`
	IsLoggedIn bool   `json:"isLoggedIn"`
}
type AcceptEmailTripInviteResponse struct {
	User   auth.User    `json:"user"`
	Cookie *http.Cookie `json:"-"`
	Err    error        `json:"error,omitempty"`
}

func (r AcceptEmailTripInviteResponse) Error() error {
	return r.Err
}

func NewAcceptEmailTripInviteEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, epReq interface{}) (interface{}, error) {
		req, ok := epReq.(AcceptEmailTripInviteRequest)
		if !ok {
			return AcceptEmailTripInviteResponse{
				Err: common.ErrEndpointReqMismatch,
			}, nil
		}
		fmt.Println("ep code", req.Code)
		fmt.Println("ep sig", req.Sig)
		user, cookie, err := svc.AcceptEmailTripInvite(
			ctx, req.ID, req.Code, req.Sig, req.IsLoggedIn,
		)
		return AcceptEmailTripInviteResponse{
			User: user, Cookie: cookie, Err: err,
		}, nil
	}
}
