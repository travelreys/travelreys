package invites

import (
	"context"
	"errors"

	"github.com/travelreys/travelreys/pkg/auth"
	"github.com/travelreys/travelreys/pkg/common"
	"github.com/travelreys/travelreys/pkg/trips"
)

var (
	ErrTripInviteMetaInfoNotSet = errors.New("invites.ErrTripInviteMetaInfoNotSet")
	ErrTripInviteInfoNotSet     = errors.New("invites.ErrTripInviteInfoNotSet")
)

type TripInviteMetaInfo struct {
	User   *auth.User
	Author *auth.User
	Trip   *trips.Trip
}

func ContextWithTripInviteMetaInfo(
	ctx context.Context,
	user *auth.User,
	author *auth.User,
	trip *trips.Trip,
) context.Context {
	return context.WithValue(
		ctx,
		common.ContextKeyTripInviteMetaInfo,
		TripInviteMetaInfo{user, author, trip},
	)
}

func TripInviteMetaInfoFromCtx(ctx context.Context) (TripInviteMetaInfo, error) {
	val := ctx.Value(common.ContextKeyTripInviteMetaInfo)
	if val == nil {
		return TripInviteMetaInfo{}, ErrTripInviteMetaInfoNotSet
	}

	ti, _ := val.(TripInviteMetaInfo)
	return ti, nil
}

type TripInviteInfo struct {
	Invite TripInvite
}

func ContextWithTripInviteInfo(
	ctx context.Context,
	invite TripInvite,
) context.Context {
	return context.WithValue(
		ctx,
		common.ContextKeyTripInviteInfo,
		TripInviteInfo{invite},
	)
}

func TripInviteInfoFromCtx(ctx context.Context) (TripInviteInfo, error) {
	val := ctx.Value(common.ContextKeyTripInviteInfo)
	if val == nil {
		return TripInviteInfo{}, ErrTripInviteInfoNotSet
	}

	ti, _ := val.(TripInviteInfo)
	return ti, nil
}
