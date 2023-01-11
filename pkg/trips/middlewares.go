package trips

import (
	"fmt"
	"time"

	"github.com/tiinyplanet/tiinyplanet/pkg/reqctx"
	"go.uber.org/zap"
)

type loggingMiddleware struct {
	next   Service
	logger *zap.Logger
}

func ServiceWithLoggingMiddleware(svc Service, logger *zap.Logger) Service {
	return &loggingMiddleware{svc, logger}
}

func (mw *loggingMiddleware) CreateTripPlan(ctx reqctx.Context, creator TripMember, name string, start, end time.Time) (TripPlan, error) {
	result, err := mw.next.CreateTripPlan(ctx, creator, name, start, end)
	if err != nil {
		mw.logger.Error(
			"CreateTripPlan",
			zap.String("creator", fmt.Sprintf("%+v", creator)),
			zap.String("name", name),
			zap.Time("start", start),
			zap.Time("end", end),
			zap.Error(err),
		)
	}
	return result, err
}

func (mw *loggingMiddleware) ReadTripPlan(ctx reqctx.Context, ID string) (TripPlan, error) {
	result, err := mw.next.ReadTripPlan(ctx, ID)
	if err != nil {
		mw.logger.Error("ReadTripPlan", zap.String("id", ID), zap.Error(err))
	}
	return result, err
}

func (mw *loggingMiddleware) ListTripPlans(ctx reqctx.Context, ff ListTripPlansFilter) ([]TripPlan, error) {
	result, err := mw.next.ListTripPlans(ctx, ff)
	if err != nil {
		mw.logger.Error("ListTripPlans", zap.String("ff", fmt.Sprintf("%+v", ff)), zap.Error(err))
	}
	return result, err
}

func (mw *loggingMiddleware) DeleteTripPlan(ctx reqctx.Context, ID string) error {
	err := mw.next.DeleteTripPlan(ctx, ID)
	if err != nil {
		mw.logger.Error("ListTripPlans", zap.Error(err))
	}
	return err
}
