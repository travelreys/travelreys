package utils

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/tiinyplanet/tiinyplanet/pkg/common"
)

var (
	ErrBadPath         = errors.New("http-bad-path")
	ErrInvalidRequest  = errors.New("http-invalid-request")
	ErrInvalidJSONBody = errors.New("http-invalid-json-body")
)

func writeErrorHeader(_ context.Context, err error, w http.ResponseWriter) {
	switch err {
	case ErrInvalidRequest:
	case ErrInvalidJSONBody:
		w.WriteHeader(http.StatusBadRequest)
	case ErrBadPath:
		w.WriteHeader(http.StatusNotFound)
	default:
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func EncodeErrorFactory(errToCode func(error) int) func(context.Context, error, http.ResponseWriter) {
	return func(ctx context.Context, err error, w http.ResponseWriter) {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		code := errToCode(err)
		if code < 0 {
			writeErrorHeader(ctx, err, w)
		} else {
			w.WriteHeader(code)
		}
		json.NewEncoder(w).Encode(common.GenericJSON{"error": err.Error()})
	}
}

func ErrorToHTTPCodeFactory(notFoundErrors, appErrors, authErrors []error) func(err error) int {
	return func(err error) int {
		if ErrorContains(notFoundErrors, err) {
			return http.StatusNotFound
		}
		if ErrorContains(appErrors, err) {
			return http.StatusUnprocessableEntity
		}
		if ErrorContains(authErrors, err) {
			return http.StatusUnauthorized
		}
		return http.StatusInternalServerError

	}
}
