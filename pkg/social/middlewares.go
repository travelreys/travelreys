package social

import (
	"context"
	"errors"
	"time"

	"github.com/travelreys/travelreys/pkg/common"
	"github.com/travelreys/travelreys/pkg/reqctx"
	"github.com/travelreys/travelreys/pkg/trips"
	"go.uber.org/zap"
)

var (
	ErrRBAC                  = errors.New("trip.ErrRBAC")
	ErrTripSharingNotEnabled = errors.New("trip.ErrTripSharingNotEnabled")
)

type rbacMiddleware struct {
	next    Service
	tripSvc trips.Service
	logger  *zap.Logger
}

func SvcWithRBACMw(svc Service, tripSvc trips.Service, logger *zap.Logger) Service {
	return &rbacMiddleware{svc, tripSvc, logger}
}

func (mw rbacMiddleware) GetProfile(ctx context.Context, id string) (UserProfile, error) {
	return mw.next.GetProfile(ctx, id)
}

func (mw rbacMiddleware) SendFollowRequest(ctx context.Context, initiatorID, targetID string) error {
	ci, err := reqctx.ClientInfoFromCtx(ctx)
	if err != nil || ci.HasEmptyID() {
		return ErrRBAC
	}
	if initiatorID != ci.UserID {
		return ErrRBAC
	}
	return mw.next.SendFollowRequest(ctx, initiatorID, targetID)
}

func (mw rbacMiddleware) GetFollowRequestByID(ctx context.Context, id string) (FollowRequest, error) {
	ci, err := reqctx.ClientInfoFromCtx(ctx)
	if err != nil || ci.HasEmptyID() {
		return FollowRequest{}, ErrRBAC
	}
	return mw.next.GetFollowRequestByID(ctx, id)
}

func (mw rbacMiddleware) AcceptFollowRequest(
	ctx context.Context,
	userID,
	initiatorID,
	reqID string,
) error {
	ci, err := reqctx.ClientInfoFromCtx(ctx)
	if err != nil || ci.HasEmptyID() {
		return ErrRBAC
	}
	req, err := mw.next.GetFollowRequestByID(ctx, reqID)
	if err != nil {
		return err
	}
	if !(req.TargetID == ci.UserID && ci.UserID == userID) {
		return ErrRBAC
	}
	return mw.next.AcceptFollowRequest(
		ContextWithFollowRequestInfo(ctx, req),
		userID,
		initiatorID,
		reqID,
	)
}

func (mw rbacMiddleware) DeleteFollowRequest(ctx context.Context, userID, reqID string) error {
	ci, err := reqctx.ClientInfoFromCtx(ctx)
	if err != nil || ci.HasEmptyID() {
		return ErrRBAC
	}
	req, err := mw.next.GetFollowRequestByID(ctx, reqID)
	if err != nil {
		return err
	}
	if !(req.TargetID == ci.UserID && ci.UserID == userID) {
		return ErrRBAC
	}
	ctxWithFollowRequestInfo := ContextWithFollowRequestInfo(ctx, req)
	return mw.next.DeleteFollowRequest(ctxWithFollowRequestInfo, userID, reqID)
}

func (mw rbacMiddleware) ListFollowRequests(ctx context.Context, ff ListFollowRequestsFilter) (FollowRequestList, error) {
	ci, err := reqctx.ClientInfoFromCtx(ctx)
	if err != nil || ci.HasEmptyID() {
		return nil, ErrRBAC
	}
	if ff.InitiatorID != nil && *ff.InitiatorID != ci.UserID {
		return nil, ErrRBAC
	}
	if ff.TargetID != nil && *ff.TargetID != ci.UserID {
		return nil, ErrRBAC
	}
	return mw.next.ListFollowRequests(ctx, ff)
}

func (mw rbacMiddleware) ListFollowers(ctx context.Context, userID string) (FollowingsList, error) {
	ci, err := reqctx.ClientInfoFromCtx(ctx)
	if err != nil || ci.HasEmptyID() {
		return nil, ErrRBAC
	}

	if ci.UserID == userID {
		return mw.next.ListFollowers(ctx, userID)
	}

	ok, err := mw.next.IsFollowing(ctx, ci.UserID, userID)
	if err != nil {
		return FollowingsList{}, err
	}
	if !ok {
		return FollowingsList{}, ErrRBAC
	}
	return mw.next.ListFollowers(ctx, userID)

}

func (mw rbacMiddleware) ListFollowing(ctx context.Context, userID string) (FollowingsList, error) {
	ci, err := reqctx.ClientInfoFromCtx(ctx)
	if err != nil || ci.HasEmptyID() {
		return nil, ErrRBAC
	}

	if ci.UserID == userID {
		return mw.next.ListFollowing(ctx, userID)
	}

	ok, err := mw.next.IsFollowing(ctx, ci.UserID, userID)
	if err != nil {
		return FollowingsList{}, err
	}
	if !ok {
		return FollowingsList{}, ErrRBAC
	}

	return mw.next.ListFollowing(ctx, userID)
}

