package invites

import (
	"compress/gzip"
	"context"
	"encoding/json"
	"errors"
	"net/http"

	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	"github.com/travelreys/travelreys/pkg/common"
	"github.com/travelreys/travelreys/pkg/reqctx"
)

const (
	URLPathVarID = "id"
)

func errToHttpCode(err error) int {
	notFoundErrors := []error{ErrInviteNotFound}
	appErrors := []error{ErrUnexpectedStoreError}

	if common.ErrorContains(notFoundErrors, err) {
		return http.StatusNotFound
	}
	if common.ErrorContains(appErrors, err) {
		return http.StatusUnprocessableEntity
	}
	if errors.Is(err, ErrRBAC) {
		return http.StatusUnauthorized
	}
	if errors.Is(err, common.ErrValidation) {
		return http.StatusBadRequest
	}
	return http.StatusInternalServerError
}

func encodeResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	if e, ok := response.(common.Errorer); ok && e.Error() != nil {
		common.EncodeErrorFactory(errToHttpCode)(ctx, e.Error(), w)
		return nil
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("Cache-Control", "no-store")
	w.Header().Set("Content-Encoding", "gzip")

	gw := gzip.NewWriter(w)
	defer gw.Close()

	return json.NewEncoder(gw).Encode(response)
}

func MakeHandler(svc Service) http.Handler {
	r := mux.NewRouter()

	opts := []kithttp.ServerOption{
		kithttp.ServerBefore(reqctx.ContextWithClientInfo),
		kithttp.ServerErrorEncoder(common.EncodeErrorFactory(errToHttpCode)),
	}

	sendTripInviteHandler := kithttp.NewServer(
		NewSendTripInviteEndpoint(svc),
		decodeSendTripInviteRequest,
		encodeResponse, opts...,
	)
	listTripInviteHandler := kithttp.NewServer(
		NewListTripInvitesEndpoint(svc),
		decodeListTripInvitesRequest,
		encodeResponse, opts...,
	)
	acceptTripInviteHandler := kithttp.NewServer(
		NewAcceptTripInviteEndpoint(svc),
		decodeAcceptTripInviteRequest,
		encodeResponse, opts...,
	)
	declineTripInviteHandler := kithttp.NewServer(
		NewDeclineTripInviteEndpoint(svc),
		decodeDeclineTripInviteRequest,
		encodeResponse, opts...,
	)

	r.Handle("/api/v1/invites/trip", sendTripInviteHandler).Methods(http.MethodPost)
	r.Handle("/api/v1/invites/trip", listTripInviteHandler).Methods(http.MethodGet)
	r.Handle("/api/v1/invites/trip/{id}/accept", acceptTripInviteHandler).Methods(http.MethodPut)
	r.Handle("/api/v1/invites/trip/{id}/decline", declineTripInviteHandler).Methods(http.MethodPut)

	sendEmailTripInviteHandler := kithttp.NewServer(
		NewSendEmailTripInviteEndpoint(svc),
		decodeSendEmailTripInviteRequest,
		encodeResponse,
		opts...,
	)

	acceptEmailTripInviteHandler := kithttp.NewServer(
		NewAcceptEmailTripInviteEndpoint(svc),
		decodeAcceptEmailTripInviteRequest,
		encodeAcceptEmailTripInviteResponse, opts...,
	)

	r.Handle("/api/v1/invites/email-trip", sendEmailTripInviteHandler).Methods(http.MethodPost)
	r.Handle("/api/v1/invites/email-trip/{id}/accept", acceptEmailTripInviteHandler).Methods(http.MethodPut)

	return r
}

// Trip Invites

func decodeSendTripInviteRequest(_ context.Context, r *http.Request) (interface{}, error) {
	req := SendTripInviteRequest{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, common.ErrInvalidJSONBody
	}
	return req, nil
}

func decodeListTripInvitesRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	req := ListTripInvitesRequest{}
	ff := MakeListTripInvitesFilterFromURLParams(r.URL.Query())
	req.ListTripInvitesFilter = ff
	return req, nil
}

func decodeDeclineTripInviteRequest(_ context.Context, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	ID, ok := vars[URLPathVarID]
	if !ok {
		return nil, common.ErrInvalidRequest
	}
	return DeclineTripInviteRequest{ID}, nil
}

func decodeAcceptTripInviteRequest(_ context.Context, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	ID, ok := vars[URLPathVarID]
	if !ok {
		return nil, common.ErrInvalidRequest
	}
	return AcceptTripInviteRequest{ID}, nil
}

// Email Trip Invites

func decodeSendEmailTripInviteRequest(_ context.Context, r *http.Request) (interface{}, error) {
	req := SendEmailTripInviteRequest{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, common.ErrInvalidJSONBody
	}
	return req, nil
}

func decodeAcceptEmailTripInviteRequest(_ context.Context, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	ID, ok := vars[URLPathVarID]
	if !ok {
		return nil, common.ErrInvalidRequest
	}
	req := AcceptEmailTripInviteRequest{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, common.ErrInvalidJSONBody
	}
	req.ID = ID
	return req, nil
}

func encodeAcceptEmailTripInviteResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	if e, ok := response.(common.Errorer); ok && e.Error() != nil {
		common.EncodeErrorFactory(errToHttpCode)(ctx, e.Error(), w)
		return nil
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	resp, ok := response.(AcceptEmailTripInviteResponse)
	if !ok {
		return common.ErrorEncodeInvalidResponse
	}

	if resp.Cookie != nil {
		http.SetCookie(w, resp.Cookie)
	}
	return json.NewEncoder(w).Encode(response)
}
