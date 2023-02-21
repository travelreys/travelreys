package images

import (
	"context"
	"errors"

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

func (mw rbacMiddleware) Search(ctx context.Context, query string) (MetadataList, error) {
	ci, err := reqctx.ClientInfoFromCtx(ctx)
	if err != nil || ci.HasEmptyID() {
		return MetadataList{}, ErrRBAC
	}
	return mw.next.Search(ctx, query)
}
