package tripssync

import (
	"fmt"

	"github.com/tiinyplanet/tiinyplanet/pkg/trips"
	"github.com/tiinyplanet/tiinyplanet/pkg/utils"
)

// Pub/Sub Subjects

func syncSessConnectionsKey(planID string) string {
	return fmt.Sprintf("sync-session:%s:connections", planID)
}

func syncSessCounterKey(planID string) string {
	return fmt.Sprintf("sync-session:%s:counter", planID)
}

// syncSessRequestSub is the subj for client -> coordinator
func syncSessRequestSubj(planID string) string {
	return fmt.Sprintf("sync-session:%s:requests", planID)
}

// syncSessTOBSubj is the subj for coordinator -> client
func syncSessTOBSubj(planID string) string {
	return fmt.Sprintf("sync-session:%s:tob", planID)
}

// Sync Session

type SyncSession struct {
	Members trips.TripMembersList `json:"members"` // who is in the current session
}

const (
	SyncOpJoinSession  = "SyncOpJoinSession"
	SyncOpLeaveSession = "SyncOpLeaveSession"
	SyncOpPingSession  = "SyncOpPingSession"
	SyncOpFetchTrip    = "SyncOpFetchTrip"
	SyncOpUpdateTrip   = "SyncOpUpdateTrip"
)

func isValidSyncOpType(opType string) bool {
	return utils.StringContains([]string{
		SyncOpJoinSession,
		SyncOpLeaveSession,
		SyncOpPingSession,
		SyncOpFetchTrip,
		SyncOpUpdateTrip,
	}, opType)
}

type SyncMessage struct {
	ID         string `json:"id"`      // Users' connection id
	Counter    uint64 `json:"counter"` // Should be monotonically increasing
	TripPlanID string `json:"tripPlanID"`
	OpType     string `json:"opType"`

	SyncDataJoinSession           `json:"syncDataJoinSession"`
	SyncDataJoinSessionBroadcast  `json:"syncDataJoinSessionBroadcast"`
	SyncDataLeaveSession          `json:"syncDataLeaveSession"`
	SyncDataLeaveSessionBroadcast `json:"syncDataLeaveSessionBroadcast"`
	SyncDataPing                  `json:"syncDataPing"`
	SyncDataUpdateTrip            `json:"syncDataUpdateTrip"`
}

type SyncDataJoinSession struct {
	trips.TripMember
}

type SyncDataJoinSessionBroadcast struct {
	trips.TripMembersList
}

type SyncDataLeaveSession struct {
	trips.TripMember
}

type SyncDataLeaveSessionBroadcast struct {
	trips.TripMembersList
}

type SyncDataPing struct{}

type SyncDataUpdateTrip struct {
	Op    string `json:"op"`
	Path  string `json:"path"`
	Value string `json:"value"` // JSON string
}

// Sync Connection

type SyncConnection struct {
	PlanID       string
	ConnectionID string
	Member       trips.TripMember
}
