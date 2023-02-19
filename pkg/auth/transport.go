package auth

import (
	"context"
	"encoding/json"
	"net/http"

	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"

	"github.com/tiinyplanet/tiinyplanet/pkg/common"
)

func errToHttpCode() func(err error) int {
	notFoundErrors := []error{ErrUserNotFound}
	appErrors := []error{ErrProviderNotSupported, ErrProviderGoogleError}
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

func encodeResponse(ctx context.Context, w http.ResponseWriter, resp interface{}) error {
	if e, ok := resp.(common.Errorer); ok && e.Error() != nil {
		common.EncodeErrorFactory(errToHttpCode())(ctx, e.Error(), w)
		return nil
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	return json.NewEncoder(w).Encode(resp)
}

func MakeHandler(svc Service) http.Handler {
	r := mux.NewRouter()

	opts := []kithttp.ServerOption{
		kithttp.ServerBefore(common.AddClientInfoToCtx),
		kithttp.ServerErrorEncoder(common.EncodeErrorFactory(errToHttpCode())),
	}

	loginHandler := kithttp.NewServer(NewLoginEndpoint(svc), decodeLoginRequest, encodeResponse, opts...)
	readUserHandler := kithttp.NewServer(NewReadUserEndpoint(svc), decodeReadUserRequest, encodeResponse, opts...)
	listUsersHandler := kithttp.NewServer(NewListUsersEndpoint(svc), decodeListUsersRequest, encodeResponse, opts...)
	updateUserHandler := kithttp.NewServer(NewUpdateUserEndpoint(svc), decodeUpdateUserRequest, encodeResponse, opts...)

	r.Handle("/api/v1/auth/login", loginHandler).Methods(http.MethodPost)
	r.Handle("/api/v1/auth/users", listUsersHandler).Methods(http.MethodGet)
	r.Handle("/api/v1/auth/users/{id}", readUserHandler).Methods(http.MethodGet)
	r.Handle("/api/v1/auth/users/{id}", updateUserHandler).Methods(http.MethodPut)

	return r
}

// Request Decoders

func decodeLoginRequest(_ context.Context, r *http.Request) (interface{}, error) {
	req := LoginRequest{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, common.ErrInvalidJSONBody
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
		return nil, common.ErrInvalidJSONBody
	}
	return req, nil
}

func decodeListUsersRequest(_ context.Context, r *http.Request) (interface{}, error) {
	ff := ListUsersFilter{}
	req := ListUsersRequest{ff}
	if r.URL.Query().Has("emai") {
		ff.Email = common.StringPtr(r.URL.Query().Get("emai"))
	}
	return req, nil
}
