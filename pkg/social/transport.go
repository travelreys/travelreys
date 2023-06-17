package social

import (
	"context"
	"encoding/json"
	"net/http"

	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/travelreys/travelreys/pkg/common"
	"github.com/travelreys/travelreys/pkg/reqctx"
	"github.com/travelreys/travelreys/pkg/trips"

	"github.com/gorilla/mux"
)

const (
	URLPathVarUserID    = "uid"
	URLPathVarTripID    = "tid"
	URLPathVarRequestID = "rid"
	URLPathBindingKey   = "bindingKey"
)

func errToHttpCode() func(err error) int {
	notFoundErrors := []error{
		ErrFriendNotFound,
		trips.ErrTripNotFound,
	}
	appErrors := []error{
		ErrInvalidFriendRequest,
		ErrUnexpectedStoreError,
	}
	authErrors := []error{
		ErrTripSharingNotEnabled,
		ErrRBAC,
	}

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

	getProfileHandler := kithttp.NewServer(
		NewGetProfileRequestEndpoint(svc),
		decodeGetProfileRequest,
		encodeResponse, opts...,
	)

	listFollowingTripsHandler := kithttp.NewServer(
		NewListFollowingTripsEndpoint(svc),
		decodeListFollowingTripsRequest,
		encodeResponse, opts...,
	)

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

	getFriendHandler := kithttp.NewServer(
		NewAreTheyFriendsResponseEndpoint(svc),
		decodeAreTheyFriendsRequest,
		encodeResponse, opts...,
	)

	listFollowersHandler := kithttp.NewServer(
		NewListFollowersRequestEndpoint(svc),
		decodeListFriendsRequest,
		encodeResponse, opts...,
	)

	listFollowingHandler := kithttp.NewServer(
		NewListFollowingRequestEndpoint(svc),
		decodeListFriendsRequest,
		encodeResponse, opts...,
	)

	deleteFriendHandler := kithttp.NewServer(
		NewDeleteFriendEndpoint(svc),
		decodeDeleteFriendRequest,
		encodeResponse, opts...,
	)

	ListTripPublicInfo := kithttp.NewServer(
		NewListTripPublicInfoEndpoint(svc), decodeListTripPublicInfoRequest, encodeResponse, opts...,
	)

	ReadTripPublicInfoHandler := kithttp.NewServer(
		NewReadTripPublicInfoEndpoint(svc), decodeReadTripPublicInfoRequest, encodeResponse, opts...,
	)

	r.Handle("/api/v1/social/{uid}", listFollowingTripsHandler).Methods(http.MethodGet)
	r.Handle("/api/v1/social/{uid}/profile", getProfileHandler).Methods(http.MethodGet)

	r.Handle("/api/v1/social/{uid}/requests", sendFriendReqHandler).Methods(http.MethodPost)
	r.Handle("/api/v1/social/{uid}/requests", listFriendReqHandler).Methods(http.MethodGet)
	r.Handle("/api/v1/social/{uid}/requests/{rid}/accept", acceptFriendReqHandler).Methods(http.MethodPut)
	r.Handle("/api/v1/social/{uid}/requests/{rid}", deleteFriendReqHandler).Methods(http.MethodDelete)

	r.Handle("/api/v1/social/{uid}/friends/following", listFollowingHandler).Methods(http.MethodGet)
	r.Handle("/api/v1/social/{uid}/friends/followers", listFollowersHandler).Methods(http.MethodGet)
	r.Handle("/api/v1/social/{uid}/friends/{targetID}", getFriendHandler).Methods(http.MethodGet)
	r.Handle("/api/v1/social/{uid}/friends/{bindingKey}", deleteFriendHandler).Methods(http.MethodDelete)

	r.Handle("/api/v1/social/{uid}/trips", ListTripPublicInfo).Methods(http.MethodGet)
	r.Handle("/api/v1/social/{uid}/trips/{tid}", ReadTripPublicInfoHandler).Methods(http.MethodGet)

	return r
}

func decodeGetProfileRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	ID, ok := vars[URLPathVarUserID]
	if !ok {
		return nil, common.ErrInvalidRequest
	}
	req := GetProfileRequest{ID: ID}
	return req, nil
}

func decodeListFollowingTripsRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	ID, ok := vars[URLPathVarUserID]
	if !ok {
		return nil, common.ErrInvalidRequest
	}
	req := ListFollowingTripsRequest{InitiatorID: ID}
	return req, nil
}

func decodeSendFriendRequestRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	userID, ok := vars[URLPathVarUserID]
	if !ok {
		return nil, common.ErrInvalidRequest
	}
	req := SendFriendRequestRequest{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, common.ErrInvalidJSONBody
	}
	req.InitiatorID = userID
	return req, nil
}

func decodeAcceptFriendRequestRequest(_ context.Context, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	userID, ok := vars[URLPathVarUserID]
	if !ok {
		return nil, common.ErrInvalidRequest
	}
	requestID, ok := vars[URLPathVarRequestID]
	if !ok {
		return nil, common.ErrInvalidRequest
	}
	req := AcceptFriendRequestRequest{
		UserID:    userID,
		RequestID: requestID,
	}
	return req, nil
}

func decodeDeleteFriendRequestRequest(_ context.Context, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	userID, ok := vars[URLPathVarUserID]
	if !ok {
		return nil, common.ErrInvalidRequest
	}
	requestID, ok := vars[URLPathVarRequestID]
	if !ok {
		return nil, common.ErrInvalidRequest
	}
	req := DeleteFriendRequestRequest{UserID: userID, RequestID: requestID}
	return req, nil
}

func decodeListFriendRequestsRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	userID, ok := vars[URLPathVarUserID]
	if !ok {
		return nil, common.ErrInvalidRequest
	}
	req := ListFriendRequestsRequest{
		ListFriendRequestsFilter: ListFriendRequestsFilter{
			TargetID: common.StringPtr(userID),
		},
	}
	return req, nil
}

func decodeAreTheyFriendsRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	userID, ok := vars[URLPathVarUserID]
	if !ok {
		return nil, common.ErrInvalidRequest
	}
	targetID, ok := vars["targetID"]
	if !ok {
		return nil, common.ErrInvalidRequest
	}

	req := AreTheyFriendsRequest{
		InitiatorID: userID,
		TargetID:    targetID,
	}
	return req, nil
}

func decodeListFriendsRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	userID, ok := vars[URLPathVarUserID]
	if !ok {
		return nil, common.ErrInvalidRequest
	}
	req := ListFriendsRequest{UserID: userID}
	return req, nil
}

func decodeDeleteFriendRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	userID, ok := vars[URLPathVarUserID]
	if !ok {
		return nil, common.ErrInvalidRequest
	}
	bindingKey, ok := vars[URLPathBindingKey]
	if !ok {
		return nil, common.ErrInvalidRequest
	}
	req := DeleteFriendRequest{
		UserID:     userID,
		BindingKey: bindingKey,
	}
	return req, nil
}

func decodeReadTripPublicInfoRequest(_ context.Context, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	userID, ok := vars[URLPathVarUserID]
	if !ok {
		return nil, common.ErrInvalidRequest
	}
	tripID, ok := vars[URLPathVarTripID]
	if !ok {
		return nil, common.ErrInvalidRequest
	}
	return ReadRequest{
		TripID:     tripID,
		ReferrerID: userID,
	}, nil
}

func decodeListTripPublicInfoRequest(_ context.Context, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	userID, ok := vars[URLPathVarUserID]
	if !ok {
		return nil, common.ErrInvalidRequest
	}

	ff := trips.ListFilter{UserID: common.StringPtr(userID)}
	return ListRequest{ff}, nil
}
