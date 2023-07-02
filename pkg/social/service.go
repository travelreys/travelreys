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
	"github.com/travelreys/travelreys/pkg/trips"
	"go.uber.org/zap"
)

const (
	defaultLoginSender                = "login@travelreys.com"
	defaultFriendReqEmailTmplFilePath = "assets/friendReqEmail.tmpl.html"
	defaultFriendReqEmailTmplFileName = "friendReqEmail.tmpl.html"
)

var (
	ErrInvalidFriendRequest    = errors.New("social.svc.InvalidFriendRequest")
	ErrAlreadyFriends          = errors.New("social.svc.AlreadyFriends")
	friendReqEmailTmplFilePath = os.Getenv("TRAVELREYS_FRIEND_REQ_EMAIL_PATH")
)

type Service interface {
	GetProfile(ctx context.Context, id string) (UserProfile, error)

	SendFriendRequest(ctx context.Context, initiatorID, targetID string) error
	GetFriendRequestByID(ctx context.Context, id string) (FriendRequest, error)
	AcceptFriendRequest(ctx context.Context, userid, reqid string) error
	ListFriendRequests(ctx context.Context, ff ListFriendRequestsFilter) (FriendRequestList, error)
	DeleteFriendRequest(ctx context.Context, userid, reqid string) error

	ListFollowers(ctx context.Context, userID string) (FriendsList, error)
	ListFollowing(ctx context.Context, userID string) (FriendsList, error)
	DeleteFriend(ctx context.Context, userID, friendID string) error
	AreTheyFriends(ctx context.Context, initiatorID, targetID string) (bool, error)

	ReadTripPublicInfo(ctx context.Context, tripID, referrerID string) (trips.Trip, UserProfile, error)
	ListTripPublicInfo(ctx context.Context, ff trips.ListFilter) (trips.TripsList, error)
	ListFollowingTrips(ctx context.Context, initiatorID string) (trips.TripsList, UserProfileMap, error)
}

type service struct {
	store   Store
	authSvc auth.Service
	tripSvc trips.Service
	mailSvc email.Service

	logger *zap.Logger
}

func NewService(
	store Store,
	authSvc auth.Service,
	tripSvc trips.Service,
	mailSvc email.Service,
	logger *zap.Logger,
) Service {
	return &service{store, authSvc, tripSvc, mailSvc, logger}
}

func (svc service) GetProfile(ctx context.Context, ID string) (UserProfile, error) {
	user, err := svc.authSvc.Read(ctx, ID)
	if err != nil {
		return UserProfile{}, err
	}
	return UserProfileFromUser(user), nil
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

	ok, err := svc.AreTheyFriends(ctx, initiatorID, targetID)
	if err != nil {
		return err
	}
	if ok {
		return ErrAlreadyFriends
	}

	req := NewFriendRequest(initiatorID, targetID)
	if err := svc.store.UpsertFriendRequest(ctx, req); err != nil {
		return err
	}
	go svc.sendFriendRequestEmail(ctx, initiator, target)
	return nil
}

func (svc service) GetFriendRequestByID(ctx context.Context, id string) (FriendRequest, error) {
	return svc.store.GetFriendRequestByID(ctx, id)
}

func (svc service) AcceptFriendRequest(ctx context.Context, userid, reqid string) error {
	var req FriendRequest
	info, err := FriendRequestInfoFromCtx(ctx)
	if err == nil {
		req = info.Req
	} else {
		req, err = svc.store.GetFriendRequestByID(ctx, reqid)
		if err != nil {
			return err
		}
	}
	if err := svc.DeleteFriendRequest(ctx, userid, reqid); err != nil {
		return err
	}
	return svc.store.SaveFriend(ctx, NewFriendFromRequest(req))
}

func (svc service) DeleteFriendRequest(ctx context.Context, userid, reqid string) error {
	return svc.store.DeleteFriendRequest(ctx, reqid)
}

func (svc service) ListFriendRequests(ctx context.Context, ff ListFriendRequestsFilter) (FriendRequestList, error) {
	reqs, err := svc.store.ListFriendRequests(ctx, ff)
	if err != nil {
		return nil, err
	}

	initiatorIDs := reqs.GetInitiatorIDs()
	users, err := svc.authSvc.List(ctx, auth.ListFilter{IDs: initiatorIDs})
	if err != nil {
		return nil, err
	}
	profiles := map[string]UserProfile{}
	for _, usr := range users {
		profiles[usr.ID] = UserProfileFromUser(usr)
	}
	for i := 0; i < len(reqs); i++ {
		reqs[i].InitiatorProfile = profiles[reqs[i].InitiatorID]
	}
	return reqs, nil
}

