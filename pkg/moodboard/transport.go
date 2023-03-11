package moodboard

import (
	"context"
	"encoding/json"
	"net/http"

	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	"github.com/travelreys/travelreys/pkg/common"
	"github.com/travelreys/travelreys/pkg/reqctx"
)

const (
	URLPathVarPinID = "pid"
)

func errToHttpCode() func(err error) int {
	notFoundErrors := []error{ErrMoodboardNotFound}
	appErrors := []error{ErrUnexpectedStoreError}
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

	readHandler := kithttp.NewServer(NewReadAndCreateIfNotExistsEndpoint(svc), decodeReadRequest, encodeResponse, opts...)
	updateHandler := kithttp.NewServer(NewUpdateEndpoint(svc), decodeUpdateRequest, encodeResponse, opts...)
	addPinHandler := kithttp.NewServer(NewAddPinEndpoint(svc), decodeAddPinRequest, encodeResponse, opts...)
	updatePinHandler := kithttp.NewServer(NewUpdatePinEndpoint(svc), decodeUpdatePinRequest, encodeResponse, opts...)
	deletePinHandler := kithttp.NewServer(NewDeletePinEndpoint(svc), decodeDeletePinRequest, encodeResponse, opts...)

	r.Handle("/api/v1/moodboards", readHandler).Methods(http.MethodGet)
	r.Handle("/api/v1/moodboards", updateHandler).Methods(http.MethodPut)
	r.Handle("/api/v1/moodboards/pins", addPinHandler).Methods(http.MethodPost)
	r.Handle("/api/v1/moodboards/pins/{pid}", updatePinHandler).Methods(http.MethodPut)
	r.Handle("/api/v1/moodboards/pins/{pid}", deletePinHandler).Methods(http.MethodDelete)

	return r
}

func decodeReadRequest(_ context.Context, r *http.Request) (interface{}, error) {
	return ReadAndCreateIfNotExistsRequest{}, nil
}

func decodeUpdateRequest(_ context.Context, r *http.Request) (interface{}, error) {
	req := UpdateRequest{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, common.ErrInvalidJSONBody
	}
	return req, nil
}

func decodeAddPinRequest(_ context.Context, r *http.Request) (interface{}, error) {
	req := AddPinRequest{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, common.ErrInvalidJSONBody
	}
	return req, nil
}

func decodeUpdatePinRequest(_ context.Context, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	ID, ok := vars[URLPathVarPinID]
	if !ok {
		return nil, common.ErrInvalidRequest
	}
	req := UpdatePinRequest{PinID: ID}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, common.ErrInvalidJSONBody
	}
	return req, nil
}

func decodeDeletePinRequest(_ context.Context, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	ID, ok := vars[URLPathVarPinID]
	if !ok {
		return nil, common.ErrInvalidRequest
	}
	return DeletePinRequest{PinID: ID}, nil
}
