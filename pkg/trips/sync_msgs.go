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
	MemberID string
}

type SessionContextList []SessionContext

func (l SessionContextList) ToMembers() []string {
	mList := []string{}
	for _, ctx := range l {
		mList = append(mList, ctx.MemberID)
	}
	return mList
}

const (
	SyncMsgTypeControl = "control"
	SyncMsgTypeData    = "data"
)

type SyncMsg struct {
	Type     string `json:"type"`
	ConnID   string `json:"connId"`
	TripID   string `json:"tripId"`
	MemberID string `json:"memberId"`
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

	Join  SyncMsgControlPayloadJoin  `json:"join"`
	Leave SyncMsgControlPayloadLeave `json:"leave"`
	Ping  SyncMsgControlPayloadPing  `json:"ping"`
}

type SyncMsgControlPayloadJoin struct {
	// Latest Snapshot of the trip
	Trip *Trip `json:"trip"`

	// List of members updated (presence)
	Members []string `json:"members"`
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
	mem string,
) SyncMsgControl {
	return SyncMsgControl{
		SyncMsg: SyncMsg{
			Type:     SyncMsgTypeControl,
			ConnID:   connID,
			TripID:   tripID,
			MemberID: mem,
		},
		Topic: SyncMsgControlTopicLeave,
	}
}

func MakeSyncMsgControlTopicPing(
	connID,
	tripID string,
	mem string,
) SyncMsgControl {
	return SyncMsgControl{
		SyncMsg: SyncMsg{
			Type:     SyncMsgTypeControl,
			ConnID:   connID,
			TripID:   tripID,
			MemberID: mem,
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

	Topic   string         `json:"topic"`
	Counter uint64         `json:"counter"`
	Ops     []jsonpatch.Op `json:"ops"`
}
