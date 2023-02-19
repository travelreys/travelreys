package common

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"os"
)

const (
	AuthHeader   = "authorization"
	EnvJwtSecret = "TIINYPLANET_JWT_SECRET"
)

var (
	ErrMissingAuthHeader    = errors.New("missing-auth-header")
	ErrInvalidAuthToken     = errors.New("invalid-auth-token")
	ErrInvalidSigningMethod = errors.New("invalid-signing-method")
	ErrMissingJWTClaims     = errors.New("missing-jwt-claims")
)

// Client Info

type ClientInfo struct {
	RawToken  string
	UserID    string
	UserEmail string
}

func (ci ClientInfo) HasEmptyID() bool {
	return ci.UserID == ""
}

func makeClientInfo(r *http.Request) ClientInfo {
	ci := ClientInfo{}

	authorization := r.Header.Get(AuthHeader)
	if authorization == "" {
		return ci
	}
	jwtToken, err := ParseBearerAndToken(authorization)
	if err != nil {
		return ci
	}
	ci.RawToken = jwtToken
	claims, err := ParseJWT(jwtToken, os.Getenv(EnvJwtSecret))
	if err != nil {
		return ci
	}
	if id, ok := claims[JWTClaimSub].(string); ok {
		ci.UserID = id
	} else {
		return ci
	}
	if email, ok := claims[JWTClaimEmail].(string); ok {
		ci.UserEmail = email
	} else {
		return ci
	}
	return ci
}

func AddClientInfoToCtx(ctx context.Context, r *http.Request) context.Context {
	ci := makeClientInfo(r)
	return context.WithValue(ctx, ContextKeyClientInfo, ci)
}

func ReadClientInfoFromCtx(ctx context.Context) (ClientInfo, error) {
	ci, ok := ctx.Value(ContextKeyClientInfo).(ClientInfo)
	if !ok {
		return ci, ErrMissingJWTClaims
	}
	return ci, nil
}

// Transport Errors

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
