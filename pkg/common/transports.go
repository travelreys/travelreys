package common

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
)

var (
	ErrBadPath                 = errors.New("http.badpath")
	ErrInvalidRequest          = errors.New("http.invalidrequest")
	ErrInvalidJSONBody         = errors.New("http.invalidjson")
	ErrMissingAuthHeader       = errors.New("jwt.missing-auth-header")
	ErrInvalidAuthToken        = errors.New("jwt.invalid-auth-token")
	ErrInvalidSigningMethod    = errors.New("jwt.invalid-signing-method")
	ErrMissingJWTClaims        = errors.New("jwt.missing-claims")
	ErrorEndpointReqMismatch   = errors.New("endpoint.invalidrequest")
	ErrorEncodeInvalidResponse = errors.New("encode.invalidresponse")
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
