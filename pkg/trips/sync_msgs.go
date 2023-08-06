package trips

import "github.com/travelreys/travelreys/pkg/jsonpatch"

// SessionContext represents a user's participation
// to the multiplayer collaboration session.
type SessionContext struct {
	// ConnID tracks the connection from an instance of
	// the travelreys client app.
	ConnID string

	// TripID represents the trip currently being updated
	TripID string

	// Participating member
	Member Member
}

type SessionContextList []SessionContext

const (
	SyncMsgTypeControl = "control"
	SyncMsgTypeData    = "data"
)

type SyncMsg struct {
	Type   string `json:"type"`
	ConnID string `json:"connId"`
	TripID string `json:"tripId"`

	Member Member `json:"member"`
}

const (
	SyncMsgControlTopicJoin        = "SyncMsgControlTopicJoin"
	SyncMsgControlTopicLeave       = "SyncMsgControlTopicLeave"
	SyncMsgControlTopicPing        = "SyncMsgControlTopicPing"
	SyncMsgControlTopiCursor       = "SyncMsgControlTopicCursor"
	SyncMsgControlTopiFormPresence = "SyncMsgControlTopicFormPresence"
)

type SyncMsgControl struct {
	SyncMsg

	Topic string `json:"topic"`

	JoinSession  SyncMsgControlPayloadJoin  `json:"join"`
	LeaveSession SyncMsgControlPayloadLeave `json:"leave"`
	Ping         SyncMsgControlPayloadPing  `json:"ping"`
}

type SyncMsgControlPayloadJoin struct {
	// Latest Snapshot of the trip
	Trip *Trip `json:"trip"`

	// List of members updated (presence)
	Members MembersList `json:"members"`
}

type SyncMsgControlPayloadLeave struct {
}

type SyncMsgControlPayloadPing struct {
}

type SyncMsgControlPayloadCursor struct{}

type SyncMsgControlPayloadFormPresence struct {
	IsActive bool   `json:"isActive"` // toggle on/off
	EditPath string `json:"editPath"`
}

func MakeSyncMsgControlTopicLeave(
	connID,
	tripID string,
	mem Member,
) SyncMsgControl {
	return SyncMsgControl{
		SyncMsg: SyncMsg{
			Type:   SyncMsgTypeControl,
			ConnID: connID,
			TripID: tripID,
			Member: mem,
		},
		Topic: SyncMsgControlTopicLeave,
	}
}

func MakeSyncMsgControlTopicPing(
	connID,
	tripID string,
	mem Member,
) SyncMsgControl {
	return SyncMsgControl{
		SyncMsg: SyncMsg{
			Type:   SyncMsgTypeControl,
			ConnID: connID,
			TripID: tripID,
			Member: mem,
		},
		Topic: SyncMsgControlTopicPing,
	}
}

const (
	// Trip
	SyncMsgDataTopicDeleteTrip        = "SyncMsgDataTopicDeleteTrip"
	SyncMsgDataTopicUpdateTripDates   = "SyncMsgDataTopicUpdateTripDates"
	SyncMsgDataTopicUpdateTripMembers = "SyncMsgDataTopicUpdateTripMembers"

	// Lodgings
	SyncMsgDataTopicAddLodging    = "SyncMsgDataTopicAddLodging"
	SyncMsgDataTopicDeleteLodging = "SyncMsgDataTopicDeleteLodging"
	SyncMsgDataTopicUpdateLodging = "SyncMsgDataTopicUpdateLodging"

	// Itinerary
	SyncMsgDataTopicDeleteActivity              = "SyncMsgDataTopicDeleteActivity"
	SyncMsgDataTopicOptimizeItinerary           = "SyncMsgDataTopicOptimizeItinerary"
	SyncMsgDataTopicReorderActivityToAnotherDay = "SyncMsgDataTopicReorderActivityToAnotherDay"
	SyncMsgDataTopicReorderItinerary            = "SyncMsgDataTopicReorderItinerary"
	SyncMsgDataTopicUpdateActivityPlace         = "SyncMsgDataTopicUpdateActivityPlace"

	// Media
	SyncMsgDataTopicAddMediaItem = "SyncMsgDataTopicAddMediaItem"
)

type SyncMsgData struct {
	SyncMsg

	Topic   string             `json:"topic"`
	Counter uint64             `json:"counter"`
	Payload SyncMsgDataPayload `json:"payload"`
}

type SyncMsgDataPayload struct {
	Ops []jsonpatch.Op `json:"ops"`
}
