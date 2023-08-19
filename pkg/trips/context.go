package trips

import (
	"context"
	"errors"

	"github.com/travelreys/travelreys/pkg/common"
)

var (
	ErrNoTripInfoSet       = errors.New("invites.ErrNoTripInfoSet")
	ErrNoTripInviteInfoSet = errors.New("invites.ErrNoTripInviteInfoSet")
)

type TripInfo struct {
	Trip *Trip
}

func ContextWithTripInfo(ctx context.Context, trip *Trip) context.Context {
	return context.WithValue(
		ctx,
		common.ContextKeyTripInfo,
		TripInfo{Trip: trip},
	)
}

func TripInfoFromCtx(ctx context.Context) (TripInfo, error) {
	val := ctx.Value(common.ContextKeyTripInfo)
	if val == nil {
		return TripInfo{}, ErrNoTripInfoSet
	}

	ti, _ := val.(TripInfo)
	return ti, nil
}
