package footprints

import (
	"context"
	"errors"

	"github.com/travelreys/travelreys/pkg/reqctx"
	"github.com/travelreys/travelreys/pkg/trips"
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

func (mw rbacMiddleware) CheckIn(ctx context.Context, userID, tripID string, activity trips.Activity) error {
	ci, err := reqctx.ClientInfoFromCtx(ctx)
	if err != nil || ci.HasEmptyID() {
		return ErrRBAC
	}
	return mw.next.CheckIn(ctx, userID, tripID, activity)
}

func (mw rbacMiddleware) List(ctx context.Context, ff ListFootprintsFilter) (FootprintList, error) {
	return mw.next.List(ctx, ff)
}
