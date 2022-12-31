package trips

import (
	"time"

	"github.com/awhdesmond/tiinyplanet/pkg/common"
	"github.com/awhdesmond/tiinyplanet/pkg/utils"
	"github.com/google/uuid"
)

// Trip Plan

type TripPlan struct {
	ID   string `json:"id"`
	Name string `json:"name"`

	CoverImageURL string    `json:"coverImageURL"`
	StartDate     time.Time `json:"startDate"`
	EndDate       time.Time `json:"endDate"`

	Creator TripMember            `json:"creator"`
	Members map[string]TripMember `json:"members"`

	Contents map[string]TripContentList `json:"contents"` // Map of trip content's list
	Flights  map[string]Flight          `json:"flights"`
	Transits map[string]BaseTransit     `json:"transits"`
	Lodgings map[string]Lodging         `json:"lodgings"`

	UpdatedAt time.Time `json:"updatedAt"`
	CreatedAt time.Time `json:"createdAt"`

	IsArchived bool          `json:"isArchived"`
	Labels     common.Labels `json:"labels"`
	Tags       common.Tags   `json:"tags"`
}

type TripPlansList []TripPlan

func NewTripPlan(creator TripMember, name string) TripPlan {
	return TripPlan{
		ID:            uuid.New().String(),
		Name:          name,
		CoverImageURL: "",
		StartDate:     time.Time{},
		EndDate:       time.Time{},

		Creator: creator,
		Members: map[string]TripMember{},

		Contents: map[string]TripContentList{},
		Flights:  map[string]Flight{},
		Transits: map[string]BaseTransit{},
		Lodgings: map[string]Lodging{},

		UpdatedAt: time.Now(),
		CreatedAt: time.Now(),

		IsArchived: false,
		Labels:     common.Labels{},
		Tags:       common.Tags{},
	}
}

func NewTripPlanWithDates(creator TripMember, name string, start, end time.Time) TripPlan {
	plan := NewTripPlan(creator, name)
	plan.StartDate = start
	plan.EndDate = end
	return plan
}

func (tp *TripPlan) UpdateName(name string) {
	tp.Name = name
}

func (tp *TripPlan) UpdateCoverImageURL(url string) {
	tp.CoverImageURL = url
}

func (tp *TripPlan) UpdateMember(member TripMember) {
	tp.Members[member.MemberEmail] = member
}

// Members

const (
	TripMemberPermCollaborator = "collaborator"
	TripMemberPermParticipant  = "participant"
)

type TripMember struct {
	MemberID    string `json:"memberID"`
	MemberEmail string `json:"memberEmail"`
	Permission  string `json:"permission"`
}

// Transit

type TransitCenter struct {
	ID                string             `json:"id"`
	TransitCenterType string             `json:"transitCenterType"`
	Name              string             `json:"name"`
	Positioning       common.Positioning `json:"positioning"`
	Labels            common.Labels      `json:"labels"`
}

type BaseTransit struct {
	ID string `json:"id"`

	Cost           float64       `json:"cost"`
	ConfirmationID string        `json:"confirmationID"`
	Notes          string        `json:"notes"`
	Tags           common.Tags   `json:"tags"`
	Labels         common.Labels `json:"labels"`

	DepartTime         time.Time           `json:"departTime"`
	DepartPositioning  common.Positioning  `json:"departPositioning"`
	ArrivalTime        time.Time           `json:"arrivalTime"`
	ArrivalPositioning common.Positioning  `json:"arrivalPositioning"`
	Attachments        []common.FileObject `json:"attachments"`
}

// Flights

type Flight struct {
	BaseTransit
	AirplaneID string `json:"airplaneID"`
}

// Lodging

type Lodging struct {
	ID string `json:"id"`

	NumGuests int32 `json:"numGuests"`

	Cost           float64       `json:"cost"`
	ConfirmationID string        `json:"confirmationID"`
	Notes          string        `json:"notes"`
	Tags           common.Tags   `json:"tags"`
	Labels         common.Labels `json:"labels"`

	Positioning  common.Positioning `json:"positiioning"`
	CheckinTime  time.Time          `json:"checkinTime"`
	CheckoutTime time.Time          `json:"checkoutTime"`
}

// Trip Content

type TripContent struct {
	ID          string           `json:"id"`
	ContentType string           `json:"contentType"`
	Location    LocationContent  `json:"locationContent"`
	Notes       NoteContent      `json:"noteContent"`
	Checklist   ChecklistContent `json:"checklistContent"`

	Comments []TripContentComment `json:"comments"`
	Labels   common.Labels        `json:"labels"`
}

type TripContentList []TripContent

type LocationContent struct {
	Name        string             `json:"name"`
	Description string             `json:"description"`
	Positioning common.Positioning `json:"positioning"`
}

type NoteContent struct {
	ImageURL string `json:"imageURL"`
	Note     string `json:"note"`
}

type ChecklistContent struct {
	ImageURL string   `json:"imageURL"`
	Items    []string `json:"items"`
}

type TripContentComment struct {
	ID      string `json:"id"`
	Comment string `json:"comment"`

	Member TripMember `json:"member"`

	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// Itinerary

type ItineraryList struct {
	ID       string               `json:"id"`
	Name     string               `json:"name"`
	Date     time.Time            `json:"date"`
	Contents ItineraryContentList `json:"contents"`
}

type ItineraryContent struct {
	TripContent
	StartTime time.Time `json:"startTime"`
	EndTime   time.Time `json:"endTime"`
}

type ItineraryContentList []ItineraryContent

// Collaboration

type CollabSession struct {
	Members []TripMember `json:"members"`
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
	TripMember
}

type CollabOpLeaveSessionRequest struct {
	TripMember
}

type CollabOpPingSessionRequest struct{}

type CollabOpUpdateTripRequest struct {
	Op    string `json:"op"`
	Path  string `json:"path"`
	Value string `json:"value"` // JSON string
}
