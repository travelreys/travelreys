package auth

import (
	"context"
	"encoding/json"
	"net/http"

	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"

	"github.com/tiinyplanet/tiinyplanet/pkg/common"
	"github.com/tiinyplanet/tiinyplanet/pkg/utils"
)

var (
	notFoundErrors = []error{}
	appErrors      = []error{}
	authErrors     = []error{}
)

var (
	encodeErrFn = utils.EncodeErrorFactory(
		utils.ErrorToHTTPCodeFactory(notFoundErrors, appErrors, authErrors),
	)
	opts = []kithttp.ServerOption{
		// kithttp.ServerBefore(reqctx.MakeContextFromHTTPRequest),
		kithttp.ServerErrorEncoder(encodeErrFn),
	}
)

func encodeResponse(ctx context.Context, w http.ResponseWriter, resp interface{}) error {
	if e, ok := resp.(common.Errorer); ok && e.Error() != nil {
		encodeErrFn(ctx, e.Error(), w)
		return nil
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	return json.NewEncoder(w).Encode(resp)
}

func MakeHandler(svc Service) http.Handler {
	r := mux.NewRouter()

	loginHandler := kithttp.NewServer(NewLoginEndpoint(svc), decodeLoginRequest, encodeResponse, opts...)
	readUserHandler := kithttp.NewServer(NewReadUserEndpoint(svc), decodeReadUserRequest, encodeResponse)
	updateUserHandler := kithttp.NewServer(NewUpdateUserEndpoint(svc), decodeUpdateUserRequest, encodeResponse)

	r.Handle("/api/v1/auth/login", loginHandler).Methods(http.MethodPost)
	r.Handle("/api/v1/auth/user/{id}", readUserHandler).Methods(http.MethodGet)
	r.Handle("/api/v1/auth/user/{id}", updateUserHandler).Methods(http.MethodPut)

	return r
}

// Request Decoders

func decodeLoginRequest(_ context.Context, r *http.Request) (interface{}, error) {
	req := LoginRequest{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, utils.ErrInvalidJSONBody
	}
	return req, nil
}

func decodeReadUserRequest(_ context.Context, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	req := ReadUserRequest{
		ID: vars["id"],
	}
	return req, nil
}

func decodeUpdateUserRequest(_ context.Context, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	req := UpdateUserRequest{
		ID: vars["id"],
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, utils.ErrInvalidJSONBody
	}
	return req, nil
}
