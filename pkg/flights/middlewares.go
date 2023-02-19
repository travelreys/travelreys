package flights

import (
	"context"
	"errors"
	"time"

	"github.com/tiinyplanet/tiinyplanet/pkg/common"
	"go.uber.org/zap"
)

var (
	ErrRBAC        = errors.New("flights.rbac.error")
	ErrRBACMissing = errors.New("auth.rbac.missing")
)

type rbacMiddleware struct {
	next   Service
	logger *zap.Logger
}

func ServiceWithRBACMiddleware(svc Service, logger *zap.Logger) Service {
	return &rbacMiddleware{svc, logger}
}

func (mw rbacMiddleware) Search(
	ctx context.Context,
	origIATA,
	destIATA string,
	numAdults uint64,
	departDate time.Time,
	opts FlightsSearchOptions,
) (Itineraries, error) {
	ci, err := common.ReadClientInfoFromCtx(ctx)
	if err != nil || ci.HasEmptyID() {
		return Itineraries{}, ErrRBACMissing
	}
	return mw.next.Search(ctx, origIATA, destIATA, numAdults, departDate, opts)
}
