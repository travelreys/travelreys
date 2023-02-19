package trips

import (
	context "context"
	"errors"
	"time"

	"github.com/tiinyplanet/tiinyplanet/pkg/auth"
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

func (mw rbacMiddleware) CreateTrip(ctx context.Context, creator Member, name string, start, end time.Time) (TripPlan, error) {
	ci, err := common.ReadClientInfoFromCtx(ctx)
	if err != nil || ci.HasEmptyID() {
		return TripPlan{}, ErrRBACMissing
	}
	return mw.next.CreateTrip(ctx, creator, name, start, end)
}

func (mw rbacMiddleware) ReadTrip(ctx context.Context, ID string) (TripPlan, error) {
	ci, err := common.ReadClientInfoFromCtx(ctx)
	if err != nil || ci.HasEmptyID() {
		return TripPlan{}, ErrRBACMissing
	}
	return mw.next.ReadTrip(ctx, ID)
}

func (mw rbacMiddleware) ReadTripWithUsers(ctx context.Context, ID string) (TripPlan, auth.UsersMap, error) {
	ci, err := common.ReadClientInfoFromCtx(ctx)
	if err != nil || ci.HasEmptyID() {
		return TripPlan{}, nil, ErrRBACMissing
	}
	return mw.next.ReadTripWithUsers(ctx, ID)
}

func (mw rbacMiddleware) ListTrips(ctx context.Context, ff ListTripsFilter) (TripPlansList, error) {
	ci, err := common.ReadClientInfoFromCtx(ctx)
	if err != nil || ci.HasEmptyID() {
		return TripPlansList{}, ErrRBACMissing
	}
	return mw.next.ListTrips(ctx, ff)
}

func (mw rbacMiddleware) DeleteTrip(ctx context.Context, ID string) error {
	ci, err := common.ReadClientInfoFromCtx(ctx)
	if err != nil || ci.HasEmptyID() {
		return ErrRBACMissing
	}
	return mw.next.DeleteTrip(ctx, ID)
}
