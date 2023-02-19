package images

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

func (mw rbacMiddleware) Search(ctx context.Context, query string) (ImageMetadataList, error) {
	ci, err := common.ReadClientInfoFromCtx(ctx)
	if err != nil || ci.HasEmptyID() {
		return ImageMetadataList{}, ErrRBACMissing
	}
	return mw.next.Search(ctx, query)
}
