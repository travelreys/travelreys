package flights

import (
	"context"
	"encoding/json"
	"net/http"

	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/tiinyplanet/tiinyplanet/pkg/common"
	"github.com/tiinyplanet/tiinyplanet/pkg/reqctx"

	"github.com/gorilla/mux"
)

func errToHttpCode() func(err error) int {
	notFoundErrors := []error{}
	appErrors := []error{}
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
	searchHandler := kithttp.NewServer(
		NewSearchEndpoint(svc),
		decodeSearchRequest,
		encodeResponse,
		opts...,
	)
	r.Handle("/api/v1/flights/search", searchHandler).Methods(http.MethodGet)
	return r
}

func decodeSearchRequest(_ context.Context, r *http.Request) (interface{}, error) {
	q := r.URL.Query()
	opts, err := SearchOptionsFromURLValues(q)
	if err != nil {
		return nil, common.ErrInvalidRequest
	}
	return SearchRequest{opts}, nil
}
