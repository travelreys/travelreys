package social

import (
	"fmt"

	"github.com/google/uuid"
)

type FriendRequest struct {
	ID          string `json:"id" bson:"id"`
	InitiatorID string `json:"initiatorID" bson:"initiatorID"`
	TargetID    string `json:"targetID" bson:"targetID"`
}

func NewFriendRequest(initiator, target string) FriendRequest {
	return FriendRequest{
		ID:          uuid.NewString(),
		InitiatorID: initiator,
		TargetID:    target,
	}
}

type FriendRequestList []FriendRequest

type Friend struct {
	ID            string `json:"id" bson:"id"`
	BindingKey    string `json:"binding" bson:"binding"`
	RevBindingKey string `json:"revbinding" bson:"revbinding"`
	UserOneID     string `json:"userOneID" bson:"userOneID"`
	UserTwoID     string `json:"userTwoID" bson:"userTwoID"`
}

func NewFriendFromRequest(req FriendRequest) Friend {
	return Friend{
		ID:            uuid.NewString(),
		BindingKey:    fmt.Sprintf("%s|%s", req.InitiatorID, req.TargetID),
		RevBindingKey: fmt.Sprintf("%s|%s", req.TargetID, req.InitiatorID),
		UserOneID:     req.InitiatorID,
		UserTwoID:     req.TargetID,
	}
}

type FriendsList []Friend

type Statistic struct {
	ID         string `json:"id"`
	NumFriends uint64 `json:"numFriends"`
}
