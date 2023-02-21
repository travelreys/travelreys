package trips

import (
	"context"
	"encoding/json"
	"net/http"

	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	"github.com/tiinyplanet/tiinyplanet/pkg/common"
	"github.com/tiinyplanet/tiinyplanet/pkg/reqctx"
)

const (
	URLPathVarID = "id"
)

func errToHttpCode(err error) int {
	notFoundErrors := []error{ErrPlanNotFound}
	appErrors := []error{ErrUnexpectedStoreError}
	authErrors := []error{ErrRBAC}

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

func encodeResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	if e, ok := response.(common.Errorer); ok && e.Error() != nil {
		common.EncodeErrorFactory(errToHttpCode)(ctx, e.Error(), w)
		return nil
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	return json.NewEncoder(w).Encode(response)
}

func MakeHandler(svc Service) http.Handler {
	r := mux.NewRouter()

	opts := []kithttp.ServerOption{
		kithttp.ServerBefore(reqctx.ContextWithClientInfo),
		kithttp.ServerErrorEncoder(common.EncodeErrorFactory(errToHttpCode)),
	}

	createTripHandler := kithttp.NewServer(NewCreateEndpoint(svc), decodeCreateRequest, encodeResponse, opts...)
	listTripsHandler := kithttp.NewServer(NewListEndpoint(svc), decodeListRequest, encodeResponse, opts...)
	readTripHandler := kithttp.NewServer(NewReadEndpoint(svc), decodeReadRequest, encodeResponse, opts...)
	deleteTripHandler := kithttp.NewServer(NewDeleteEndpoint(svc), decodeDeleteRequest, encodeResponse, opts...)

	r.Handle("/api/v1/trips", createTripHandler).Methods(http.MethodPost)
	r.Handle("/api/v1/trips", listTripsHandler).Methods(http.MethodGet)
	r.Handle("/api/v1/trips/{id}", readTripHandler).Methods(http.MethodGet)
	r.Handle("/api/v1/trips/{id}", deleteTripHandler).Methods(http.MethodDelete)

	return r
}

func decodeCreateRequest(_ context.Context, r *http.Request) (interface{}, error) {
	req := CreateRequest{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, common.ErrInvalidJSONBody
	}
	return req, nil
}
func decodeReadRequest(_ context.Context, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	ID, ok := vars[URLPathVarID]
	if !ok {
		return nil, common.ErrInvalidRequest
	}

	req := ReadRequest{ID: ID}
	if r.URL.Query().Get("withUsers") == "true" {
		req.WithUsers = true
	}

	return req, nil
}
func decodeListRequest(_ context.Context, r *http.Request) (interface{}, error) {
	return ListRequest{}, nil

}
func decodeDeleteRequest(_ context.Context, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	ID, ok := vars[URLPathVarID]
	if !ok {
		return nil, common.ErrInvalidRequest
	}
	return DeleteRequest{ID}, nil
}
