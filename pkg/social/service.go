package social

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"html/template"
	"math/rand"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/travelreys/travelreys/pkg/auth"
	"github.com/travelreys/travelreys/pkg/common"
	"github.com/travelreys/travelreys/pkg/email"
	"github.com/travelreys/travelreys/pkg/images"
	"github.com/travelreys/travelreys/pkg/trips"
	"go.uber.org/zap"
)

const (
	defaultLoginSender         = "login@travelreys.com"
	friendReqEmailTmplFilePath = "assets/friendReqEmail.tmpl.html"
	friendReqEmailTmplFileName = "friendReqEmail.tmpl.html"
)

var (
	ErrInvalidFollowRequest = errors.New("social.ErrInvalidFollowRequest")
	ErrAlreadyFriends       = errors.New("social.ErrAlreadyFriends")
	ErrFollowRequestExists  = errors.New("social.ErrFollowRequestExists")
)

type Service interface {
	GetProfile(ctx context.Context, id string) (UserProfile, error)

	SendFollowRequest(ctx context.Context, initiatorID, targetID string) error
	GetFollowRequestByID(ctx context.Context, id string) (FollowRequest, error)
	AcceptFollowRequest(ctx context.Context, userID, initiatorID, reqID string) error
	ListFollowRequests(ctx context.Context, ff ListFollowRequestsFilter) (FollowRequestList, error)
	DeleteFollowRequest(ctx context.Context, userID, reqID string) error

	ListFollowers(ctx context.Context, userID string) (FollowingsList, error)
	ListFollowing(ctx context.Context, userID string) (FollowingsList, error)
	DeleteFollowing(ctx context.Context, userID, friendID string) error
	IsFollowing(ctx context.Context, initiatorID, targetID string) (bool, error)

	ReadTripPublicInfo(ctx context.Context, tripID, referrerID string) (*trips.Trip, UserProfile, error)
	ListTripPublicInfo(ctx context.Context, ff trips.ListFilter) (trips.TripsList, error)
	ListFollowingTrips(ctx context.Context, initiatorID string) (trips.TripsList, UserProfileMap, error)

	DuplicateTrip(ctx context.Context, initiatorID, referrerID, tripID, name string, startDate time.Time) (string, error)
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

func (svc *service) tripFromContext(ctx context.Context, ID string) (*trips.Trip, error) {
	var (
		trip *trips.Trip
		err  error
	)
	ti, err := trips.TripInfoFromCtx(ctx)
	if err == nil {
		trip = ti.Trip
	} else {
		trip, err = svc.tripSvc.Read(ctx, ID)
		if err != nil {
			return nil, err
		}
	}
	return trip, err
}

func (svc service) GetProfile(ctx context.Context, ID string) (UserProfile, error) {
	user, err := svc.authSvc.Read(ctx, ID)
	if err != nil {
		return UserProfile{}, err
	}
	return UserProfileFromUser(user), nil
}

func (svc service) SendFollowRequest(ctx context.Context, initiatorID, targetID string) error {
	// 1. Validate IDs exist
	userFF := auth.ListFilter{IDs: []string{initiatorID, targetID}}
	users, err := svc.authSvc.List(ctx, userFF)
	if err != nil {
		return err
	}
	if len(users) != 2 {
		return ErrInvalidFollowRequest
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

	// 2. Validate if following relationship already exists
	ok, err := svc.IsFollowing(ctx, initiatorID, targetID)
	if err != nil {
		return err
	}
	if ok {
		return ErrAlreadyFriends
	}

	reqs, err := svc.ListFollowRequests(ctx, ListFollowRequestsFilter{
		InitiatorID: common.StringPtr(initiatorID),
		TargetID:    common.StringPtr(targetID),
	})
	if err != nil {
		return nil
	}
	if len(reqs) > 0 {
		return ErrFollowRequestExists
	}

	req := NewFollowRequest(initiatorID, targetID)

	// 3. If target has verified, can automatically follow,
	// no need to wait for acceptance.
	targetProfile := UserProfileFromUser(target)
	if targetProfile.IsVerified() {
		return svc.store.SaveFollowing(ctx, NewFollowingFromRequest(req))
	}

	// 4. Else, send request
	if err := svc.store.UpsertFollowRequest(ctx, req); err != nil {
		return err
	}
	go svc.sendFollowRequestEmail(ctx, initiator, target, req)
	return nil
}

func (svc service) GetFollowRequestByID(ctx context.Context, id string) (FollowRequest, error) {
	return svc.store.GetFollowRequestByID(ctx, id)
}

func (svc service) AcceptFollowRequest(
	ctx context.Context,
	userID,
	initiatorID,
	reqID string,
) error {
	ok, err := svc.IsFollowing(ctx, initiatorID, userID)
	if err != nil {
		return err
	}
	if ok {
		return nil
	}

	var req FollowRequest
	info, err := FollowRequestInfoFromCtx(ctx)
	if err == nil {
		req = info.Req
	} else {
		req, err = svc.store.GetFollowRequestByID(ctx, reqID)
		if err != nil {
			return err
		}
	}
	if err := svc.DeleteFollowRequest(ctx, userID, reqID); err != nil {
		return err
	}
	return svc.store.SaveFollowing(ctx, NewFollowingFromRequest(req))
}

func (svc service) DeleteFollowRequest(ctx context.Context, userID, reqID string) error {
	return svc.store.DeleteFollowRequest(ctx, reqID)
}

func (svc service) ListFollowRequests(ctx context.Context, ff ListFollowRequestsFilter) (FollowRequestList, error) {
	reqs, err := svc.store.ListFollowRequests(ctx, ff)
	if err != nil {
		return nil, err
	}
	if len(reqs) == 0 {
		return FollowRequestList{}, nil
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

func (svc service) ListFollowers(ctx context.Context, userID string) (FollowingsList, error) {
	followers, err := svc.store.ListFollowers(ctx, userID)
	if err != nil {
		return nil, err
	}

	if len(followers) == 0 {
		return FollowingsList{}, nil
	}

	initiatorIDs := followers.GetInitiatorIDs()
	initiators, err := svc.authSvc.List(ctx, auth.ListFilter{IDs: initiatorIDs})
	if err != nil {
		return nil, err
	}
	profiles := map[string]UserProfile{}
	for _, usr := range initiators {
		profiles[usr.ID] = UserProfileFromUser(usr)
	}
	for i := 0; i < len(followers); i++ {
		followers[i].InitiatorProfile = profiles[followers[i].InitiatorID]
	}
	return followers, err
}

func (svc service) ListFollowing(ctx context.Context, userID string) (FollowingsList, error) {
	followings, err := svc.store.ListFollowing(ctx, userID)
	if err != nil {
		return nil, err
	}

	if len(followings) == 0 {
		return FollowingsList{}, nil
	}

	targets, err := svc.authSvc.List(
		ctx,
		auth.ListFilter{IDs: followings.GetTargetIDs()},
	)
	if err != nil {
		return nil, err
	}
	profiles := map[string]UserProfile{}
	for _, usr := range targets {
		profiles[usr.ID] = UserProfileFromUser(usr)
	}
	for i := 0; i < len(followings); i++ {
		followings[i].TargetProfile = profiles[followings[i].TargetID]
	}

	return followings, err
}

func (svc service) DeleteFollowing(ctx context.Context, userID, bindingKey string) error {
	return svc.store.DeleteFollowing(ctx, bindingKey)
}

func (svc service) IsFollowing(ctx context.Context, initiatorID, targetID string) (bool, error) {
	_, err := svc.store.GetFollowing(ctx, fmt.Sprintf("%s|%s", initiatorID, targetID))
	if err == ErrFollowingNotFound {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, nil
}

func (svc *service) ReadTripPublicInfo(ctx context.Context, tripID, referrerID string) (*trips.Trip, UserProfile, error) {
	trip, err := svc.tripFromContext(ctx, tripID)
	if err != nil {
		return nil, UserProfile{}, err
	}

	ff := auth.ListFilter{IDs: []string{referrerID}}
	users, err := svc.authSvc.List(ctx, ff)
	if err != nil {
		return trip, UserProfile{}, err
	}

	profile := UserProfileFromUser(users[0])
	pubInfo := MakeTripPublicInfo(trip)
	return pubInfo, profile, nil
}

func (svc *service) ListTripPublicInfo(ctx context.Context, ff trips.ListFilter) (trips.TripsList, error) {
	tripslist, err := svc.tripSvc.List(ctx, ff)
	if err != nil {
		return nil, err
	}

	publicInfo := trips.TripsList{}
	for _, t := range tripslist {
		publicInfo = append(publicInfo, MakeTripPublicInfo(t))
	}

	return publicInfo, nil
}

func (svc *service) ListFollowingTrips(ctx context.Context, initiatorID string) (trips.TripsList, UserProfileMap, error) {
	followings, err := svc.ListFollowing(ctx, initiatorID)
	if err != nil {
		return trips.TripsList{}, UserProfileMap{}, err
	}
	if len(followings) == 0 {
		return trips.TripsList{}, UserProfileMap{}, nil
	}

	profileMap := UserProfileMap{}
	for _, friend := range followings {
		profileMap[friend.TargetID] = friend.TargetProfile
	}

	tripslist, err := svc.tripSvc.List(
		ctx,
		trips.ListFilter{UserIDs: followings.GetTargetIDs()},
	)
	if err != nil {
		return trips.TripsList{}, UserProfileMap{}, err
	}
	publicInfo := trips.TripsList{}
	for _, t := range tripslist {
		publicInfo = append(publicInfo, MakeTripPublicInfoWithUserProfiles(t, profileMap))
	}
	return publicInfo, profileMap, nil
}

func (svc *service) DuplicateTrip(
	ctx context.Context,
	initiatorID,
	referrerID,
	tripID,
	name string,
	startDate time.Time,
) (string, error) {
	trip, err := svc.tripFromContext(ctx, tripID)
	if err != nil {
		return "", err
	}
	numDays := int(trip.EndDate.Sub(trip.StartDate).Hours() / 24)
	endDate := startDate.Add(time.Duration(numDays*24) * time.Hour)

	creator := trips.NewMember(initiatorID, trips.MemberRoleCreator)
	newTrip := trips.NewTripWithDates(creator, name, startDate, endDate)
	newTrip.CoverImage = &trips.CoverImage{
		Source:   trips.CoverImageSourceWeb,
		WebImage: images.CoverStockImageList[rand.Intn(len(images.CoverStockImageList))],
	}

	for _, lodging := range trip.Lodgings {
		numDays := int(lodging.CheckoutTime.Sub(lodging.CheckinTime).Hours() / 24)
		checkinDtDiff := int(lodging.CheckinTime.Sub(trip.StartDate).Hours() / 24)
		checkinDt := startDate
		if checkinDtDiff < 0 {
			checkinDt = checkinDt.Add(time.Duration(checkinDtDiff*-24) * time.Hour)
		} else {
			checkinDt = checkinDt.Add(time.Duration(checkinDtDiff*24) * time.Hour)
		}
		checkoutDt := checkinDt.Add(time.Duration(numDays*24) * time.Hour)

		newLodging := &trips.Lodging{
			ID:           uuid.NewString(),
			CheckinTime:  checkinDt,
			CheckoutTime: checkoutDt,
			PriceItem:    lodging.PriceItem,
			Notes:        lodging.Notes,
			Place:        lodging.Place,
			Labels: common.Labels{
				trips.LabelCreatedBy: initiatorID,
			},
		}
		newTrip.Lodgings[newLodging.ID] = newLodging
	}

	for _, itin := range trip.Itineraries {
		daysDiff := int(itin.Date.Sub(trip.StartDate).Hours() / 24)
		newDate := startDate.Add(time.Duration(daysDiff*24) * time.Hour)

		newItin := trips.NewItinerary(newDate)
		actIDMap := map[string]string{}

		for actKey, act := range itin.Activities {
			newAct := &trips.Activity{
				ID:        uuid.NewString(),
				Title:     act.Title,
				Place:     act.Place,
				Notes:     act.Notes,
				PriceItem: act.PriceItem,
				StartTime: act.StartTime,
				EndTime:   act.EndTime,
				Labels: common.Labels{
					trips.LabelCreatedBy:       initiatorID,
					trips.LabelFractionalIndex: act.Labels[trips.LabelFractionalIndex],
				},
			}
			newItin.Activities[actKey] = newAct
			actIDMap[act.ID] = newAct.ID
		}
		newItin.Labels[trips.LabelUiColor] = itin.Labels[trips.LabelUiColor]

		for rKey, route := range itin.Routes {
			tkns := strings.Split(rKey, "|")
			newKey := fmt.Sprintf("%s|%s", actIDMap[tkns[0]], actIDMap[tkns[1]])
			newItin.Routes[newKey] = route
		}

		newTrip.Itineraries[newDate.Format(trips.ItineraryDtKeyFormat)] = newItin
	}

	return newTrip.ID, svc.tripSvc.Save(ctx, newTrip)
}

func (svc service) sendFollowRequestEmail(
	ctx context.Context,
	initiator,
	target auth.User,
	req FollowRequest,
) {
	svc.logger.Info("sending friend req email", zap.String("to", target.Email))

	t, err := template.
		New(friendReqEmailTmplFileName).
		ParseFiles(friendReqEmailTmplFilePath)
	if err != nil {
		svc.logger.Error("sendFollowRequestEmail", zap.Error(err))
		return
	}

	var doc bytes.Buffer
	data := struct {
		InitiatorID            string
		InitiatorProfileImgURL string
		InitiatorName          string
		ReqID                  string
	}{
		initiator.ID,
		initiator.GetProfileImgURL(),
		initiator.Username,
		req.ID,
	}
	if err := t.Execute(&doc, data); err != nil {
		svc.logger.Error("sendFollowRequestEmail", zap.Error(err))
		return
	}

	mailContentBody := doc.String()
	mailBody, err := svc.mailSvc.InsertContentOnTemplate(mailContentBody)
	if err != nil {
		svc.logger.Error("sendFollowRequestEmail", zap.Error(err))
		return
	}

	subj := "New Follow Request!"
	if err := svc.mailSvc.SendMail(
		ctx, target.Email, defaultLoginSender, subj, mailBody,
	); err != nil {
		svc.logger.Error("sendFollowRequestEmail", zap.Error(err))
	}
}
