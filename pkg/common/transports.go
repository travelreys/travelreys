package common

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
)

var (
	ErrBadPath                 = errors.New("common.ErrBadPath")
	ErrInvalidRequest          = errors.New("common.ErrInvalidRequest")
	ErrInvalidJSONBody         = errors.New("common.ErrInvalidJSONBody")
	ErrMissingAuthHeader       = errors.New("common.ErrMissingAuthHeader")
	ErrInvalidAuthToken        = errors.New("common.ErrInvalidAuthToken")
	ErrInvalidSigningMethod    = errors.New("common.ErrInvalidSigningMethod")
	ErrMissingJWTClaims        = errors.New("common.ErrMissingJWTClaims")
	ErrEndpointReqMismatch     = errors.New("common.ErrEndpointReqMismatch")
	ErrorEncodeInvalidResponse = errors.New("common.ErrorEncodeInvalidResponse")
)

func EncodeErrorFactory(errToCode func(error) int) func(context.Context, error, http.ResponseWriter) {
	return func(ctx context.Context, err error, w http.ResponseWriter) {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		code := errToCode(err)
		if code < 0 {
			switch err {
			case ErrInvalidRequest:
			case ErrInvalidJSONBody:
				w.WriteHeader(http.StatusBadRequest)
			case ErrBadPath:
				w.WriteHeader(http.StatusNotFound)
			default:
				w.WriteHeader(http.StatusInternalServerError)
			}
		} else {
			w.WriteHeader(code)
		}
		json.NewEncoder(w).Encode(GenericJSON{"error": err.Error()})
	}
}
