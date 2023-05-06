package footprints

import (
	"context"

	"github.com/travelreys/travelreys/pkg/maps"
	"github.com/travelreys/travelreys/pkg/trips"
	"go.uber.org/zap"
)

const (
	svcLoggerName = "footprint.svc"
)

type Service interface {
	CheckIn(ctx context.Context, userID, tripID string, activity trips.Activity) error
	List(ctx context.Context, ff ListFootprintsFilter) (FootprintList, error)
}

type service struct {
	store  Store
	logger *zap.Logger
}

func NewService(store Store, logger *zap.Logger) Service {
	return &service{store, logger.Named(svcLoggerName)}
}

func (svc *service) CheckIn(ctx context.Context, userID, tripID string, activity trips.Activity) error {
	placeID := activity.Place.Labels[maps.LabelPlaceID]

	var (
		fp  Footprint
		err error
	)
	fp, err = svc.store.Read(ctx, userID, placeID)
	if err != nil && err != ErrFpNotFound {
		return err
	}
	if err == ErrFpNotFound {
		fp = NewFootprint(userID, tripID, activity)
	} else {
		fp.AddNewCheckin(tripID, activity)
	}

	return svc.store.Save(ctx, userID, fp)
}

func (svc *service) List(ctx context.Context, ff ListFootprintsFilter) (FootprintList, error) {
	return svc.store.List(ctx, ff)
}
