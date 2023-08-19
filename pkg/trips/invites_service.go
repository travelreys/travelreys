package trips

import (
	"bytes"
	"context"
	"html/template"
	"time"

	"github.com/google/uuid"
	"github.com/travelreys/travelreys/pkg/email"
	"go.uber.org/zap"
)

const (
	SyncMsgWaitInterval    = 500 * time.Millisecond
	defaultLoginSender     = "login@travelreys.com"
	tripInviteTmplFilePath = "assets/tripInviteEmail.tmpl.html"
	tripInviteTmplFileName = "tripInviteEmail.tmpl.html"
)

type InviteService interface {
	Send(ctx context.Context, tripID, authorID, userID string) error
	Accept(ctx context.Context, ID string) error
	Decline(ctx context.Context, ID string) error
	Read(ctx context.Context, ID string) (Invite, error)
	List(ctx context.Context, ff ListInvitesFilter) (InviteList, error)
}

type inviteService struct {
	syncSvc SyncService
	mailSvc email.Service
	store   InviteStore
	logger  *zap.Logger
}

func NewInviteService(
	syncSvc SyncService,
	mailSvc email.Service,
	store InviteStore,
	logger *zap.Logger,
) InviteService {
	return &inviteService{
		syncSvc,
		mailSvc,
		store,
		logger,
	}
}

func (svc *inviteService) Send(
	ctx context.Context,
	tripID,
	authorID,
	userID string,
) error {
	inviteMeta, err := InviteMetaInfoFromCtx(ctx)
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

	err = svc.store.Save(ctx, invite)
	if err == nil {
		go func() {
			svc.sendTripInviteEmail(ctx, invite)
		}()
	}
	return err
}

func (svc *inviteService) Accept(ctx context.Context, ID string) error {
	invite, err := svc.store.Read(ctx, ID)
	if err != nil {
		return err
	}

	connID := uuid.NewString()
	joinMsg := MakeSyncMsgTOBTopicJoin(
		connID,
		invite.TripID,
		invite.AuthorID,
	)
	if err := svc.syncSvc.Join(ctx, &joinMsg); err != nil {
		return err
	}
	time.Sleep(SyncMsgWaitInterval)

	member := NewMember(invite.UserID, MemberRoleCollaborator)
	addMemMsg := MakeSyncMsgTOBTopicUpdate(
		connID,
		invite.TripID,
		invite.AuthorID,
		SyncMsgTOBUpdateOpUpdateTripMembers,
		MakeSyncMsgTOBUpdateOpUpdateTripMembersOps(member),
	)
	if err := svc.syncSvc.Update(ctx, &addMemMsg); err != nil {
		return err
	}

	return svc.store.Delete(ctx, ID)
}

func (svc *inviteService) Decline(ctx context.Context, ID string) error {
	return svc.store.Delete(ctx, ID)
}

func (svc *inviteService) Read(ctx context.Context, ID string) (Invite, error) {
	return svc.store.Read(ctx, ID)
}

func (svc *inviteService) List(ctx context.Context, ff ListInvitesFilter) (InviteList, error) {
	return svc.store.List(ctx, ff)
}

func (svc *inviteService) sendTripInviteEmail(
	ctx context.Context,
	invite Invite,
) {
	svc.logger.Info("sending trip invite email", zap.String("to", invite.UserEmail))
	t, err := template.
		New(tripInviteTmplFileName).
		ParseFiles(tripInviteTmplFilePath)
	if err != nil {
		svc.logger.Error("sendTripInviteEmail", zap.Error(err))
		return
	}

	var doc bytes.Buffer
	data := struct {
		ID         string
		AuthorName string
		TripName   string
	}{
		invite.ID,
		invite.AuthorName,
		invite.TripName,
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
