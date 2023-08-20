package invites

import (
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/travelreys/travelreys/pkg/common"
)

var (
	ErrInvalidInvite = errors.New("invites.ErrInvalidInvite")
)

type AppInvite struct {
	ID         string        `json:"id" bson:"id"`
	AuthorID   string        `json:"authorID" bson:"authorID"`
	AuthorName string        `json:"authorName" bson:"-"`
	UserEmail  string        `json:"userEmail" bson:"-"`
	CreatedAt  time.Time     `json:"createdAt" bson:"createdAt"`
	Labels     common.Labels `json:"labels" bson:"labels"`
}

func NewAppInvite(
	authorID,
	authorName,
	userEmail string,
) AppInvite {
	return AppInvite{
		ID:         uuid.NewString(),
		AuthorID:   authorID,
		AuthorName: authorName,
		UserEmail:  userEmail,
		CreatedAt:  time.Now(),
		Labels:     common.Labels{},
	}
}

type AppInvitesList []AppInvite

type TripInvite struct {
	ID         string        `json:"id" bson:"id"`
	AuthorID   string        `json:"authorID" bson:"authorID"`
	AuthorName string        `json:"authorName" bson:"-"`
	TripID     string        `json:"tripID" bson:"tripID"`
	TripName   string        `json:"tripName" bson:"-"`
	UserID     string        `json:"userID" bson:"userID"`
	UserEmail  string        `json:"userEmail" bson:"-"`
	Labels     common.Labels `json:"labels" bson:"labels"`
}

func NewTripInvite(
	tripID,
	tripName,
	authorID,
	authorName,
	userID,
	userEmail string,
) TripInvite {
	return TripInvite{
		ID:         fmt.Sprintf("%s|%s", tripID, userID),
		AuthorID:   authorID,
		AuthorName: authorName,
		TripID:     tripID,
		TripName:   tripName,
		UserID:     userID,
		UserEmail:  userEmail,
		Labels:     common.Labels{},
	}
}

type TripInviteList []TripInvite

type EmailTripInvite struct {
	ID         string        `json:"id" bson:"id"`
	AuthorID   string        `json:"authorID" bson:"authorID"`
	AuthorName string        `json:"authorName" bson:"-"`
	TripID     string        `json:"tripID" bson:"tripID"`
	TripName   string        `json:"tripName" bson:"-"`
	UserEmail  string        `json:"userEmail" bson:"-"`
	CreatedAt  time.Time     `json:"createdAt" bson:"createdAt"`
	Labels     common.Labels `json:"labels" bson:"labels"`
}

func NewEmailTripInvite(
	tripID,
	tripName,
	authorID,
	authorName,
	userEmail string,
) EmailTripInvite {
	return EmailTripInvite{
		ID:         fmt.Sprintf("%s|%s", tripID, userEmail),
		AuthorID:   authorID,
		AuthorName: authorName,
		TripID:     tripID,
		TripName:   tripName,
		UserEmail:  userEmail,
		CreatedAt:  time.Now(),
		Labels:     common.Labels{},
	}
}

type EmailTripInviteList []EmailTripInvite
