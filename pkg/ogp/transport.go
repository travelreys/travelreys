package ogp

import (
	"context"
	"encoding/json"
	"net/http"

	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	"github.com/travelreys/travelreys/pkg/common"
	"github.com/travelreys/travelreys/pkg/reqctx"
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

	fetchHandler := kithttp.NewServer(NewFetchEndpoint(svc), decodeFetchRequest, encodeResponse, opts...)

	r.Handle("/api/v1/ogp", fetchHandler).Methods(http.MethodGet)
	return r
}

func decodeFetchRequest(_ context.Context, r *http.Request) (interface{}, error) {
	q := r.URL.Query()
	return FetchRequest{
		Url: q.Get("url"),
	}, nil

}
