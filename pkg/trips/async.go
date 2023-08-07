package trips

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/travelreys/travelreys/pkg/jsonpatch"
)

type JoinTripFromEmailRequest struct {
	ID           string `json:"id" bson:"id"`
	Impersonatee string `json:"impersonatee" bson:"impersonatee"`
	Email        string `json:"email" bson:"email"`
}

type JoinTripFromMsgRequest struct {
	ID           string `json:"id" bson:"id"`
	Impersonatee string `json:"impersonatee" bson:"impersonatee"`
	TargetUserID string `json:"targetUserID" bson:"targetUserID"`
}

type AsyncService interface {
	SendJoinTripEmail(ctx context.Context, tripID, impersonatee, email string)
	JoinTripFromEmail(ctx context.Context, tripID, impersonatee, email string)

	SendJoinTripMsg(ctx context.Context, tripID, impersonatee, targetUserID string)
	JoinTripFromMsg(ctx context.Context, tripID, impersonatee, targetUserID string)
}

type asyncService struct {
	ctrlMsgStore SyncMsgControlStore
	dataMsgStore SyncMsgDataStore
}

// need to check if user has an account or not on frontend
// if have, once click

func (svc *asyncService) SendJoinTripEmail(ctx context.Context, tripID, impersonatee, email string) {

}

func (svc *asyncService) JoinTripFromMsg(
	ctx context.Context,
	tripID,
	impersonatee,
	targetUserID string,
) {
	connID := uuid.NewString()
	// 1. impersonate join
	joinMsg := MakeSyncMsgControlTopicJoin(
		connID,
		tripID,
		impersonatee,
	)
	svc.ctrlMsgStore.PubReq(tripID, joinMsg)
	time.Sleep(1 * time.Second)

	// 2. Send Add Member Op
	addMemMsg := MakeSyncMsgData(
		connID,
		tripID,
		impersonatee,
		SyncMsgDataTopicUpdateTripMembers,
		[]jsonpatch.Op{},
	)
	svc.dataMsgStore.PubReq(tripID, addMemMsg)
}
