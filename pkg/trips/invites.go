package trips

import (
	"errors"
	"fmt"

	"github.com/travelreys/travelreys/pkg/common"
)

var (
	ErrInvalidInvite = errors.New("trips.ErrInvalidInvite")
)

type Invite struct {
	ID         string        `json:"id" bson:"id"`
	AuthorID   string        `json:"authorID" bson:"authorID"`
	AuthorName string        `json:"authorName" bson:"-"`
	TripID     string        `json:"tripID" bson:"tripID"`
	TripName   string        `json:"tripName" bson:"-"`
	UserID     string        `json:"userID" bson:"userID"`
	UserEmail  string        `json:"userEmail" bson:"-"`
	Labels     common.Labels `json:"labels" bson:"labels"`
}

func NewInvite(
	tripID,
	tripName,
	authorID,
	authorName,
	userID,
	userEmail string,
) Invite {
	return Invite{
		ID:         fmt.Sprintf("%s|%s", tripID, userID),
		AuthorID:   authorID,
		AuthorName: authorName,
		TripID:     tripID,
		TripName:   tripName,
		UserID:     userID,
		UserEmail:  userEmail,
	}
}

type InviteList []Invite
