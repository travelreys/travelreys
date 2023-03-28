package auth

import (
	"context"
	"errors"
	"net/http"

	"github.com/travelreys/travelreys/pkg/reqctx"
	"go.uber.org/zap"
)

var (
	ErrRBAC = errors.New("auth.rbac.error")
)

type rbacMiddleware struct {
	next   Service
	logger *zap.Logger
}

func ServiceWithRBACMiddleware(svc Service, logger *zap.Logger) Service {
	return &rbacMiddleware{svc, logger}
}

func (mw rbacMiddleware) Login(ctx context.Context, authCode, provider string) (User, *http.Cookie, error) {
	return mw.next.Login(ctx, authCode, provider)
}

func (mw rbacMiddleware) Read(ctx context.Context, ID string) (User, error) {
	ci, err := reqctx.ClientInfoFromCtx(ctx)
	if err != nil || ci.HasEmptyID() {
		return User{}, ErrRBAC
	}
	return mw.next.Read(ctx, ID)
}

func (mw rbacMiddleware) List(ctx context.Context, ff ListFilter) (UsersList, error) {
	ci, err := reqctx.ClientInfoFromCtx(ctx)
	if err != nil || ci.HasEmptyID() {
		return UsersList{}, ErrRBAC
	}
	return mw.next.List(ctx, ff)
}

func (mw rbacMiddleware) Update(ctx context.Context, ID string, ff UpdateFilter) error {
	ci, err := reqctx.ClientInfoFromCtx(ctx)
	if err != nil {
		return ErrRBAC
	}
	if ci.UserID != ID {
		return ErrRBAC
	}
	return mw.next.Update(ctx, ID, ff)
}
