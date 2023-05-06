package footprints

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
	URLPathVarUserID  = "userid"
	URLQueryVarUserID = "userid"
)

func errToHttpCode(err error) int {
	notFoundErrors := []error{ErrFpNotFound}
	appErrors := []error{
		ErrCheckinTripAlready,
		ErrUnexpectedStoreError,
	}
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
	w.Header().Set("Cache-Control", "no-store")
	return json.NewEncoder(w).Encode(response)
}

func MakeHandler(svc Service) http.Handler {
	r := mux.NewRouter()

	opts := []kithttp.ServerOption{
		kithttp.ServerBefore(reqctx.ContextWithClientInfo),
		kithttp.ServerErrorEncoder(common.EncodeErrorFactory(errToHttpCode)),
	}

	checkinHandler := kithttp.NewServer(
		NewCheckInEndpoint(svc), decodeCheckinRequest, encodeResponse, opts...,
	)
	listHandler := kithttp.NewServer(
		NewListEndpoint(svc), decodeListRequest, encodeResponse, opts...,
	)

	r.Handle("/api/v1/footprints", listHandler).Methods(http.MethodGet)
	r.Handle("/api/v1/footprints/checkin", checkinHandler).Methods(http.MethodPost)

	return r
}

func decodeCheckinRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	req := CheckInRequest{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, common.ErrInvalidJSONBody
	}
	return req, nil
}

func decodeListRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	req := ListRequest{ListFootprintsFilter: ListFootprintsFilter{}}
	params := r.URL.Query()
	if params.Has(URLQueryVarUserID) {
		req.UserID = common.StringPtr(params.Get(URLQueryVarUserID))
	}
	return req, nil
}
