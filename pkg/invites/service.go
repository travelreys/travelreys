package invites

import (
	"bytes"
	"context"
	"html/template"
	"time"

	"github.com/google/uuid"
	"github.com/travelreys/travelreys/pkg/email"
	"github.com/travelreys/travelreys/pkg/trips"
	"go.uber.org/zap"
)

const (
	SyncMsgWaitInterval    = 500 * time.Millisecond
	defaultLoginSender     = "login@travelreys.com"
	tripInviteTmplFilePath = "assets/tripInviteEmail.tmpl.html"
	tripInviteTmplFileName = "tripInviteEmail.tmpl.html"

	defaultCoverImgURL = "https://images.unsplash.com/photo-1476514525535-07fb3b4ae5f1?crop=entropy&cs=tinysrgb&fit=max&fm=jpg&ixid=M3wzOTc1ODV8MHwxfHNlYXJjaHwzfHx0cmF2ZWx8ZW58MHwwfHx8MTY4ODc0MDUxNnww&ixlib=rb-4.0.3&q=80&w=1080"
)

type Service interface {
	SendTripInvite(ctx context.Context, tripID, authorID, userID string) error
	AcceptTripInvite(ctx context.Context, ID string) error
	DeclineTripInvite(ctx context.Context, ID string) error
	ReadTripInvite(ctx context.Context, ID string) (TripInvite, error)
	ListTripInvites(ctx context.Context, ff ListTripInvitesFilter) (TripInviteList, error)
}

type inviteService struct {
	syncSvc trips.SyncService
	mailSvc email.Service
	store   Store
	logger  *zap.Logger
}

func NewService(
	syncSvc trips.SyncService,
	mailSvc email.Service,
	store Store,
	logger *zap.Logger,
) Service {
	return &inviteService{
		syncSvc,
		mailSvc,
		store,
		logger,
	}
}

func (svc *inviteService) SendTripInvite(
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

func (svc *inviteService) AcceptTripInvite(ctx context.Context, ID string) error {
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
	time.Sleep(SyncMsgWaitInterval)

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

func (svc *inviteService) DeclineTripInvite(ctx context.Context, ID string) error {
	return svc.store.DeleteTripInvite(ctx, ID)
}

func (svc *inviteService) ReadTripInvite(ctx context.Context, ID string) (TripInvite, error) {
	return svc.store.ReadTripInvite(ctx, ID)
}

func (svc *inviteService) ListTripInvites(
	ctx context.Context,
	ff ListTripInvitesFilter,
) (TripInviteList, error) {
	return svc.store.ListTripInvites(ctx, ff)
}

func (svc *inviteService) sendTripInviteEmail(
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
