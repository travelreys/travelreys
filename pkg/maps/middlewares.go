package maps

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

func (mw rbacMiddleware) PlacesAutocomplete(ctx context.Context, query, types, sessiontoken string) (AutocompletePredictionList, error) {
	ci, err := common.ReadClientInfoFromCtx(ctx)
	if err != nil || ci.HasEmptyID() {
		return AutocompletePredictionList{}, ErrRBACMissing
	}
	return mw.next.PlacesAutocomplete(ctx, query, types, sessiontoken)
}

func (mw rbacMiddleware) PlaceDetails(ctx context.Context, placeID string, fields []string, sessiontoken string) (Place, error) {
	ci, err := common.ReadClientInfoFromCtx(ctx)
	if err != nil || ci.HasEmptyID() {
		return Place{}, ErrRBACMissing
	}
	return mw.next.PlaceDetails(ctx, placeID, fields, sessiontoken)
}

func (mw rbacMiddleware) Directions(ctx context.Context, originPlaceID, destPlaceID, mode string) (RouteList, error) {
	ci, err := common.ReadClientInfoFromCtx(ctx)
	if err != nil || ci.HasEmptyID() {
		return RouteList{}, ErrRBACMissing
	}
	return mw.next.Directions(ctx, originPlaceID, destPlaceID, mode)
}

func (mw rbacMiddleware) OptimizeRoute(ctx context.Context, originPlaceID, destPlaceID string, waypointsPlaceID []string) (RouteList, error) {
	ci, err := common.ReadClientInfoFromCtx(ctx)
	if err != nil || ci.HasEmptyID() {
		return RouteList{}, ErrRBACMissing
	}
	return mw.next.OptimizeRoute(ctx, originPlaceID, destPlaceID, waypointsPlaceID)
}
