package maps

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

func (mw rbacMiddleware) PlacesAutocomplete(ctx context.Context, query, types, sessiontoken, lang string) (AutocompletePredictionList, error) {
	ci, err := reqctx.ClientInfoFromCtx(ctx)
	if err != nil || ci.HasEmptyID() {
		return AutocompletePredictionList{}, ErrRBAC
	}
	return mw.next.PlacesAutocomplete(ctx, query, types, sessiontoken, lang)
}

func (mw rbacMiddleware) PlaceDetails(ctx context.Context, placeID string, fields []string, sessiontoken, lang string) (Place, error) {
	ci, err := reqctx.ClientInfoFromCtx(ctx)
	if err != nil || ci.HasEmptyID() {
		return Place{}, ErrRBAC
	}
	return mw.next.PlaceDetails(ctx, placeID, fields, sessiontoken, lang)
}

func (mw rbacMiddleware) PlaceAtmosphere(ctx context.Context, placeID string, fields []string, sessiontoken, lang string) (PlaceAtmosphere, error) {
	ci, err := reqctx.ClientInfoFromCtx(ctx)
	if err != nil || ci.HasEmptyID() {
		return PlaceAtmosphere{}, ErrRBAC
	}
	return mw.next.PlaceAtmosphere(ctx, placeID, fields, sessiontoken, lang)
}

func (mw rbacMiddleware) Directions(ctx context.Context, originPlaceID, destPlaceID string, modes []string) (RouteList, error) {
	ci, err := reqctx.ClientInfoFromCtx(ctx)
	if err != nil || ci.HasEmptyID() {
		return RouteList{}, ErrRBAC
	}
	return mw.next.Directions(ctx, originPlaceID, destPlaceID, modes)
}

func (mw rbacMiddleware) OptimizeRoute(ctx context.Context, originPlaceID, destPlaceID string, waypointsPlaceID []string) (RouteList, []int, error) {
	ci, err := reqctx.ClientInfoFromCtx(ctx)
	if err != nil || ci.HasEmptyID() {
		return RouteList{}, nil, ErrRBAC
	}
	return mw.next.OptimizeRoute(ctx, originPlaceID, destPlaceID, waypointsPlaceID)
}
