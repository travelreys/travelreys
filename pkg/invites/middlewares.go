package invites

import (
	"context"
	"errors"
	"net/http"

	"github.com/travelreys/travelreys/pkg/auth"
	"github.com/travelreys/travelreys/pkg/common"
	"github.com/travelreys/travelreys/pkg/reqctx"
	"github.com/travelreys/travelreys/pkg/trips"
	"go.uber.org/zap"
)

var (
	ErrRBAC = errors.New("invites.ErrRBAC")
)

type validationMiddleware struct {
	next   Service
	logger *zap.Logger
}

func SvcWithValidationMw(
	svc Service,
	logger *zap.Logger,
) Service {
	return &validationMiddleware{
		svc,
		logger.Named("invites.validationMiddleware"),
	}
}

// App Invites

func (mw *validationMiddleware) SendAppInvite(
	ctx context.Context,
	authorID,
	userEmail string,
) error {
	if authorID == "" || userEmail == "" {
		mw.logger.Warn("SendAppInvite")
		return common.ErrValidation
	}
	return mw.next.SendAppInvite(ctx, authorID, userEmail)
}

func (mw *validationMiddleware) AcceptAppInvite(
	ctx context.Context,
	ID,
	code,
	sig string,
) (auth.User, *http.Cookie, error) {
	if ID == "" || code == "" || sig == "" {
		mw.logger.Warn("AcceptAppInvite")
		return auth.User{}, nil, common.ErrValidation
	}
	return mw.next.AcceptAppInvite(ctx, ID, code, sig)
}

// Trip Invites

func (mw *validationMiddleware) SendTripInvite(
	ctx context.Context,
	tripID,
	authorID,
	userID string,
) error {
	if tripID == "" || authorID == "" || userID == "" {
		mw.logger.Warn("SendTripInvite")
		return common.ErrValidation
	}
	return mw.next.SendTripInvite(ctx, tripID, authorID, userID)
}

func (mw *validationMiddleware) AcceptTripInvite(ctx context.Context, ID string) error {
	if ID == "" {
		mw.logger.Warn("AcceptTripInvite")
		return common.ErrValidation
	}
	return mw.next.AcceptTripInvite(ctx, ID)
}

func (mw *validationMiddleware) ReadTripInvite(
	ctx context.Context,
	ID string,
) (TripInvite, error) {
	if ID == "" {
		mw.logger.Warn("ReadTripInvite")
		return TripInvite{}, common.ErrValidation
	}
	return mw.next.ReadTripInvite(ctx, ID)
}

func (mw *validationMiddleware) DeclineTripInvite(ctx context.Context, ID string) error {
	if ID == "" {
		mw.logger.Warn("DeclineTripInvite")
		return common.ErrValidation
	}
	return mw.next.DeclineTripInvite(ctx, ID)
}

func (mw *validationMiddleware) ListTripInvites(
	ctx context.Context,
	ff ListTripInvitesFilter,
) (TripInviteList, error) {
	if err := ff.Validate(); err != nil {
		mw.logger.Warn("ListTripInvites")
		return nil, common.ErrValidation
	}
	return mw.next.ListTripInvites(ctx, ff)
}

// Email Trip Invites

func (mw *validationMiddleware) SendEmailTripInvite(
	ctx context.Context,
	tripID,
	authorID,
	userEmail string,
) error {
	if tripID == "" || authorID == "" || userEmail == "" {
		mw.logger.Warn("SendEmailTripInvite")
		return common.ErrValidation
	}
	return mw.next.SendEmailTripInvite(ctx, tripID, authorID, userEmail)
}

func (mw *validationMiddleware) AcceptEmailTripInvite(
	ctx context.Context,
	ID,
	code,
	sig string,
) (auth.User, *http.Cookie, error) {
	if ID == "" || sig == "" || code == "" {
		mw.logger.Warn("AcceptEmailTripInvite")
		return auth.User{}, nil, common.ErrValidation
	}
	return mw.next.AcceptEmailTripInvite(ctx, ID, code, sig)
}

func (mw *validationMiddleware) ReadEmailTripInvite(
	ctx context.Context,
	ID string,
) (EmailTripInvite, error) {
	if ID == "" {
		mw.logger.Warn("ReadEmailTripInvite")
		return EmailTripInvite{}, common.ErrValidation
	}
	return mw.next.ReadEmailTripInvite(ctx, ID)
}

func (mw *validationMiddleware) ListEmailTripInvites(
	ctx context.Context,
	ff ListEmailTripInvitesFilter,
) (EmailTripInviteList, error) {
	if err := ff.Validate(); err != nil {
		mw.logger.Warn("ListEmailTripInvites")
		return EmailTripInviteList{}, common.ErrValidation
	}
	return mw.next.ListEmailTripInvites(ctx, ff)
}

type rbacMiddleware struct {
	next    Service
	tripSvc trips.Service
	authSvc auth.Service
	logger  *zap.Logger
}

