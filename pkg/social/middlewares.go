package social

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

func (mw rbacMiddleware) AcceptFriendRequest(ctx context.Context, id string) error {
	ci, err := reqctx.ClientInfoFromCtx(ctx)
	if err != nil || ci.HasEmptyID() {
		return ErrRBAC
	}
	req, err := mw.next.GetFriendRequestByID(ctx, id)
	if err != nil {
		return err
	}
	if req.TargetID != ci.UserID {
		return ErrRBAC
	}
	ctxWithFriendRequestInfo := ContextWithFriendRequestInfo(ctx, req)
	return mw.next.AcceptFriendRequest(ctxWithFriendRequestInfo, id)
}

func (mw rbacMiddleware) GetFriendRequestByID(ctx context.Context, id string) (FriendRequest, error) {
	ci, err := reqctx.ClientInfoFromCtx(ctx)
	if err != nil || ci.HasEmptyID() {
		return FriendRequest{}, ErrRBAC
	}
	return mw.next.GetFriendRequestByID(ctx, id)
}

func (mw rbacMiddleware) DeleteFriendRequest(ctx context.Context, id string) error {
	ci, err := reqctx.ClientInfoFromCtx(ctx)
	if err != nil || ci.HasEmptyID() {
		return ErrRBAC
	}
	req, err := mw.next.GetFriendRequestByID(ctx, id)
	if err != nil {
		return err
	}
	if req.TargetID != ci.UserID {
		return ErrRBAC
	}
	ctxWithFriendRequestInfo := ContextWithFriendRequestInfo(ctx, req)
	return mw.next.DeleteFriendRequest(ctxWithFriendRequestInfo, id)
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

func (mw rbacMiddleware) ListFriends(ctx context.Context, userID string) (FriendsList, error) {
	ci, err := reqctx.ClientInfoFromCtx(ctx)
	if err != nil || ci.HasEmptyID() {
		return nil, ErrRBAC
	}

	if ci.UserID != ci.UserID {
		return nil, ErrRBAC
	}
	return mw.next.ListFriends(ctx, userID)
}

func (mw rbacMiddleware) DeleteFriend(ctx context.Context, userID, targetID string) error {
	ci, err := reqctx.ClientInfoFromCtx(ctx)
	if err != nil || ci.HasEmptyID() {
		return ErrRBAC
	}
	if ci.UserID != ci.UserID {
		return ErrRBAC
	}
	return mw.next.DeleteFriend(ctx, userID, targetID)
}

func (mw rbacMiddleware) AreTheyFriends(ctx context.Context, userOneID, userTwoID string) error {
	return ErrRBAC
}
