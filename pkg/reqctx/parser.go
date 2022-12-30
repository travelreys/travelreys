package reqctx

import (
	"context"
	"errors"
	"net/http"
	"os"
)

const (
	AuthHeader = "authorization"
	JwtIssuer  = "tiinyPlanet"

	EnvJwtSecret = "JWT_SECRET"
)

var (
	ErrMissingAuthHeader    = errors.New("missing-auth-header")
	ErrInvalidAuthToken     = errors.New("invalid-auth-token")
	ErrInvalidSigningMethod = errors.New("invalid-signing-method")
	ErrMissingJWTClaims     = errors.New("missing-jwt-claims")
)

func processCallerInfo(rctx *Context, r *http.Request) {
	authorization := r.Header.Get(AuthHeader)
	if authorization == "" {
		rctx.CallerInfo.Err = ErrMissingAuthHeader
		return
	}

	jwtToken, err := ParseBearerAndToken(authorization)
	if err != nil {
		rctx.CallerInfo.Err = err
		return
	}
	rctx.CallerInfo.RawToken = jwtToken
	claims, err := ParseJWT(jwtToken, os.Getenv(EnvJwtSecret))
	if err != nil {
		rctx.CallerInfo.Err = err
		return
	}
	if email, ok := claims["email"].(string); ok {
		rctx.CallerInfo.UserEmail = email
	} else {
		rctx.CallerInfo.Err = ErrMissingJWTClaims
		return
	}
	if id, ok := claims["id"].(string); ok {
		rctx.CallerInfo.UserID = id
	} else {
		rctx.CallerInfo.Err = ErrMissingJWTClaims
		return
	}
}

func MakeContextFromHTTPRequest(ctx context.Context, r *http.Request) context.Context {
	rctx := NewContextFromContext(ctx)
	processCallerInfo(&rctx, r)
	return rctx
}
