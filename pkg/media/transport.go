package media

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"

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
	notFoundErrors := []error{}
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
	w.Header().Set("Cache-Control", "no-store")
	return json.NewEncoder(w).Encode(response)
}

func MakeHandler(svc Service) http.Handler {
	r := mux.NewRouter()

	opts := []kithttp.ServerOption{
		kithttp.ServerBefore(reqctx.ContextWithClientInfo),
		kithttp.ServerErrorEncoder(common.EncodeErrorFactory(errToHttpCode)),
	}
	generateMediaItemsHandler := kithttp.NewServer(
		NewGenerateMediaItemsEndpoint(svc), decodeGenerateMediaItemsRequest, encodeResponse, opts...,
	)
	saveForUserHandler := kithttp.NewServer(
		NewSaveForUserEndpoint(svc), decodeSaveForUserRequest, encodeResponse, opts...,
	)

	listHandler := kithttp.NewServer(
		NewListEndpoint(svc), decodeListRequest, encodeResponse, opts...,
	)

	deleteHandler := kithttp.NewServer(
		NewDeleteEndpoint(svc), decodeDeleteRequest, encodeResponse, opts...,
	)

	r.Handle("/api/v1/media", listHandler).Methods(http.MethodGet)
	r.Handle("/api/v1/media/user/{userid}/generate", generateMediaItemsHandler).Methods(http.MethodPost)
	r.Handle("/api/v1/media/user/{userid}", saveForUserHandler).Methods(http.MethodPost)
	r.Handle("/api/v1/media/user/{userid}", deleteHandler).Methods(http.MethodDelete)

	return r
}

func decodeGenerateMediaItemsRequest(_ context.Context, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	userID, ok := vars[URLPathVarUserID]
	if !ok {
		return nil, common.ErrInvalidRequest
	}

	req := GenerateMediaItemsRequest{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, common.ErrInvalidJSONBody
	}
	req.UserID = userID
	return req, nil
}

func decodeSaveForUserRequest(_ context.Context, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	userID, ok := vars[URLPathVarUserID]
	if !ok {
		return nil, common.ErrInvalidRequest
	}

	req := SaveForUserRequest{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, common.ErrInvalidJSONBody
	}
	req.UserID = userID
	return req, nil
}

func decodeListRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	req := ListRequest{ListMediaFilter: ListMediaFilter{}}
	params := r.URL.Query()
	if params.Has(URLQueryVarUserID) {
		req.UserID = common.StringPtr(params.Get(URLQueryVarUserID))
	}
	if params.Has("withURLs") {
		req.WithURLs = true
	}
	return req, nil
}

func decodeDeleteRequest(_ context.Context, r *http.Request) (interface{}, error) {
	params := r.URL.Query()
	req := DeleteRequest{
		DeleteMediaFilter: DeleteMediaFilter{
			UserID: common.StringPtr(params.Get(URLQueryVarUserID)),
			IDs:    strings.Split(params.Get("ids"), ","),
		},
	}
	return req, nil
}
