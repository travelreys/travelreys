package tripssync

import (
	"encoding/json"
	"fmt"

	"github.com/awhdesmond/tiinyplanet/pkg/trips"
	"github.com/awhdesmond/tiinyplanet/pkg/utils"
)

// Pub/Sub Subjects

func collabSessMembersKey(planID string) string {
	return fmt.Sprintf("collab-session:%s:members", planID)
}

func collabSessUpdatesKey(planID string) string {
	return fmt.Sprintf("collab-session:%s:updates", planID)
}

func collabSessTOBKey(planID string) string {
	return fmt.Sprintf("collab-session:%s:tob", planID)
}

func collabSessCounterKey(planID string) string {
	return fmt.Sprintf("collab-session:%s:counter", planID)
}

// Collaboration Session

type CollabSession struct {
	Members trips.TripMembersList `json:"members"`
}

const (
	CollabOpJoinSession  = "CollabOpJoinSession"
	CollabOpLeaveSession = "CollabOpLeaveSession"
	CollabOpPingSession  = "CollabOpPingSession"
	CollabOpFetchTrip    = "CollabOpFetchTrip"
	CollabOpUpdateTrip   = "CollabOpUpdateTrip"
)

func isValidCollabOpType(opType string) bool {
	return utils.StringContains([]string{
		CollabOpJoinSession,
		CollabOpLeaveSession,
		CollabOpPingSession,
		CollabOpFetchTrip,
		CollabOpUpdateTrip,
	}, opType)
}

type CollabOpMessage struct {
	ID         string `json:"id"`
	Counter    uint64 `json:"ts"` // Should be monotonically increasing
	TripPlanID string `json:"tripPlanID"`
	OpType     string `json:"opType"`

	JoinSessionReq  CollabOpJoinSessionRequest  `json:"joinSessionReq"`
	LeaveSessionReq CollabOpLeaveSessionRequest `json:"leaveSessionReq"`
	PingSessionReq  CollabOpPingSessionRequest  `json:"pingSessionReq"`
	UpdateTripReq   CollabOpUpdateTripRequest   `json:"updateTripReq"`
}

type CollabOpJoinSessionRequest struct {
	trips.TripMember
}

type CollabOpLeaveSessionRequest struct {
	trips.TripMember
}

type CollabOpPingSessionRequest struct{}

type CollabOpUpdateTripRequest struct {
	Op    string `json:"op"`
	Path  string `json:"path"`
	Value string `json:"value"` // JSON string
}

func (req CollabOpUpdateTripRequest) Bytes() []byte {
	bytes, _ := json.Marshal(req)
	return bytes
}
