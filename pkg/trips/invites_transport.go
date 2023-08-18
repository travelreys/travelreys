package trips

import (
	"context"
	"encoding/json"
	"net/http"

	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	"github.com/travelreys/travelreys/pkg/common"
	"github.com/travelreys/travelreys/pkg/reqctx"
)

func MakeInviteHandler(svc InviteService) http.Handler {
	r := mux.NewRouter()

	opts := []kithttp.ServerOption{
		kithttp.ServerBefore(reqctx.ContextWithClientInfo),
		kithttp.ServerErrorEncoder(common.EncodeErrorFactory(errToHttpCode)),
	}

	sendHandler := kithttp.NewServer(
		NewSendEndpoint(svc),
		decodeSendRequest,
		encodeResponse, opts...,
	)
	listHandler := kithttp.NewServer(
		NewListInvitesEndpoint(svc),
		decodeListInvitesRequest,
		encodeResponse, opts...,
	)
	acceptHandler := kithttp.NewServer(
		NewAcceptEndpoint(svc),
		decodeAcceptInviteRequest,
		encodeResponse, opts...,
	)
	declineHandler := kithttp.NewServer(
		NewDeclineEndpoint(svc),
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
	req := SendInviteRequest{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, common.ErrInvalidJSONBody
	}
	return req, nil
}

func decodeListInvitesRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	req := ListInvitesRequest{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, common.ErrInvalidJSONBody
	}
	return req, nil
}

func decodeDeclineInviteRequest(_ context.Context, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	ID, ok := vars[URLPathVarID]
	if !ok {
		return nil, common.ErrInvalidRequest
	}
	return DeclineInviteRequest{ID}, nil
}

func decodeAcceptInviteRequest(_ context.Context, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	ID, ok := vars[URLPathVarID]
	if !ok {
		return nil, common.ErrInvalidRequest
	}
	return AcceptInviteRequest{ID}, nil
}
