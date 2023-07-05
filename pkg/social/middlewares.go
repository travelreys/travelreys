package social

import (
	"context"
	"errors"

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

func ServiceWithRBACMiddleware(svc Service, tripSvc trips.Service, logger *zap.Logger) Service {
	return &rbacMiddleware{svc, tripSvc, logger}
}

func (mw rbacMiddleware) GetProfile(ctx context.Context, id string) (UserProfile, error) {
	return mw.next.GetProfile(ctx, id)
}

func (mw rbacMiddleware) SendFriendRequest(ctx context.Context, initiatorID, targetID string) error {
	ci, err := reqctx.ClientInfoFromCtx(ctx)
	if err != nil || ci.HasEmptyID() {
		return ErrRBAC
	}
	if initiatorID != ci.UserID {
		return ErrRBAC
	}
	return mw.next.SendFriendRequest(ctx, initiatorID, targetID)
}

func (mw rbacMiddleware) GetFriendRequestByID(ctx context.Context, id string) (FriendRequest, error) {
	ci, err := reqctx.ClientInfoFromCtx(ctx)
	if err != nil || ci.HasEmptyID() {
		return FriendRequest{}, ErrRBAC
	}
	return mw.next.GetFriendRequestByID(ctx, id)
}

func (mw rbacMiddleware) AcceptFriendRequest(ctx context.Context, userid, reqid string) error {
	ci, err := reqctx.ClientInfoFromCtx(ctx)
	if err != nil || ci.HasEmptyID() {
		return ErrRBAC
	}
	req, err := mw.next.GetFriendRequestByID(ctx, reqid)
	if err != nil {
		return err
	}
	if !(req.TargetID == ci.UserID && ci.UserID == userid) {
		return ErrRBAC
	}
	ctxWithFriendRequestInfo := ContextWithFriendRequestInfo(ctx, req)
	return mw.next.AcceptFriendRequest(ctxWithFriendRequestInfo, userid, reqid)
}

func (mw rbacMiddleware) DeleteFriendRequest(ctx context.Context, userid, reqid string) error {
	ci, err := reqctx.ClientInfoFromCtx(ctx)
	if err != nil || ci.HasEmptyID() {
		return ErrRBAC
	}
	req, err := mw.next.GetFriendRequestByID(ctx, reqid)
	if err != nil {
		return err
	}
	if !(req.TargetID == ci.UserID && ci.UserID == userid) {
		return ErrRBAC
	}
	ctxWithFriendRequestInfo := ContextWithFriendRequestInfo(ctx, req)
	return mw.next.DeleteFriendRequest(ctxWithFriendRequestInfo, userid, reqid)
}

func (mw rbacMiddleware) ListFriendRequests(ctx context.Context, ff ListFriendRequestsFilter) (FriendRequestList, error) {
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
	return mw.next.ListFriendRequests(ctx, ff)
}

func (mw rbacMiddleware) ListFollowers(ctx context.Context, userID string) (FriendsList, error) {
	ci, err := reqctx.ClientInfoFromCtx(ctx)
	if err != nil || ci.HasEmptyID() {
		return nil, ErrRBAC
	}

	if ci.UserID != ci.UserID {
		return nil, ErrRBAC
	}
	return mw.next.ListFollowers(ctx, userID)
}

func (mw rbacMiddleware) ListFollowing(ctx context.Context, userID string) (FriendsList, error) {
	ci, err := reqctx.ClientInfoFromCtx(ctx)
	if err != nil || ci.HasEmptyID() {
		return nil, ErrRBAC
	}

	if ci.UserID != ci.UserID {
		return nil, ErrRBAC
	}
	return mw.next.ListFollowing(ctx, userID)
}

func (mw rbacMiddleware) DeleteFriend(ctx context.Context, userID, bindingKey string) error {
	ci, err := reqctx.ClientInfoFromCtx(ctx)
	if err != nil || ci.HasEmptyID() {
		return ErrRBAC
	}
	if ci.UserID != userID {
		return ErrRBAC
	}
	return mw.next.DeleteFriend(ctx, userID, bindingKey)
}

func (mw rbacMiddleware) AreTheyFriends(ctx context.Context, initiatorID, targetID string) (bool, error) {
	ci, err := reqctx.ClientInfoFromCtx(ctx)
	if err != nil || ci.HasEmptyID() {
		return false, ErrRBAC
	}
	if ci.UserID != initiatorID && ci.UserID != targetID {
		return false, ErrRBAC
	}
	return mw.next.AreTheyFriends(ctx, initiatorID, targetID)
}

func (mw rbacMiddleware) ReadTripPublicInfo(ctx context.Context, tripID, referrerID string) (trips.Trip, UserProfile, error) {
	trip, err := mw.tripSvc.Read(ctx, tripID)
	if err != nil {
		return trips.Trip{}, UserProfile{}, err
	}
	ctxWithTripInfo := trips.ContextWithTripInfo(ctx, trip)

	// Allow access if the trip is public
	if trip.IsSharingEnabled() {
		return mw.next.ReadTripPublicInfo(ctxWithTripInfo, tripID, referrerID)
	}

	// Allow access if you are a member of the trip
	ci, err := reqctx.ClientInfoFromCtx(ctx)
	if err != nil || ci.HasEmptyID() {
		return trips.Trip{}, UserProfile{}, ErrTripSharingNotEnabled
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
		if _, err := mw.next.AreTheyFriends(ctx, ci.UserID, referrerID); err == nil {
			return mw.next.ReadTripPublicInfo(ctxWithTripInfo, tripID, referrerID)
		}
	}

	return trips.Trip{}, UserProfile{}, ErrTripSharingNotEnabled
}

func (mw rbacMiddleware) ListTripPublicInfo(ctx context.Context, ff trips.ListFilter) (trips.TripsList, error) {
	ci, err := reqctx.ClientInfoFromCtx(ctx)
	if err != nil || ci.HasEmptyID() {
		return nil, ErrRBAC
	}

	if ci.UserID == *ff.UserID {
		return mw.next.ListTripPublicInfo(ctx, ff)
	}

	ok, err := mw.next.AreTheyFriends(ctx, ci.UserID, *ff.UserID)
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
