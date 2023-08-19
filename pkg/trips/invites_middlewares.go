package trips

import (
	"context"

	"github.com/travelreys/travelreys/pkg/auth"
	"github.com/travelreys/travelreys/pkg/common"
	"github.com/travelreys/travelreys/pkg/reqctx"
	"go.uber.org/zap"
)

type inviteValidationMiddleware struct {
	next   InviteService
	logger *zap.Logger
}

func InviteSvcWithValidationMw(
	svc InviteService,
	logger *zap.Logger,
) InviteService {
	return &inviteValidationMiddleware{
		svc,
		logger.Named("trips.inviteValidationMiddleware"),
	}
}

func (mw *inviteValidationMiddleware) Send(
	ctx context.Context,
	tripID,
	authorID,
	userID string,
) error {
	if tripID == "" || authorID == "" || userID == "" {
		return ErrValidation
	}
	return mw.next.Send(ctx, tripID, authorID, userID)
}

func (mw *inviteValidationMiddleware) Accept(ctx context.Context, ID string) error {
	if ID == "" {
		return ErrValidation
	}
	return mw.next.Accept(ctx, ID)
}

func (mw *inviteValidationMiddleware) Read(ctx context.Context, ID string) (Invite, error) {
	if ID == "" {
		return Invite{}, ErrValidation
	}
	return mw.next.Read(ctx, ID)
}

func (mw *inviteValidationMiddleware) Decline(ctx context.Context, ID string) error {
	if ID == "" {
		return ErrValidation
	}
	return mw.next.Decline(ctx, ID)
}

func (mw *inviteValidationMiddleware) List(ctx context.Context, ff ListInvitesFilter) (InviteList, error) {
	if err := ff.Validate(); err != nil {
		return nil, ErrValidation
	}
	return mw.next.List(ctx, ff)
}

type inviteRBACMiddleware struct {
	next    InviteService
	tripSvc Service
	authSvc auth.Service
	logger  *zap.Logger
}

func InviteSvcWithRBACMw(
	svc InviteService,
	tripSvc Service,
	authSvc auth.Service,
	logger *zap.Logger,
) InviteService {
	return &inviteRBACMiddleware{
		svc,
		tripSvc,
		authSvc,
		logger.Named("trips.inviteRBACMiddleware"),
	}
}

func (mw *inviteRBACMiddleware) Send(
	ctx context.Context,
	tripID,
	authorID,
	userID string,
) error {
	ci, err := reqctx.ClientInfoFromCtx(ctx)
	if err != nil || ci.HasEmptyID() {
		return ErrRBAC
	}

	if _, err := mw.authSvc.Read(ctx, userID); err != nil {
		return err
	}
	t, err := mw.tripSvc.Read(ctx, tripID)
	if err != nil {
		return err
	}

	if !common.StringContains(t.GetMemberIDs(), authorID) {
		return ErrRBAC
	}

	return mw.next.Send(ctx, tripID, authorID, userID)
}

func (mw *inviteRBACMiddleware) Read(ctx context.Context, ID string) (Invite, error) {
	return Invite{}, ErrRBAC
}

func (mw *inviteRBACMiddleware) Accept(ctx context.Context, ID string) error {
	ci, err := reqctx.ClientInfoFromCtx(ctx)
	if err != nil || ci.HasEmptyID() {
		return ErrRBAC
	}
	invite, err := mw.next.Read(ctx, ID)
	if err != nil {
		return err
	}

	if invite.UserID != ID {
		return ErrRBAC
	}

	return mw.next.Accept(ContextWithInviteInfo(ctx, invite), ID)
}

func (mw *inviteRBACMiddleware) Decline(ctx context.Context, ID string) error {
	ci, err := reqctx.ClientInfoFromCtx(ctx)
	if err != nil || ci.HasEmptyID() {
		return ErrRBAC
	}
	invite, err := mw.next.Read(ctx, ID)
	if err != nil {
		return err
	}

	t, err := mw.tripSvc.Read(ctx, invite.TripID)
	if err != nil {
		return err
	}

	if invite.UserID != ci.UserID && !common.StringContains(t.GetMemberIDs(), ci.UserID) {
		return ErrRBAC
	}

	return mw.next.Decline(ContextWithInviteInfo(ctx, invite), ID)
}

func (mw *inviteRBACMiddleware) List(ctx context.Context, ff ListInvitesFilter) (InviteList, error) {
	ci, err := reqctx.ClientInfoFromCtx(ctx)
	if err != nil || ci.HasEmptyID() {
		return nil, ErrRBAC
	}
	return mw.next.List(ctx, ff)
}
