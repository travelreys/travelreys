package invites

import (
	"bytes"
	"context"
	"html/template"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/travelreys/travelreys/pkg/auth"
	"github.com/travelreys/travelreys/pkg/email"
	"github.com/travelreys/travelreys/pkg/trips"
	"go.uber.org/zap"
)

const (
	syncMsgWaitInterval = 500 * time.Millisecond
	defaultLoginSender  = "login@travelreys.com"

	tripInviteTmplFilePath      = "assets/tripInviteEmail.tmpl.html"
	tripInviteTmplFileName      = "tripInviteEmail.tmpl.html"
	emailTripInviteTmplFilePath = "assets/emailTripInviteEmail.tmpl.html"
	emailTripInviteTmplFileName = "emailTripInviteEmail.tmpl.html"

	defaultCoverImgURL  = "https://images.unsplash.com/photo-1476514525535-07fb3b4ae5f1?crop=entropy&cs=tinysrgb&fit=max&fm=jpg&ixid=M3wzOTc1ODV8MHwxfHNlYXJjaHwzfHx0cmF2ZWx8ZW58MHwwfHx8MTY4ODc0MDUxNnww&ixlib=rb-4.0.3&q=80&w=1080"
	emailInviteDuration = 168 * time.Hour // 7 days
)

type Service interface {
	SendTripInvite(ctx context.Context, tripID, authorID, userID string) error
	AcceptTripInvite(ctx context.Context, ID string) error
	DeclineTripInvite(ctx context.Context, ID string) error
	ReadTripInvite(ctx context.Context, ID string) (TripInvite, error)
	ListTripInvites(ctx context.Context, ff ListTripInvitesFilter) (TripInviteList, error)

	SendEmailTripInvite(ctx context.Context, tripID, authorID, userEmail string) error
	AcceptEmailTripInvite(ctx context.Context, ID, sig, code string, isLoggedIn bool) (*http.Cookie, error)
	ReadEmailTripInvite(ctx context.Context, ID string) (EmailTripInvite, error)
	ListEmailTripInvites(ctx context.Context, ff ListEmailTripInvitesFilter) (EmailTripInviteList, error)
}

type service struct {
	authSvc auth.Service
	syncSvc trips.SyncService
	mailSvc email.Service
	store   Store
	logger  *zap.Logger
}

func NewService(
	authSvc auth.Service,
	syncSvc trips.SyncService,
	mailSvc email.Service,
	store Store,
	logger *zap.Logger,
) Service {
	return &service{
		authSvc,
		syncSvc,
		mailSvc,
		store,
		logger,
	}
}

// SendTripInvite sends a trip invite to an existing
// user. Such invite operation should be idempotent.
func (svc *service) SendTripInvite(
	ctx context.Context,
	tripID,
	authorID,
	userID string,
) error {
	inviteMeta, err := TripInviteMetaInfoFromCtx(ctx)
	if err != nil {
		return err
	}
	invite := NewInvite(
		tripID,
		inviteMeta.Trip.Name,
		authorID,
		inviteMeta.Author.Username,
		userID,
		inviteMeta.User.Email,
	)

	err = svc.store.SaveTripInvite(ctx, invite)
	if err == nil {
		go func() {
			svc.sendTripInviteEmail(ctx, invite, inviteMeta.Trip)
		}()
	}
	return err
}

func (svc *service) AcceptTripInvite(ctx context.Context, ID string) error {
	invite, err := svc.store.ReadTripInvite(ctx, ID)
	if err != nil {
		return err
	}

	connID := uuid.NewString()
	joinMsg := trips.MakeSyncMsgTOBTopicJoin(
		connID,
		invite.TripID,
		invite.AuthorID,
	)
	if err := svc.syncSvc.Join(ctx, &joinMsg); err != nil {
		return err
	}
	time.Sleep(syncMsgWaitInterval)

	member := trips.NewMember(invite.UserID, trips.MemberRoleCollaborator)
	addMemMsg := trips.MakeSyncMsgTOBTopicUpdate(
		connID,
		invite.TripID,
		invite.AuthorID,
		trips.SyncMsgTOBUpdateOpUpdateTripMembers,
		trips.MakeSyncMsgTOBUpdateOpUpdateTripMembersOps(member),
	)
	if err := svc.syncSvc.Update(ctx, &addMemMsg); err != nil {
		return err
	}

	return svc.store.DeleteTripInvite(ctx, ID)
}

func (svc *service) DeclineTripInvite(ctx context.Context, ID string) error {
	return svc.store.DeleteTripInvite(ctx, ID)
}

func (svc *service) ReadTripInvite(ctx context.Context, ID string) (TripInvite, error) {
	return svc.store.ReadTripInvite(ctx, ID)
}

func (svc *service) ListTripInvites(
	ctx context.Context,
	ff ListTripInvitesFilter,
) (TripInviteList, error) {
	return svc.store.ListTripInvites(ctx, ff)
}

