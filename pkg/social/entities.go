package social

import (
	"fmt"

	"github.com/google/uuid"
)

type FriendRequest struct {
	ID          string `json:"id"`
	InitiatorID string `json:"initiatorID"`
	TargetID    string `json:"targetID"`
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
	ID            string `json:"id"`
	BindingKey    string `json:"binding"`
	RevBindingKey string `json:"revbinding"`
	UserOneID     string `json:"userOneID"`
	UserTwoID     string `json:"userTwoID"`
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
