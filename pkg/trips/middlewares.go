package trips

import (
	context "context"
	"errors"
	"time"

	"github.com/tiinyplanet/tiinyplanet/pkg/auth"
	"github.com/tiinyplanet/tiinyplanet/pkg/reqctx"
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

func (mw rbacMiddleware) Create(ctx context.Context, creator Member, name string, start, end time.Time) (Trip, error) {
	ci, err := reqctx.ClientInfoFromCtx(ctx)
	if err != nil || ci.HasEmptyID() {
		return Trip{}, ErrRBAC
	}
	return mw.next.Create(ctx, creator, name, start, end)
}

func (mw rbacMiddleware) Read(ctx context.Context, ID string) (Trip, error) {
	ci, err := reqctx.ClientInfoFromCtx(ctx)
	if err != nil || ci.HasEmptyID() {
		return Trip{}, ErrRBAC
	}
	return mw.next.Read(ctx, ID)
}

func (mw rbacMiddleware) ReadWithUsers(ctx context.Context, ID string) (Trip, auth.UsersMap, error) {
	ci, err := reqctx.ClientInfoFromCtx(ctx)
	if err != nil || ci.HasEmptyID() {
		return Trip{}, nil, ErrRBAC
	}
	return mw.next.ReadWithUsers(ctx, ID)
}

func (mw rbacMiddleware) List(ctx context.Context, ff ListFilter) (TripsList, error) {
	ci, err := reqctx.ClientInfoFromCtx(ctx)
	if err != nil || ci.HasEmptyID() {
		return TripsList{}, ErrRBAC
	}
	return mw.next.List(ctx, ff)
}

func (mw rbacMiddleware) Delete(ctx context.Context, ID string) error {
	ci, err := reqctx.ClientInfoFromCtx(ctx)
	if err != nil || ci.HasEmptyID() {
		return ErrRBAC
	}
	return mw.next.Delete(ctx, ID)
}