func (svc service) ListFollowers(ctx context.Context, userID string) (FriendsList, error) {
	friends, err := svc.store.ListFollowers(ctx, userID)
	if err != nil {
		return nil, err
	}

	initiatorIDs := friends.GetInitiatorIDs()
	initiators, err := svc.authSvc.List(ctx, auth.ListFilter{IDs: initiatorIDs})
	if err != nil {
		return nil, err
	}
	profiles := map[string]UserProfile{}
	for _, usr := range initiators {
		profiles[usr.ID] = UserProfileFromUser(usr)
	}
	for i := 0; i < len(friends); i++ {
		friends[i].InitiatorProfile = profiles[friends[i].InitiatorID]
	}
	return friends, err
}

func (svc service) ListFollowing(ctx context.Context, userID string) (FriendsList, error) {
	friends, err := svc.store.ListFollowing(ctx, userID)
	if err != nil {
		return nil, err
	}

	targetIDs := friends.GetTargetIDs()
	targets, err := svc.authSvc.List(ctx, auth.ListFilter{IDs: targetIDs})
	if err != nil {
		return nil, err
	}
	profiles := map[string]UserProfile{}
	for _, usr := range targets {
		profiles[usr.ID] = UserProfileFromUser(usr)
	}
	for i := 0; i < len(friends); i++ {
		friends[i].TargetProfile = profiles[friends[i].TargetID]
	}

	return friends, err
}

func (svc service) DeleteFriend(ctx context.Context, userID, bindingKey string) error {
	return svc.store.DeleteFriend(ctx, bindingKey)
}

func (svc service) AreTheyFriends(ctx context.Context, initiatorID, targetID string) (bool, error) {
	_, err := svc.store.GetFriend(ctx, fmt.Sprintf("%s|%s", initiatorID, targetID))
	if err == ErrFriendNotFound {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, nil
}

func (svc service) sendFriendRequestEmail(ctx context.Context, initiator, target auth.User) {
	svc.logger.Info("sending friend req email", zap.String("to", target.Email))

	if friendReqEmailTmplFilePath == "" {
		friendReqEmailTmplFilePath = defaultFriendReqEmailTmplFilePath
	}
	t, err := template.
		New(defaultFriendReqEmailTmplFileName).
		ParseFiles(friendReqEmailTmplFilePath)
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

func (svc *service) ReadTripPublicInfo(ctx context.Context, tripID, referrerID string) (trips.Trip, UserProfile, error) {
	var (
		trip trips.Trip
		err  error
	)
	ti, err := trips.TripInfoFromCtx(ctx)
	if err == nil {
		trip = ti.Trip
	} else {
		trip, err = svc.tripSvc.Read(ctx, tripID)
		if err != nil {
			return trips.Trip{}, UserProfile{}, err
		}
	}

	ff := auth.ListFilter{IDs: []string{referrerID}}
	users, err := svc.authSvc.List(ctx, ff)
	if err != nil {
		return trip, UserProfile{}, err
	}

	profile := UserProfileFromUser(users[0])
	pubInfo := MakeTripPublicInfo(&trip)
	return pubInfo, profile, nil
}

func (svc *service) ListTripPublicInfo(ctx context.Context, ff trips.ListFilter) (trips.TripsList, error) {
	tripslist, err := svc.tripSvc.List(ctx, ff)
	if err != nil {
		return nil, err
	}

	publicInfo := trips.TripsList{}
	for _, t := range tripslist {
		publicInfo = append(publicInfo, MakeTripPublicInfo(&t))
	}

	return publicInfo, nil
}

func (svc *service) ListFollowingTrips(ctx context.Context, initiatorID string) (trips.TripsList, UserProfileMap, error) {
	followings, err := svc.ListFollowing(ctx, initiatorID)
	if err != nil {
		return trips.TripsList{}, UserProfileMap{}, err
	}
	targetIDs := followings.GetTargetIDs()
	if len(targetIDs) == 0 {
		return trips.TripsList{}, UserProfileMap{}, err
	}
	ff := auth.ListFilter{IDs: targetIDs}
	targets, err := svc.authSvc.List(ctx, ff)
	if err != nil {
		return trips.TripsList{}, UserProfileMap{}, err
	}
	profileMap := UserProfileMap{}
	for _, target := range targets {
		profileMap[target.ID] = UserProfileFromUser(target)
	}

	tripslist, err := svc.tripSvc.List(ctx, trips.ListFilter{UserIDs: targetIDs})
	if err != nil {
		return trips.TripsList{}, UserProfileMap{}, err
	}
	publicInfo := trips.TripsList{}
	for _, t := range tripslist {
		publicInfo = append(publicInfo, MakeTripPublicInfoWithUserProfiles(&t, profileMap))
	}

	return publicInfo, profileMap, nil
}