func SvcWithRBACMw(
	svc Service,
	tripSvc trips.Service,
	authSvc auth.Service,
	logger *zap.Logger,
) Service {
	return &rbacMiddleware{
		svc,
		tripSvc,
		authSvc,
		logger.Named("invites.rbacMiddleware"),
	}
}

// App Invites

func (mw *rbacMiddleware) SendAppInvite(
	ctx context.Context,
	authorID,
	userEmail string,
) error {
	ci, err := reqctx.ClientInfoFromCtx(ctx)
	if err != nil || ci.HasEmptyID() {
		return ErrRBAC
	}
	if authorID != ci.UserID {
		return ErrRBAC
	}
	return mw.next.SendAppInvite(ctx, authorID, userEmail)
}

func (mw *rbacMiddleware) AcceptAppInvite(
	ctx context.Context,
	ID,
	code,
	sig string,
) (auth.User, *http.Cookie, error) {
	return mw.next.AcceptAppInvite(ctx, ID, code, sig)
}

// Trip Invites

func (mw *rbacMiddleware) SendTripInvite(
	ctx context.Context,
	tripID,
	authorID,
	userID string,
) error {
	ci, err := reqctx.ClientInfoFromCtx(ctx)
	if err != nil || ci.HasEmptyID() {
		return ErrRBAC
	}

	users, err := mw.authSvc.List(ctx, auth.ListFilter{
		IDs: []string{userID, authorID},
	})
	if err != nil {
		return err
	}
	if len(users) != 2 {
		return ErrInvalidInvite
	}
	t, err := mw.tripSvc.Read(ctx, tripID)
	if err != nil {
		return err
	}
	if !common.StringContains(t.GetMemberIDs(), authorID) {
		return ErrRBAC
	}

	author := users[0]
	user := users[1]
	if users[1].ID == authorID {
		author = users[1]
		user = users[0]
	}

	return mw.next.SendTripInvite(
		ContextWithTripInviteMetaInfo(ctx, &user, &author, t),
		tripID,
		authorID,
		userID,
	)
}

func (mw *rbacMiddleware) ReadTripInvite(
	ctx context.Context,
	ID string,
) (TripInvite, error) {
	return TripInvite{}, ErrRBAC
}

func (mw *rbacMiddleware) AcceptTripInvite(ctx context.Context, ID string) error {
	ci, err := reqctx.ClientInfoFromCtx(ctx)
	if err != nil || ci.HasEmptyID() {
		return ErrRBAC
	}
	invite, err := mw.next.ReadTripInvite(ctx, ID)
	if err != nil {
		return err
	}

	if invite.UserID != ci.UserID {
		return ErrRBAC
	}

	return mw.next.AcceptTripInvite(ContextWithTripInviteInfo(ctx, invite), ID)
}

func (mw *rbacMiddleware) DeclineTripInvite(ctx context.Context, ID string) error {
	ci, err := reqctx.ClientInfoFromCtx(ctx)
	if err != nil || ci.HasEmptyID() {
		return ErrRBAC
	}
	invite, err := mw.next.ReadTripInvite(ctx, ID)
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

	return mw.next.DeclineTripInvite(ContextWithTripInviteInfo(ctx, invite), ID)
}

func (mw *rbacMiddleware) ListTripInvites(
	ctx context.Context,
	ff ListTripInvitesFilter,
) (TripInviteList, error) {
	ci, err := reqctx.ClientInfoFromCtx(ctx)
	if err != nil || ci.HasEmptyID() {
		return nil, ErrRBAC
	}
	return mw.next.ListTripInvites(ctx, ff)
}

// Email Trip Invites

func (mw *rbacMiddleware) SendEmailTripInvite(
	ctx context.Context,
	tripID,
	authorID,
	userEmail string,
) error {
	ci, err := reqctx.ClientInfoFromCtx(ctx)
	if err != nil || ci.HasEmptyID() {
		return ErrRBAC
	}

	author, err := mw.authSvc.Read(ctx, authorID)
	if err != nil {
		return err
	}
	t, err := mw.tripSvc.Read(ctx, tripID)
	if err != nil {
		return err
	}
	if !common.StringContains(t.GetMemberIDs(), authorID) {
		return ErrRBAC
	}

	return mw.next.SendEmailTripInvite(
		ContextWithTripInviteMetaInfo(ctx, nil, &author, t),
		tripID,
		authorID,
		userEmail,
	)
}

func (mw *rbacMiddleware) AcceptEmailTripInvite(
	ctx context.Context,
	ID,
	code,
	sig string,
) (auth.User, *http.Cookie, error) {
	return mw.next.AcceptEmailTripInvite(ctx, ID, code, sig)
}

func (mw *rbacMiddleware) ReadEmailTripInvite(
	ctx context.Context,
	ID string,
) (EmailTripInvite, error) {
	return EmailTripInvite{}, ErrRBAC
}

func (mw *rbacMiddleware) ListEmailTripInvites(
	ctx context.Context,
	ff ListEmailTripInvitesFilter,
) (EmailTripInviteList, error) {
	return EmailTripInviteList{}, ErrRBAC
}
