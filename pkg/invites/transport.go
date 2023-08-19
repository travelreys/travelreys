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

	sendHandler := kithttp.NewServer(
		NewSendTripInviteEndpoint(svc),
		decodeSendRequest,
		encodeResponse, opts...,
	)
	listHandler := kithttp.NewServer(
		NewListTripInvitesEndpoint(svc),
		decodeListInvitesRequest,
		encodeResponse, opts...,
	)
	acceptHandler := kithttp.NewServer(
		NewAcceptTripInviteEndpoint(svc),
		decodeAcceptInviteRequest,
		encodeResponse, opts...,
	)
	declineHandler := kithttp.NewServer(
		NewDeclineTripInviteEndpoint(svc),
		decodeDeclineInviteRequest,
		encodeResponse, opts...,
	)

	r.Handle("/api/v1/trip-invites", sendHandler).Methods(http.MethodPost)
	r.Handle("/api/v1/trip-invites", listHandler).Methods(http.MethodGet)
	r.Handle("/api/v1/trip-invites/{id}/accept", acceptHandler).Methods(http.MethodPut)
	r.Handle("/api/v1/trip-invites/{id}/decline", declineHandler).Methods(http.MethodPut)

	return r
}

func decodeSendRequest(_ context.Context, r *http.Request) (interface{}, error) {
	req := SendTripInviteRequest{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, common.ErrInvalidJSONBody
	}
	return req, nil
}

func decodeListInvitesRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	req := ListTripInvitesRequest{}
	ff := MakeListTripInvitesFilterFromURLParams(r.URL.Query())
	req.ListTripInvitesFilter = ff
	return req, nil
}

func decodeDeclineInviteRequest(_ context.Context, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	ID, ok := vars[URLPathVarID]
	if !ok {
		return nil, common.ErrInvalidRequest
	}
	return DeclineTripInviteRequest{ID}, nil
}

func decodeAcceptInviteRequest(_ context.Context, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	ID, ok := vars[URLPathVarID]
	if !ok {
		return nil, common.ErrInvalidRequest
	}
	return AcceptTripInviteRequest{ID}, nil
}
