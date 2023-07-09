package social

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/travelreys/travelreys/pkg/auth"
	"github.com/travelreys/travelreys/pkg/common"
	"github.com/travelreys/travelreys/pkg/trips"
)

const (
	LabelVerified = "verified"
)

type UserProfile struct {
	ID         string        `json:"id"`
	Username   string        `json:"username"`
	ProfileImg string        `json:"profileImg"`
	Labels     common.Labels `json:"labels"`
}

func UserProfileFromUser(user auth.User) UserProfile {
	return UserProfile{
		ID:         user.ID,
		Username:   user.Username,
		ProfileImg: user.GetProfileImgURL(),
		Labels: common.Labels{
			LabelVerified: user.Labels[LabelVerified],
		},
	}
}

func (up UserProfile) IsVerified() bool {
	return up.Labels[LabelVerified] == "true"
}

type UserProfileMap map[string]UserProfile

type FriendRequest struct {
	ID          string `json:"id" bson:"id"`
	BindingKey  string `json:"binding" bson:"binding"`
	InitiatorID string `json:"initiatorID" bson:"initiatorID"`
	TargetID    string `json:"targetID" bson:"targetID"`

	InitiatorProfile UserProfile `json:"initiatorProfile" bson:"-"`
}

func NewFriendRequest(initiator, target string) FriendRequest {
	return FriendRequest{
		ID:          uuid.NewString(),
		BindingKey:  fmt.Sprintf("%s|%s", initiator, target),
		InitiatorID: initiator,
		TargetID:    target,
	}
}

type FriendRequestList []FriendRequest

func (reqList FriendRequestList) GetInitiatorIDs() []string {
	ids := []string{}
	for _, req := range reqList {
		ids = append(ids, req.InitiatorID)
	}
	return ids
}

type Friend struct {
	ID          string `json:"id" bson:"id"`
	BindingKey  string `json:"binding" bson:"binding"`
	InitiatorID string `json:"initiatorID" bson:"initiatorID"`
	TargetID    string `json:"targetID" bson:"targetID"`

	InitiatorProfile UserProfile `json:"initiatorProfile" bson:"-"`
	TargetProfile    UserProfile `json:"targetProfile" bson:"-"`
}

func NewFriendFromRequest(req FriendRequest) Friend {
	return Friend{
		ID:          uuid.NewString(),
		BindingKey:  fmt.Sprintf("%s|%s", req.InitiatorID, req.TargetID),
		InitiatorID: req.InitiatorID,
		TargetID:    req.TargetID,
	}
}

type FriendsList []Friend

func (l FriendsList) GetTargetIDs() []string {
	ids := []string{}
	for _, req := range l {
		ids = append(ids, req.TargetID)
	}
	return ids
}

func (l FriendsList) GetInitiatorIDs() []string {
	ids := []string{}
	for _, req := range l {
		ids = append(ids, req.InitiatorID)
	}
	return ids
}

func MakeTripPublicInfo(trip *trips.Trip) trips.Trip {
	newTrip := trips.NewTrip(trip.Creator, trip.Name)
	newTrip.ID = trip.ID
	newTrip.CoverImage = trip.CoverImage
	newTrip.Lodgings = trip.Lodgings
	for key, lod := range newTrip.Lodgings {
		lod.CheckinTime = time.Time{}
		lod.CheckoutTime = time.Time{}
		newTrip.Lodgings[key] = lod
	}
	newTrip.MediaItems = trip.MediaItems

	newTrip.Itineraries = map[string]trips.Itinerary{}
	sortedItinKey := trips.GetSortedItineraryKeys(trip)
	for idx, key := range sortedItinKey {
		itin := trip.Itineraries[key]
		itin.Date = time.Time{}
		newActivities := map[string]trips.Activity{}
		for aKey, act := range trip.Itineraries[key].Activities {
			act.StartTime = time.Time{}
			act.EndTime = time.Time{}
			newActivities[aKey] = act
		}
		itin.Activities = newActivities
		newTrip.Itineraries[fmt.Sprintf("%d", idx)] = itin
	}
	if _, ok := trip.Labels[trips.LabelSharingAccess]; ok {
		newTrip.Labels[trips.LabelSharingAccess] = trip.Labels[trips.LabelSharingAccess]
	}
	return newTrip
}

func MakeTripPublicInfoWithUserProfiles(trip *trips.Trip, profiles UserProfileMap) trips.Trip {
	newTrip := MakeTripPublicInfo(trip)

	if _, ok := profiles[trip.Creator.ID]; ok {
		newTrip.Members[trip.Creator.ID] = trips.Member{}
	}
	for key := range trip.Members {
		if _, ok := profiles[key]; ok {
			newTrip.Members[key] = trips.Member{}
		}
	}
	return newTrip
}
