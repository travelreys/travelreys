package tripssync

import (
	"github.com/tiinyplanet/tiinyplanet/pkg/trips"
	"github.com/tiinyplanet/tiinyplanet/pkg/utils"
)

// Sync Session

type SyncSession struct {
	// Members is a list of members in the current session
	Members trips.TripMembersList `json:"members"`
}

// Sync Message

const (
	SyncOpJoinSession           = "SyncOpJoinSession"
	SyncOpJoinSessionBroadcast  = "SyncOpJoinSessionBroadcast"
	SyncOpLeaveSession          = "SyncOpLeaveSession"
	SyncOpLeaveSessionBroadcast = "SyncOpLeaveSessionBroadcast"
	SyncOpPingSession           = "SyncOpPingSession"
	SyncOpFetchTrip             = "SyncOpFetchTrip"
	SyncOpUpdateTrip            = "SyncOpUpdateTrip"
)

func isValidSyncMessageType(opType string) bool {
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
