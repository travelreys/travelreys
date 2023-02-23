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

	createHandler := kithttp.NewServer(NewCreateEndpoint(svc), decodeCreateRequest, encodeResponse, opts...)
	listHandler := kithttp.NewServer(NewListEndpoint(svc), decodeListRequest, encodeResponse, opts...)
	readHandler := kithttp.NewServer(NewReadEndpoint(svc), decodeReadRequest, encodeResponse, opts...)
	readMembersHandler := kithttp.NewServer(NewReadMembersEndpoint(svc), decodeReadMembersRequest, encodeResponse, opts...)
	deleteHandler := kithttp.NewServer(NewDeleteEndpoint(svc), decodeDeleteRequest, encodeResponse, opts...)

	r.Handle("/api/v1/trips", createHandler).Methods(http.MethodPost)
	r.Handle("/api/v1/trips", listHandler).Methods(http.MethodGet)
	r.Handle("/api/v1/trips/{id}", readHandler).Methods(http.MethodGet)
	r.Handle("/api/v1/trips/{id}/members", readMembersHandler).Methods(http.MethodGet)
	r.Handle("/api/v1/trips/{id}", deleteHandler).Methods(http.MethodDelete)

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
	if r.URL.Query().Get("withMembers") == "true" {
		req.WithMembers = true
	}
	return req, nil
}

func decodeReadMembersRequest(_ context.Context, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	ID, ok := vars[URLPathVarID]
	if !ok {
		return nil, common.ErrInvalidRequest
	}
	return ReadMembersRequest{ID: ID}, nil
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
