package reqctx

import (
	"context"
	"net/http"

	"github.com/travelreys/travelreys/pkg/common"
)

const (
	AuthHeader     = "authorization"
	AuthCookieName = "_travelreysAuth"
)

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

	cookie, err := r.Cookie(AuthCookieName)
	if err != nil || cookie.Value == "" {
		return ci
	}

	jwtToken := cookie.Value
	ci.RawToken = jwtToken
	claims, err := common.ParseJWT(jwtToken, common.GetJwtSecret())
	if err != nil {
		return ci
	}
	if id, ok := claims[common.JwtClaimSub].(string); ok {
		ci.UserID = id
	} else {
		return ci
	}
	if email, ok := claims[common.JwtClaimEmail].(string); ok {
		ci.UserEmail = email
	} else {
		return ci
	}
	return ci
}

func ContextWithClientInfo(ctx context.Context, r *http.Request) context.Context {
	ci := makeClientInfo(r)
	return context.WithValue(ctx, common.ContextKeyClientInfo, ci)
}

func ClientInfoFromCtx(ctx context.Context) (ClientInfo, error) {
	ci, ok := ctx.Value(common.ContextKeyClientInfo).(ClientInfo)
	if !ok {
		return ci, common.ErrMissingJWTClaims
	}
	return ci, nil
}
