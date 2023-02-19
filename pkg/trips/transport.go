package trips

import (
	"context"
	"encoding/json"
	"net/http"

	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	"github.com/tiinyplanet/tiinyplanet/pkg/common"
)

const (
	URLPathVarID = "id"
)

func errToHttpCode() func(err error) int {
	notFoundErrors := []error{ErrPlanNotFound}
	appErrors := []error{ErrUnexpectedStoreError}
	authErrors := []error{ErrRBAC, ErrRBACMissing}

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
		kithttp.ServerBefore(common.AddClientInfoToCtx),
		kithttp.ServerErrorEncoder(common.EncodeErrorFactory(errToHttpCode())),
	}

	createTripPlanHandler := kithttp.NewServer(NewCreateTripEndpoint(svc), decodeCreateTripRequest, encodeResponse, opts...)
	listTripPlansHandler := kithttp.NewServer(NewListTripsEndpoint(svc), decodeListTripsRequest, encodeResponse, opts...)
	readTripPlanHandler := kithttp.NewServer(NewReadTripEndpoint(svc), decodeReadTripRequest, encodeResponse, opts...)
	deleteTripPlanHandler := kithttp.NewServer(NewDeleteTripEndpoint(svc), decodeDeleteTripRequest, encodeResponse, opts...)

	r.Handle("/api/v1/trips", createTripPlanHandler).Methods(http.MethodPost)
	r.Handle("/api/v1/trips", listTripPlansHandler).Methods(http.MethodGet)
	r.Handle("/api/v1/trips/{id}", readTripPlanHandler).Methods(http.MethodGet)
	r.Handle("/api/v1/trips/{id}", deleteTripPlanHandler).Methods(http.MethodDelete)

	return r
}

func decodeCreateTripRequest(_ context.Context, r *http.Request) (interface{}, error) {
	req := CreateTripRequest{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, common.ErrInvalidJSONBody
	}
	return req, nil
}
func decodeReadTripRequest(_ context.Context, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	ID, ok := vars[URLPathVarID]
	if !ok {
		return nil, common.ErrInvalidRequest
	}

	req := ReadTripRequest{ID: ID}
	if r.URL.Query().Get("withUsers") == "true" {
		req.WithUsers = true
	}

	return req, nil
}
func decodeListTripsRequest(_ context.Context, r *http.Request) (interface{}, error) {
	return ListTripsRequest{}, nil

}
func decodeDeleteTripRequest(_ context.Context, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	ID, ok := vars[URLPathVarID]
	if !ok {
		return nil, common.ErrInvalidRequest
	}
	return DeleteTripRequest{ID}, nil
}
