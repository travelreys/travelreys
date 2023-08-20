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

	appInviteTmplFilePath       = "assets/appInviteEmail.tmpl.html"
	appInviteTmplFileName       = "appInviteEmail.tmpl.html"
	tripInviteTmplFilePath      = "assets/tripInviteEmail.tmpl.html"
	tripInviteTmplFileName      = "tripInviteEmail.tmpl.html"
	emailTripInviteTmplFilePath = "assets/emailTripInviteEmail.tmpl.html"
	emailTripInviteTmplFileName = "emailTripInviteEmail.tmpl.html"

	defaultCoverImgURL = "https://images.unsplash.com/photo-1476514525535-07fb3b4ae5f1?crop=entropy&cs=tinysrgb&fit=max&fm=jpg&ixid=M3wzOTc1ODV8MHwxfHNlYXJjaHwzfHx0cmF2ZWx8ZW58MHwwfHx8MTY4ODc0MDUxNnww&ixlib=rb-4.0.3&q=80&w=1080"

	appInviteDuration   = 720 * time.Hour // 30 days
	emailInviteDuration = 168 * time.Hour // 7 days
)

type Service interface {
	SendAppInvite(ctx context.Context, authorID, userEmail string) error
	AcceptAppInvite(ctx context.Context, ID, code, sig string) (auth.User, *http.Cookie, error)

	SendTripInvite(ctx context.Context, tripID, authorID, userID string) error
	AcceptTripInvite(ctx context.Context, ID string) error
	DeclineTripInvite(ctx context.Context, ID string) error
	ReadTripInvite(ctx context.Context, ID string) (TripInvite, error)
	ListTripInvites(ctx context.Context, ff ListTripInvitesFilter) (TripInviteList, error)

	SendEmailTripInvite(ctx context.Context, tripID, authorID, userEmail string) error
	AcceptEmailTripInvite(ctx context.Context, ID, code, sig string, isLoggedIn bool) (auth.User, *http.Cookie, error)
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
	return &service{authSvc, syncSvc, mailSvc, store, logger}
}

// App Invites

func (svc *service) SendAppInvite(ctx context.Context, authorID, userEmail string) error {
	inviteMeta, err := TripInviteMetaInfoFromCtx(ctx)
	if err != nil {
		return err
	}

	code, sig, err := svc.authSvc.GenerateOTPAuthCodeAndSig(
		ctx,
		userEmail,
		appInviteDuration,
	)
	if err != nil {
		return err
	}

	inv := NewAppInvite(
		authorID,
		inviteMeta.Author.Username,
		userEmail,
	)

	err = svc.store.SaveAppInvite(ctx, inv)
	if err == nil {
		go func() {
			svc.sendAppInviteEmail(
				context.Background(),
				inv,
				*inviteMeta.Author,
				code,
				sig,
			)
		}()
	}
	return err
}

func (svc *service) AcceptAppInvite(
	ctx context.Context,
	ID,
	code,
	sig string,
) (auth.User, *http.Cookie, error) {
	_, err := svc.store.ReadAppInvite(ctx, ID)
	if err != nil {
		return auth.User{}, nil, err
	}

	go svc.store.DeleteAppInvite(ctx, ID)
	return svc.authSvc.EmailLogin(ctx, code, sig, false)
}

func (svc *service) sendAppInviteEmail(
	ctx context.Context,
	inv AppInvite,
	author auth.User,
	code,
	sig string,
) {
	svc.logger.Info("sending app invite email", zap.String("to", inv.UserEmail))
	t, err := template.
		New(appInviteTmplFileName).
		ParseFiles(appInviteTmplFilePath)
	if err != nil {
		svc.logger.Error("sendAppInviteEmail", zap.Error(err))
		return
	}

	var doc bytes.Buffer
	data := struct {
		ID                  string
		AuthorName          string
		AuthorProfileImgURL string
		Code                string
		Sig                 string
	}{
		inv.ID,
		inv.AuthorName,
		author.GetProfileImgURL(),
		code,
		sig,
	}
	if err := t.Execute(&doc, data); err != nil {
		svc.logger.Error("sendAppInviteEmail", zap.Error(err))
		return
	}

	mailContentBody := doc.String()
	mailBody, err := svc.mailSvc.InsertContentOnTemplate(mailContentBody)
	if err != nil {
		svc.logger.Error("sendAppInviteEmail", zap.Error(err))
		return
	}

	subj := "Join Travelreys!"
	if err := svc.mailSvc.SendMail(
		ctx,
		inv.UserEmail,
		defaultLoginSender,
		subj,
		mailBody,
	); err != nil {
		svc.logger.Error("sendAppInviteEmail", zap.Error(err))
	}

}

// Trip Invites

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
	invite := NewTripInvite(
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
			svc.sendTripInviteEmail(
				context.Background(),
				invite,
				inviteMeta.Trip,
			)
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
				context.Background(),
				userEmail,
				emailInviteDuration,
			)
			if err != nil {
				svc.logger.Error("GenerateOTPAuthCodeAndSig", zap.Error(err))
			}

			svc.sendEmailTripInviteEmail(
				context.Background(), inv, inviteMeta.Trip, c, sig,
			)
		}()
	}
	return err
}

// AcceptEmailTripInvite accepts an email trip invitation.
// Such invite operation should be idempotent.
func (svc service) AcceptEmailTripInvite(
	ctx context.Context,
	ID,
	code,
	sig string,
	isLoggedIn bool,
) (auth.User, *http.Cookie, error) {
	invite, err := svc.store.ReadEmailTripInvite(ctx, ID)
	if err != nil {
		return auth.User{}, nil, err
	}

	usr, cookie, err := svc.authSvc.EmailLogin(ctx, code, sig, isLoggedIn)
	if err != nil {
		return auth.User{}, nil, err
	}

	connID := uuid.NewString()
	joinMsg := trips.MakeSyncMsgTOBTopicJoin(
		connID,
		invite.TripID,
		invite.AuthorID,
	)
	if err := svc.syncSvc.Join(ctx, &joinMsg); err != nil {
		return auth.User{}, nil, err
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
		return auth.User{}, nil, err
	}

	go svc.store.DeleteTripInvite(ctx, ID)
	return usr, cookie, nil
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
	code,
	sig string,
) {
	svc.logger.Info("sending email trip invite email", zap.String("to", inv.UserEmail))
	t, err := template.
		New(emailTripInviteTmplFileName).
		ParseFiles(emailTripInviteTmplFilePath)
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
		Code        string
		Sig         string
	}{
		inv.ID,
		inv.AuthorName,
		inv.TripName,
		coverImgURL,
		code,
		sig,
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
