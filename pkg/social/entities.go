package social

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/travelreys/travelreys/pkg/auth"
	"github.com/travelreys/travelreys/pkg/common"
	"github.com/travelreys/travelreys/pkg/trips"
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
	}
}

type FriendRequest struct {
	ID          string `json:"id" bson:"id"`
	InitiatorID string `json:"initiatorID" bson:"initiatorID"`
	TargetID    string `json:"targetID" bson:"targetID"`

	InitiatorProfile UserProfile `json:"initiatorProfile" bson:"-"`
}

func NewFriendRequest(initiator, target string) FriendRequest {
	return FriendRequest{
		ID:          uuid.NewString(),
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

type Statistic struct {
	ID         string `json:"id"`
	NumFriends uint64 `json:"numFriends"`
}

func MakeTripPublicInfo(trip *trips.Trip) trips.Trip {
	newTrip := trips.NewTrip(trip.Creator, trip.Name)
	newTrip.CoverImage = trip.CoverImage
	newTrip.Lodgings = trip.Lodgings
	newTrip.MediaItems = trip.MediaItems

	newTrip.Itineraries = map[string]trips.Itinerary{}
	sortedItinKey := trips.GetSortedItineraryKeys(trip)
	for idx, key := range sortedItinKey {
		newTrip.Itineraries[fmt.Sprintf("%d", idx)] = trip.Itineraries[key]
	}
	return newTrip
}
