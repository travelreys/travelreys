package social

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"html/template"
	"os"

	"github.com/travelreys/travelreys/pkg/auth"
	"github.com/travelreys/travelreys/pkg/email"
	"go.uber.org/zap"
)

const (
	defaultLoginSender                = "login@travelreys.com"
	defaultFriendReqEmailTmplFilePath = "assets/friendReqEmail.tmpl.html"
	defaultFriendReqEmailTmplFileName = "friendRequest.tmpl.html"
)

var (
	ErrInvalidFriendRequest    = errors.New("social.svc.InvalidFriendRequest")
	friendReqEmailTmplFilePath = os.Getenv("TRAVELREYS_FRIEND_REQ_EMAIL_PATH")
)

type Service interface {
	SendFriendRequest(ctx context.Context, initiatorID, targetID string) error
	GetFriendRequestByID(ctx context.Context, id string) (FriendRequest, error)
	AcceptFriendRequest(ctx context.Context, id string) error
	ListFriendRequests(ctx context.Context, ff ListFriendRequestsFilter) (FriendRequestList, error)
	DeleteFriendRequest(ctx context.Context, id string) error

	ListFriends(ctx context.Context, userID string) (FriendsList, error)
	DeleteFriend(ctx context.Context, userID, friendID string) error
	AreTheyFriends(ctx context.Context, userOneID, userTwoID string) error
}

type service struct {
	store   Store
	authSvc auth.Service
	mailSvc email.Service

	logger *zap.Logger
}

func NewService(
	store Store,
	authSvc auth.Service,
	mailSvc email.Service,
	logger *zap.Logger,
) Service {
	return &service{store, authSvc, mailSvc, logger}
}

func (svc service) SendFriendRequest(ctx context.Context, initiatorID, targetID string) error {
	userFF := auth.ListFilter{
		IDs: []string{initiatorID, targetID},
	}
	users, err := svc.authSvc.List(ctx, userFF)
	if err != nil {
		return err
	}
	if len(users) != 2 {
		return ErrInvalidFriendRequest
	}

	var (
		initiator auth.User
		target    auth.User
	)
	if users[0].ID == initiatorID {
		initiator = users[0]
		target = users[1]
	} else {
		initiator = users[1]
		target = users[0]
	}

	req := NewFriendRequest(initiatorID, targetID)
	if err := svc.store.InsertFriendRequest(ctx, req); err != nil {
		return err
	}
	go svc.sendFriendRequestEmail(ctx, initiator, target)
	return nil

}

func (svc service) AcceptFriendRequest(ctx context.Context, id string) error {
	var req FriendRequest
	info, err := FriendRequestInfoFromCtx(ctx)
	if err == nil {
		req = info.Req
	} else {
		req, err = svc.store.GetFriendRequestByID(ctx, id)
		if err != nil {
			return err
		}
	}

	if err := svc.store.SaveFriend(ctx, NewFriendFromRequest(req)); err != nil {
		return err
	}
	go svc.store.DeleteFriendRequest(ctx, id)
	return nil
}

func (svc service) GetFriendRequestByID(ctx context.Context, id string) (FriendRequest, error) {
	return svc.store.GetFriendRequestByID(ctx, id)
}

func (svc service) DeleteFriendRequest(ctx context.Context, id string) error {
	return svc.store.DeleteFriendRequest(ctx, id)
}

func (svc service) ListFriendRequests(ctx context.Context, ff ListFriendRequestsFilter) (FriendRequestList, error) {
	return svc.store.ListFriendRequests(ctx, ff)
}

func (svc service) ListFriends(ctx context.Context, userID string) (FriendsList, error) {
	return svc.store.ListFriends(ctx, userID)
}

func (svc service) DeleteFriend(ctx context.Context, userID, friendID string) error {
	return svc.store.DeleteFriend(ctx, fmt.Sprintf("%s|%s", userID, friendID))
}

func (svc service) AreTheyFriends(ctx context.Context, userOneID, userTwoID string) error {
	_, err := svc.store.GetFriendByUserIDs(ctx, userOneID, userTwoID)
	return err
}

func (svc service) sendFriendRequestEmail(ctx context.Context, initiator, target auth.User) {
	svc.logger.Info("sending friend req email", zap.String("to", target.Email))

	if friendReqEmailTmplFilePath == "" {
		friendReqEmailTmplFilePath = defaultFriendReqEmailTmplFileName
	}
	t, err := template.
		New(defaultFriendReqEmailTmplFileName).
		ParseFiles(defaultFriendReqEmailTmplFilePath)
	if err != nil {
		svc.logger.Error("sendFriendRequestEmail", zap.Error(err))
		return
	}

	var doc bytes.Buffer
	data := struct {
		InitiatorName string
		TargetName    string
	}{initiator.Name, target.Name}
	if err := t.Execute(&doc, data); err != nil {
		svc.logger.Error("sendFriendRequestEmail", zap.Error(err))
		return
	}

	mailBody := doc.String()
	subj := "New Friend Request!"
	if err := svc.mailSvc.SendMail(
		ctx, target.Email, defaultLoginSender, subj, mailBody,
	); err != nil {
		svc.logger.Error("sendFriendRequestEmail", zap.Error(err))
	}
}
