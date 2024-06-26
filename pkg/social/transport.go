package social

import (
	"compress/gzip"
	"context"
	"encoding/json"
	"errors"
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
		ErrFollowingNotFound,
		trips.ErrTripNotFound,
	}
	appErrors := []error{
		ErrInvalidFollowRequest,
		ErrFollowRequestExists,
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
		if errors.Is(err, common.ErrValidation) {
			return http.StatusBadRequest
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
		NewSendFollowRequestEndpoint(svc),
		decodeSendFollowRequestRequest,
		encodeResponse, opts...,
	)
	listFriendReqHandler := kithttp.NewServer(
		NewListFollowRequestsRequestEndpoint(svc),
		decodeListFollowRequestsRequest,
		encodeResponse, opts...,
	)
	acceptFriendReqHandler := kithttp.NewServer(
		NewAcceptFollowRequestEndpoint(svc),
		decodeAcceptFollowRequestRequest,
		encodeResponse, opts...,
	)

	deleteFollowingReqHandler := kithttp.NewServer(
		NewDeleteFollowRequestEndpoint(svc),
		decodeDeleteFollowRequestRequest,
		encodeResponse, opts...,
	)

	getFollowingHandler := kithttp.NewServer(
		NewIsFollowingResponseEndpoint(svc),
		decodeIsFollowingRequest,
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

	deleteFollowHandler := kithttp.NewServer(
		NewDeleteFollowEndpoint(svc),
		decodeDeleteFollowRequest,
		encodeResponse, opts...,
	)

	ListTripPublicInfo := kithttp.NewServer(
		NewListTripPublicInfoEndpoint(svc),
		decodeListTripPublicInfoRequest,
		encodeResponse, opts...,
	)

	readTripPublicInfoHandler := kithttp.NewServer(
		NewReadTripPublicInfoEndpoint(svc),
		decodeReadTripPublicInfoRequest,
		encodeResponse, opts...,
	)

	duplicateTripHandler := kithttp.NewServer(
		NewDuplicateTripEndpoint(svc),
		decodeDuplicateTripRequest,
		encodeResponse, opts...,
	)

	r.Handle("/api/v1/social/{uid}", listFollowingTripsHandler).Methods(http.MethodGet)
	r.Handle("/api/v1/social/{uid}/profile", getProfileHandler).Methods(http.MethodGet)

	r.Handle("/api/v1/social/{uid}/requests", sendFriendReqHandler).Methods(http.MethodPost)
	r.Handle("/api/v1/social/{uid}/requests", listFriendReqHandler).Methods(http.MethodGet)
	r.Handle("/api/v1/social/{uid}/requests/{rid}/accept", acceptFriendReqHandler).Methods(http.MethodPut)
	r.Handle("/api/v1/social/{uid}/requests/{rid}", deleteFollowingReqHandler).Methods(http.MethodDelete)

	r.Handle("/api/v1/social/{uid}/followers", listFollowersHandler).Methods(http.MethodGet)
	r.Handle("/api/v1/social/{uid}/following", listFollowingHandler).Methods(http.MethodGet)
	r.Handle("/api/v1/social/{uid}/following/{targetID}", getFollowingHandler).Methods(http.MethodGet)
	r.Handle("/api/v1/social/{uid}/following/{bindingKey}", deleteFollowHandler).Methods(http.MethodDelete)

	r.Handle("/api/v1/social/{uid}/trips", ListTripPublicInfo).Methods(http.MethodGet)
	r.Handle("/api/v1/social/{uid}/trips/{tid}", readTripPublicInfoHandler).Methods(http.MethodGet)
	r.Handle("/api/v1/social/{uid}/trips/{tid}/duplicate", duplicateTripHandler).Methods(http.MethodPost)

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

func decodeSendFollowRequestRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	userID, ok := vars[URLPathVarUserID]
	if !ok {
		return nil, common.ErrInvalidRequest
	}
	req := SendFollowRequestRequest{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, common.ErrInvalidJSONBody
	}
	req.InitiatorID = userID
	return req, nil
}

func decodeAcceptFollowRequestRequest(_ context.Context, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	userID, ok := vars[URLPathVarUserID]
	if !ok {
		return nil, common.ErrInvalidRequest
	}
	requestID, ok := vars[URLPathVarRequestID]
	if !ok {
		return nil, common.ErrInvalidRequest
	}
	req := AcceptFollowRequestRequest{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, common.ErrInvalidJSONBody
	}

	req.UserID = userID
	req.RequestID = requestID
	return req, nil
}

func decodeDeleteFollowRequestRequest(_ context.Context, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	userID, ok := vars[URLPathVarUserID]
	if !ok {
		return nil, common.ErrInvalidRequest
	}
	requestID, ok := vars[URLPathVarRequestID]
	if !ok {
		return nil, common.ErrInvalidRequest
	}
	req := DeleteFollowRequestRequest{UserID: userID, RequestID: requestID}
	return req, nil
}

func decodeListFollowRequestsRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	userID, ok := vars[URLPathVarUserID]
	if !ok {
		return nil, common.ErrInvalidRequest
	}
	req := ListFollowRequestsRequest{
		ListFollowRequestsFilter: ListFollowRequestsFilter{
			TargetID: common.StringPtr(userID),
		},
	}
	return req, nil
}

func decodeIsFollowingRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	userID, ok := vars[URLPathVarUserID]
	if !ok {
		return nil, common.ErrInvalidRequest
	}
	targetID, ok := vars["targetID"]
	if !ok {
		return nil, common.ErrInvalidRequest
	}

	req := IsFollowingRequest{
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

func decodeDeleteFollowRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	userID, ok := vars[URLPathVarUserID]
	if !ok {
		return nil, common.ErrInvalidRequest
	}
	bindingKey, ok := vars[URLPathBindingKey]
	if !ok {
		return nil, common.ErrInvalidRequest
	}
	req := DeleteFollowRequest{
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

func decodeDuplicateTripRequest(_ context.Context, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	userID, ok := vars[URLPathVarUserID]
	if !ok {
		return nil, common.ErrInvalidRequest
	}
	tripID, ok := vars[URLPathVarTripID]
	if !ok {
		return nil, common.ErrInvalidRequest
	}

	req := DuplicateRequest{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, common.ErrInvalidJSONBody
	}
	req.TripID = tripID
	req.ReferrerID = userID
	return req, nil
}
