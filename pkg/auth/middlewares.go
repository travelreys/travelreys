package auth

import (
	"context"
	"errors"

	"github.com/tiinyplanet/tiinyplanet/pkg/common"
	"go.uber.org/zap"
)

var (
	ErrRBACMissing = errors.New("auth.rbac.missing")
	ErrRBAC        = errors.New("auth.rbac.error")
)

type rbacMiddleware struct {
	next   Service
	logger *zap.Logger
}

func ServiceWithRBACMiddleware(svc Service, logger *zap.Logger) Service {
	return &rbacMiddleware{svc, logger}
}

func (mw rbacMiddleware) Login(ctx context.Context, authCode, provider string) (string, error) {
	return mw.next.Login(ctx, authCode, provider)
}

func (mw rbacMiddleware) ReadUser(ctx context.Context, ID string) (User, error) {
	ci, err := common.ReadClientInfoFromCtx(ctx)
	if err != nil || ci.HasEmptyID() {
		return User{}, ErrRBACMissing
	}
	return mw.next.ReadUser(ctx, ID)
}

func (mw rbacMiddleware) ListUsers(ctx context.Context, ff ListUsersFilter) (UsersList, error) {
	ci, err := common.ReadClientInfoFromCtx(ctx)
	if err != nil || ci.HasEmptyID() {
		return UsersList{}, ErrRBACMissing
	}
	return mw.next.ListUsers(ctx, ff)
}

func (mw rbacMiddleware) UpdateUser(ctx context.Context, ID string, ff UpdateUserFilter) error {
	ci, err := common.ReadClientInfoFromCtx(ctx)
	if err != nil {
		return ErrRBACMissing
	}
	if ci.UserID != ID {
		return ErrRBAC
	}
	return mw.next.UpdateUser(ctx, ID, ff)
}
