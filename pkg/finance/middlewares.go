package finance

import (
	"context"
	"errors"

	"github.com/travelreys/travelreys/pkg/reqctx"
	"go.uber.org/zap"
)

var (
	ErrRBAC = errors.New("finance.ErrRBAC")
)

type rbacMiddleware struct {
	next   Service
	logger *zap.Logger
}

func SvcWithRBACMw(svc Service, logger *zap.Logger) Service {
	return &rbacMiddleware{svc, logger}
}

func (mw rbacMiddleware) GetFxRates(ctx context.Context, base string) (ExchangeRates, error) {
	ci, err := reqctx.ClientInfoFromCtx(ctx)
	if err != nil || ci.HasEmptyID() {
		return ExchangeRates{}, ErrRBAC
	}
	return mw.next.GetFxRates(ctx, base)
}
