package trips

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/awhdesmond/tiinyplanet/pkg/common"
	"github.com/awhdesmond/tiinyplanet/pkg/utils"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
)

const (
	URLPathVarID = "id"
)

var (
	encodeErrFn = utils.EncodeErrorFactory(ErrorToHTTPCode)

	opts = []kithttp.ServerOption{
		// kithttp.ServerBefore(reqctx.MakeContextFromHTTPRequest),
		kithttp.ServerErrorEncoder(encodeErrFn),
	}
)

func encodeResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	if e, ok := response.(common.Errorer); ok && e.Error() != nil {
		encodeErrFn(ctx, e.Error(), w)
		return nil
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	return json.NewEncoder(w).Encode(response)
}

func MakeHandler(svc Service) http.Handler {
	r := mux.NewRouter()
	createTripPlanHandler := kithttp.NewServer(NewCreateTripPlanEndpoint(svc), decodeCreateTripPlanRequest, encodeResponse, opts...)
	listTripPlansHandler := kithttp.NewServer(NewListTripPlansEndpoint(svc), decodeListTripPlansRequest, encodeResponse, opts...)
	readTripPlanHandler := kithttp.NewServer(NewReadTripPlanEndpoint(svc), decodeReadTripPlanRequest, encodeResponse, opts...)
	deleteTripPlanHandler := kithttp.NewServer(NewDeleteTripPlanEndpoint(svc), decodeDeleteTripPlanRequest, encodeResponse, opts...)

	r.Handle("/api/v1/trips", createTripPlanHandler).Methods(http.MethodPost)
	r.Handle("/api/v1/trips", listTripPlansHandler).Methods(http.MethodGet)
	r.Handle("/api/v1/trips/{id}", readTripPlanHandler).Methods(http.MethodGet)
	r.Handle("/api/v1/trips/{id}", deleteTripPlanHandler).Methods(http.MethodDelete)

	return r
}

func decodeCreateTripPlanRequest(_ context.Context, r *http.Request) (interface{}, error) {
	req := CreateTripPlanRequest{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, utils.ErrInvalidJSONBody
	}
	return req, nil
}
func decodeReadTripPlanRequest(_ context.Context, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	ID, ok := vars[URLPathVarID]
	if !ok {
		return nil, utils.ErrInvalidRequest
	}
	return ReadTripPlanRequest{ID}, nil
}
func decodeListTripPlansRequest(_ context.Context, r *http.Request) (interface{}, error) {
	return ListTripPlansRequest{}, nil

}
func decodeDeleteTripPlanRequest(_ context.Context, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	ID, ok := vars[URLPathVarID]
	if !ok {
		return nil, utils.ErrInvalidRequest
	}
	return DeleteTripPlanRequest{ID}, nil
}
