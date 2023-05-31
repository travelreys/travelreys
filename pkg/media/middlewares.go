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

func (mw rbacMiddleware) GenerateGetSignedURLs(ctx context.Context, items MediaItemList) (MediaPresignedUrlList, error) {
	ci, err := reqctx.ClientInfoFromCtx(ctx)
	if err != nil || ci.HasEmptyID() {
		return nil, ErrRBAC
	}
	return mw.next.GenerateGetSignedURLs(ctx, items)
}

func (mw rbacMiddleware) GeneratePutSignedURLs(ctx context.Context, items MediaItemList) (MediaPresignedUrlList, error) {
	ci, err := reqctx.ClientInfoFromCtx(ctx)
	if err != nil || ci.HasEmptyID() {
		return nil, ErrRBAC
	}
	return mw.next.GeneratePutSignedURLs(ctx, items)
}

func (mw rbacMiddleware) Save(ctx context.Context, items MediaItemList) error {
	ci, err := reqctx.ClientInfoFromCtx(ctx)
	if err != nil || ci.HasEmptyID() {
		return ErrRBAC
	}
	return mw.next.Save(ctx, items)
}

func (mw rbacMiddleware) Delete(ctx context.Context, items MediaItemList) error {
	ci, err := reqctx.ClientInfoFromCtx(ctx)
	if err != nil || ci.HasEmptyID() {
		return ErrRBAC
	}
	return mw.next.Delete(ctx, items)
}
