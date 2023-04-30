package media

import (
	"context"
	"errors"

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

func (mw rbacMiddleware) GenerateMediaItems(ctx context.Context, userID string, params []NewMediaItemParams) (MediaItemList, []string, error) {
	ci, err := reqctx.ClientInfoFromCtx(ctx)
	if err != nil || ci.HasEmptyID() {
		return nil, nil, ErrRBAC
	}
	return mw.next.GenerateMediaItems(ctx, userID, params)
}

func (mw rbacMiddleware) SaveForUser(ctx context.Context, userID string, items MediaItemList) error {
	ci, err := reqctx.ClientInfoFromCtx(ctx)
	if err != nil || ci.HasEmptyID() {
		return ErrRBAC
	}
	return mw.next.SaveForUser(ctx, userID, items)
}

func (mw rbacMiddleware) List(ctx context.Context, ff ListMediaFilter, pg ListMediaPagination) (MediaItemList, string, error) {
	return mw.next.List(ctx, ff, pg)
}

func (mw rbacMiddleware) ListWithSignedURLs(ctx context.Context, ff ListMediaFilter, pg ListMediaPagination) (MediaItemList, string, []string, error) {
	return mw.next.ListWithSignedURLs(ctx, ff, pg)
}

func (mw rbacMiddleware) Delete(ctx context.Context, ff DeleteMediaFilter) error {
	ci, err := reqctx.ClientInfoFromCtx(ctx)
	if err != nil || ci.HasEmptyID() {
		return ErrRBAC
	}
	return mw.next.Delete(ctx, ff)
}

func (mw rbacMiddleware) GenerateGetSignedURLsForItems(ctx context.Context, items MediaItemList) ([]string, error) {
	ci, err := reqctx.ClientInfoFromCtx(ctx)
	if err != nil || ci.HasEmptyID() {
		return nil, ErrRBAC
	}
	return mw.next.GenerateGetSignedURLsForItems(ctx, items)
}