func (mw rbacMiddleware) DeleteFollowing(ctx context.Context, userID, bindingKey string) error {
	ci, err := reqctx.ClientInfoFromCtx(ctx)
	if err != nil || ci.HasEmptyID() {
		return ErrRBAC
	}
	if ci.UserID != userID {
		return ErrRBAC
	}
	return mw.next.DeleteFollowing(ctx, userID, bindingKey)
}

func (mw rbacMiddleware) IsFollowing(ctx context.Context, initiatorID, targetID string) (bool, error) {
	ci, err := reqctx.ClientInfoFromCtx(ctx)
	if err != nil || ci.HasEmptyID() {
		return false, ErrRBAC
	}
	if ci.UserID != initiatorID && ci.UserID != targetID {
		return false, ErrRBAC
	}
	return mw.next.IsFollowing(ctx, initiatorID, targetID)
}

func (mw rbacMiddleware) ReadTripPublicInfo(ctx context.Context, tripID, referrerID string) (*trips.Trip, UserProfile, error) {
	trip, err := mw.tripSvc.Read(ctx, tripID)
	if err != nil {
		return nil, UserProfile{}, err
	}
	ctxWithTripInfo := trips.ContextWithTripInfo(ctx, trip)

	// Allow access if the trip is public
	if trip.IsSharingEnabled() {
		return mw.next.ReadTripPublicInfo(ctxWithTripInfo, tripID, referrerID)
	}

	// Allow access if you are a member of the trip
	ci, err := reqctx.ClientInfoFromCtx(ctx)
	if err != nil || ci.HasEmptyID() {
		return nil, UserProfile{}, ErrRBAC
	}
	membersID := []string{trip.Creator.ID}
	for _, mem := range trip.Members {
		membersID = append(membersID, mem.ID)
	}
	if common.StringContains(membersID, ci.UserID) {
		return mw.next.ReadTripPublicInfo(ctxWithTripInfo, tripID, referrerID)
	}

	// ReferrerID should be a member ID.
	// Allow access if you are friend of the member of the trip

	if common.StringContains(membersID, referrerID) {
		if _, err := mw.next.IsFollowing(ctx, ci.UserID, referrerID); err == nil {
			return mw.next.ReadTripPublicInfo(ctxWithTripInfo, tripID, referrerID)
		}
	}

	return nil, UserProfile{}, ErrTripSharingNotEnabled
}

func (mw rbacMiddleware) ListTripPublicInfo(ctx context.Context, ff trips.ListFilter) (trips.TripsList, error) {
	ci, err := reqctx.ClientInfoFromCtx(ctx)
	if err != nil || ci.HasEmptyID() {
		return nil, ErrRBAC
	}

	if ci.UserID == *ff.UserID {
		return mw.next.ListTripPublicInfo(ctx, ff)
	}

	ok, err := mw.next.IsFollowing(ctx, ci.UserID, *ff.UserID)
	if err != nil {
		return nil, err
	}
	ff.OnlyPublic = !ok
	return mw.next.ListTripPublicInfo(ctx, ff)
}

func (mw rbacMiddleware) ListFollowingTrips(ctx context.Context, initiatorID string) (trips.TripsList, UserProfileMap, error) {
	ci, err := reqctx.ClientInfoFromCtx(ctx)
	if err != nil || ci.HasEmptyID() {
		return nil, nil, ErrRBAC
	}

	if ci.UserID != initiatorID {
		return nil, nil, ErrRBAC
	}
	return mw.next.ListFollowingTrips(ctx, initiatorID)
}

func (mw rbacMiddleware) DuplicateTrip(
	ctx context.Context,
	initiatorID,
	referrerID,
	tripID,
	name string,
	startDate time.Time,
) (string, error) {
	ci, err := reqctx.ClientInfoFromCtx(ctx)
	if err != nil || ci.HasEmptyID() {
		return "", ErrRBAC
	}

	trip, err := mw.tripSvc.Read(ctx, tripID)
	if err != nil {
		return "", err
	}
	ctxWithTripInfo := trips.ContextWithTripInfo(ctx, trip)

	// Allow access if the trip is public
	if trip.IsSharingEnabled() {
		return mw.next.DuplicateTrip(ctxWithTripInfo, ci.UserID, referrerID, tripID, name, startDate)
	}

	// Allow access if you are a member of the trip
	membersID := []string{trip.Creator.ID}
	for _, mem := range trip.Members {
		membersID = append(membersID, mem.ID)
	}
	if common.StringContains(membersID, ci.UserID) {
		return mw.next.DuplicateTrip(ctxWithTripInfo, ci.UserID, referrerID, tripID, name, startDate)
	}

	// ReferrerID should be a member ID.
	// Allow access if you are friend of the member of the trip
	if common.StringContains(membersID, referrerID) {
		if _, err := mw.next.IsFollowing(ctx, ci.UserID, referrerID); err == nil {
			return mw.next.DuplicateTrip(ctxWithTripInfo, ci.UserID, referrerID, tripID, name, startDate)
		}
	}

	return "", ErrTripSharingNotEnabled
}
