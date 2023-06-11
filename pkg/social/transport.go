package social

import (
	"context"
	"encoding/json"
	"net/http"

	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/travelreys/travelreys/pkg/common"
	"github.com/travelreys/travelreys/pkg/reqctx"

	"github.com/gorilla/mux"
)

const (
	URLPathVarID = "id"
)

func errToHttpCode() func(err error) int {
	notFoundErrors := []error{}
	appErrors := []error{
		ErrInvalidFriendRequest,
		ErrUnexpectedStoreError,
	}
	authErrors := []error{ErrRBAC}

	return func(err error) int {
		if common.ErrorContains(notFoundErrors, err) {
			return http.StatusNotFound
		}
		if common.ErrorContains(appErrors, err) {
			return http.StatusUnprocessableEntity
		}
		if common.ErrorContains(authErrors, err) {
			return http.StatusUnauthorized
		}
		return http.StatusInternalServerError
	}
}

func encodeResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	if e, ok := response.(common.Errorer); ok && e.Error() != nil {
		common.EncodeErrorFactory(errToHttpCode())(ctx, e.Error(), w)
		return nil
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	return json.NewEncoder(w).Encode(response)
}

func MakeHandler(svc Service) http.Handler {
	r := mux.NewRouter()

	opts := []kithttp.ServerOption{
		kithttp.ServerBefore(reqctx.ContextWithClientInfo),
		kithttp.ServerErrorEncoder(common.EncodeErrorFactory(errToHttpCode())),
	}

	sendFriendReqHandler := kithttp.NewServer(
		NewSendFriendRequestEndpoint(svc),
		decodeSendFriendRequestRequest,
		encodeResponse, opts...,
	)
	listFriendReqHandler := kithttp.NewServer(
		NewListFriendRequestsRequestEndpoint(svc),
		decodeListFriendRequestsRequest,
		encodeResponse, opts...,
	)
	acceptFriendReqHandler := kithttp.NewServer(
		NewAcceptFriendRequestEndpoint(svc),
		decodeAcceptFriendRequestRequest,
		encodeResponse, opts...,
	)

	deleteFriendReqHandler := kithttp.NewServer(
		NewDeleteFriendRequestEndpoint(svc),
		decodeDeleteFriendRequestRequest,
		encodeResponse, opts...,
	)

	listFriendsHandler := kithttp.NewServer(
		NewListFriendsRequestEndpoint(svc),
		decodeListFriendsRequest,
		encodeResponse, opts...,
	)

	deleteFriendHandler := kithttp.NewServer(
		NewDeleteFriendEndpoint(svc),
		decodeDeleteFriendRequest,
		encodeResponse, opts...,
	)

	r.Handle("/api/v1/social/requests", sendFriendReqHandler).Methods(http.MethodPost)
	r.Handle("/api/v1/social/requests", listFriendReqHandler).Methods(http.MethodGet)
	r.Handle("/api/v1/social/requests/{id}", deleteFriendReqHandler).Methods(http.MethodDelete)
	r.Handle("/api/v1/social/requests/{id}/accept", acceptFriendReqHandler).Methods(http.MethodPut)

	r.Handle("/api/v1/social/friends", listFriendsHandler).Methods(http.MethodGet)
	r.Handle("/api/v1/social/friends/{friendId}", deleteFriendHandler).Methods(http.MethodDelete)

	return r
}

func decodeSendFriendRequestRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	ci, err := reqctx.ClientInfoFromCtx(ctx)
	if err != nil {
		return nil, common.ErrInvalidRequest
	}

	req := SendFriendRequestRequest{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, common.ErrInvalidJSONBody
	}
	req.InitiatorID = ci.UserID
	return req, nil
}

func decodeAcceptFriendRequestRequest(_ context.Context, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	ID, ok := vars[URLPathVarID]
	if !ok {
		return nil, common.ErrInvalidRequest
	}
	req := AcceptFriendRequestRequest{ID: ID}
	return req, nil
}

func decodeDeleteFriendRequestRequest(_ context.Context, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	ID, ok := vars[URLPathVarID]
	if !ok {
		return nil, common.ErrInvalidRequest
	}
	req := DeleteFriendRequestRequest{ID: ID}
	return req, nil
}

func decodeListFriendRequestsRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	ci, err := reqctx.ClientInfoFromCtx(ctx)
	if err != nil {
		return nil, common.ErrInvalidRequest
	}

	req := ListFriendRequestsRequest{
		ListFriendRequestsFilter: ListFriendRequestsFilter{
			InitiatorID: common.StringPtr(ci.UserID),
			TargetID:    common.StringPtr(ci.UserID),
		},
	}
	return req, nil
}

func decodeListFriendsRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	ci, err := reqctx.ClientInfoFromCtx(ctx)
	if err != nil {
		return nil, common.ErrInvalidRequest
	}

	req := ListFriendsRequest{UserID: ci.UserID}
	return req, nil
}

func decodeDeleteFriendRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	ci, err := reqctx.ClientInfoFromCtx(ctx)
	if err != nil {
		return nil, common.ErrInvalidRequest
	}

	vars := mux.Vars(r)
	friendID, ok := vars["friendId"]
	if !ok {
		return nil, common.ErrInvalidRequest
	}
	req := DeleteFriendRequest{
		UserID:   ci.UserID,
		FriendID: friendID,
	}
	return req, nil
}