func (svc *service) sendTripInviteEmail(
	ctx context.Context,
	invite TripInvite,
	trip *trips.Trip,
) {
	svc.logger.Info("sending trip invite email", zap.String("to", invite.UserEmail))
	t, err := template.
		New(tripInviteTmplFileName).
		ParseFiles(tripInviteTmplFilePath)
	if err != nil {
		svc.logger.Error("sendTripInviteEmail", zap.Error(err))
		return
	}

	coverImgURL := defaultCoverImgURL
	if trip.CoverImage.Source == trips.CoverImageSourceWeb {
		coverImgURL = trip.CoverImage.WebImage.Urls.Regular
	}

	var doc bytes.Buffer
	data := struct {
		ID          string
		AuthorName  string
		TripName    string
		CoverImgURL string
	}{
		invite.ID,
		invite.AuthorName,
		invite.TripName,
		coverImgURL,
	}
	if err := t.Execute(&doc, data); err != nil {
		svc.logger.Error("sendTripInviteEmail", zap.Error(err))
		return
	}

	mailContentBody := doc.String()
	mailBody, err := svc.mailSvc.InsertContentOnTemplate(mailContentBody)
	if err != nil {
		svc.logger.Error("sendTripInviteEmail", zap.Error(err))
		return
	}

	subj := "New Trip Invite!"
	if err := svc.mailSvc.SendMail(
		ctx,
		invite.UserEmail,
		defaultLoginSender,
		subj,
		mailBody,
	); err != nil {
		svc.logger.Error("sendTripInviteEmail", zap.Error(err))
	}
}

// Email Trip Invite

// SendEmailTripInvite sends a email trip invite to a new or existing
// user. Such invite operation should be idempotent.
func (svc service) SendEmailTripInvite(
	ctx context.Context,
	tripID,
	authorID,
	userEmail string,
) error {
	inviteMeta, err := TripInviteMetaInfoFromCtx(ctx)
	if err != nil {
		return err
	}

	inv := NewEmailTripInvite(
		tripID,
		inviteMeta.Trip.Name,
		authorID,
		inviteMeta.Author.Username,
		userEmail,
	)

	err = svc.store.SaveEmailTripInvite(ctx, inv)
	if err == nil {
		go func() {
			c, sig, err := svc.authSvc.GenerateOTPAuthCodeAndSig(
				ctx,
				userEmail,
				emailInviteDuration,
			)
			if err != nil {
				svc.logger.Error("svc.authSvc.GenerateOTPAuthCodeAndSig", zap.Error(err))
			}

			svc.sendEmailTripInviteEmail(ctx, inv, inviteMeta.Trip, sig, c)
		}()
	}

	return err
}

func (svc service) AcceptEmailTripInvite(
	ctx context.Context,
	ID,
	sig,
	c string,
	isLoggedIn bool,
) (*http.Cookie, error) {
	usr, cookie, err := svc.authSvc.EmailLogin(ctx, c, sig, isLoggedIn)
	if err != nil {
		return nil, err
	}

	invite, err := svc.store.ReadEmailTripInvite(ctx, ID)
	if err != nil {
		return nil, err
	}

	connID := uuid.NewString()
	joinMsg := trips.MakeSyncMsgTOBTopicJoin(
		connID,
		invite.TripID,
		invite.AuthorID,
	)
	if err := svc.syncSvc.Join(ctx, &joinMsg); err != nil {
		return nil, err
	}
	time.Sleep(syncMsgWaitInterval)

	member := trips.NewMember(usr.ID, trips.MemberRoleCollaborator)
	addMemMsg := trips.MakeSyncMsgTOBTopicUpdate(
		connID,
		invite.TripID,
		invite.AuthorID,
		trips.SyncMsgTOBUpdateOpUpdateTripMembers,
		trips.MakeSyncMsgTOBUpdateOpUpdateTripMembersOps(member),
	)
	if err := svc.syncSvc.Update(ctx, &addMemMsg); err != nil {
		return nil, err
	}

	go svc.store.DeleteTripInvite(ctx, ID)
	return cookie, nil
}

func (svc service) ReadEmailTripInvite(ctx context.Context, ID string) (EmailTripInvite, error) {
	return svc.store.ReadEmailTripInvite(ctx, ID)

}

func (svc service) ListEmailTripInvites(
	ctx context.Context,
	ff ListEmailTripInvitesFilter,
) (EmailTripInviteList, error) {
	return svc.store.ListEmailTripInvites(ctx, ff)
}

func (svc *service) sendEmailTripInviteEmail(
	ctx context.Context,
	inv EmailTripInvite,
	trip *trips.Trip,
	sig,
	c string,
) {
	svc.logger.Info("sending email trip invite email", zap.String("to", inv.UserEmail))
	t, err := template.
		New(tripInviteTmplFileName).
		ParseFiles(tripInviteTmplFilePath)
	if err != nil {
		svc.logger.Error("sendEmailTripInviteEmail", zap.Error(err))
		return
	}

	coverImgURL := defaultCoverImgURL
	if trip.CoverImage.Source == trips.CoverImageSourceWeb {
		coverImgURL = trip.CoverImage.WebImage.Urls.Regular
	}

	var doc bytes.Buffer
	data := struct {
		ID          string
		AuthorName  string
		TripName    string
		CoverImgURL string
		Sig         string
		Code        string
	}{
		inv.ID,
		inv.AuthorName,
		inv.TripName,
		coverImgURL,
		sig,
		c,
	}
	if err := t.Execute(&doc, data); err != nil {
		svc.logger.Error("sendEmailTripInviteEmail", zap.Error(err))
		return
	}

	mailContentBody := doc.String()
	mailBody, err := svc.mailSvc.InsertContentOnTemplate(mailContentBody)
	if err != nil {
		svc.logger.Error("sendEmailTripInviteEmail", zap.Error(err))
		return
	}

	subj := "New Trip Invite!"
	if err := svc.mailSvc.SendMail(
		ctx,
		inv.UserEmail,
		defaultLoginSender,
		subj,
		mailBody,
	); err != nil {
		svc.logger.Error("sendEmailTripInviteEmail", zap.Error(err))
	}
}
