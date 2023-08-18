package trips

import (
	"context"
	"errors"

	"github.com/travelreys/travelreys/pkg/common"
)

var (
	ErrNoTripInfoSet   = errors.New("trips.ErrNoTripInfoSet")
	ErrNoInviteInfoSet = errors.New("trips.ErrNoInviteInfoSet")
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

type InviteInfo struct {
	Invite Invite
}

func ContextWithInviteInfo(ctx context.Context, invite Invite) context.Context {
	return context.WithValue(
		ctx,
		common.ContextKeyInviteInfo,
		InviteInfo{invite},
	)
}

func InviteInfoFromCtx(ctx context.Context) (InviteInfo, error) {
	val := ctx.Value(common.ContextKeyInviteInfo)
	if val == nil {
		return InviteInfo{}, ErrNoTripInfoSet
	}

	ti, _ := val.(InviteInfo)
	return ti, nil
}
