package tripssync

import (
	"github.com/tiinyplanet/tiinyplanet/pkg/trips"
	"github.com/tiinyplanet/tiinyplanet/pkg/utils"
)

// Sync Session State

type SyncConnection struct {
	PlanID       string
	ConnectionID string
	Member       trips.TripMember
}

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

func NewSyncMessageJoinSessionBroadcast(tripPlanID string, members trips.TripMembersList) SyncMessage {
	return SyncMessage{
		TripPlanID:                   tripPlanID,
		OpType:                       SyncOpJoinSessionBroadcast,
		SyncDataJoinSessionBroadcast: SyncDataJoinSessionBroadcast{members},
	}
}

type SyncDataLeaveSession struct {
	trips.TripMember
}

func NewSyncMessageLeaveSession(connID, tripPlanID string) SyncMessage {
	return SyncMessage{
		OpType:     SyncOpLeaveSession,
		ID:         connID,
		TripPlanID: tripPlanID,
	}
}

type SyncDataLeaveSessionBroadcast struct {
	trips.TripMembersList
}

func NewSyncMessageLeaveSessionBroadcast(tripPlanID string, members trips.TripMembersList) SyncMessage {
	return SyncMessage{
		TripPlanID:                    tripPlanID,
		OpType:                        SyncOpLeaveSessionBroadcast,
		SyncDataLeaveSessionBroadcast: SyncDataLeaveSessionBroadcast{members},
	}
}

type SyncDataPing struct{}

type SyncDataUpdateTrip struct {
	Op    string `json:"op"`
	Path  string `json:"path"`
	Value string `json:"value"` // JSON string
}
