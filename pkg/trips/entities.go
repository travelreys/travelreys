package trips

import (
	"time"

	"github.com/awhdesmond/tiinyplanet/pkg/common"
	"github.com/google/uuid"
)

// Trip Plan

type TripPlan struct {
	ID            string    `json:"id"`
	Name          string    `json:"name"`
	CreatorID     string    `json:"creatorID"` // ID of user that created the plan
	CoverImageURL string    `json:"coverImageURL"`
	StartDate     time.Time `json:"startDate"`
	EndDate       time.Time `json:"endDate"`

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

func NewTripPlan(creatorID, name string) TripPlan {
	return TripPlan{
		ID:            uuid.New().String(),
		Name:          name,
		CreatorID:     creatorID,
		CoverImageURL: "",
		StartDate:     time.Time{},
		EndDate:       time.Time{},

		Contents: map[string]TripContentList{},
		Flights:  map[string]Flight{},
		Transits: map[string]BaseTransit{},
		Lodgings: map[string]Lodging{},

		UpdatedAt: time.Now(),
		CreatedAt: time.Now(),

		Labels: common.Labels{},
		Tags:   common.Tags{},
	}
}

func NewTripPlanWithDates(creatorID, name string, start, end time.Time) TripPlan {
	plan := NewTripPlan(creatorID, name)
	plan.StartDate = start
	plan.EndDate = end
	return plan
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
}

type TripContentList []TripContent

type LocationContent struct {
	Name        string             `json:"name"`
	Description string             `json:"description"`
	Positioning common.Positioning `json:"positioning"`
}

type NoteContent struct {
	Note string `json:"note"`
}

type ChecklistContent struct {
	Items []string `json:"items"`
}

type TripContentComment struct {
	ID      string `json:"id"`
	Comment string `json:"comment"`

	AuthorID    string `json:"authorID"`
	AuthorName  string `json:"authorName"`
	AuthorEmail string `json:"authorEmail"`

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

type CollabOpType int32

const (
	JoinSession CollabOpType = iota
	LeaveSession
	PingSession
	FetchTrip
	UpdateTrip
)

type CollabOpMessage struct {
	ID         string       `json:"id"`
	TripPlanID string       `json:"tripPlanID"`
	OpType     CollabOpType `json:"opType"`

	CollabOpJoinSessionRequest  CollabOpJoinSessionRequest  `json:"collabOpJoinSessionRequest"`
	CollabOpLeaveSessionRequest CollabOpLeaveSessionRequest `json:"collabOpLeaveSessionRequest"`
	CollabOpPingSessionRequest  CollabOpPingSessionRequest  `json:"collabOpPingSessionRequest"`
	CollabOpUpdateTripRequest   CollabOpUpdateTripRequest   `json:"collabOpUpdateTripRequest"`
}

type CollabOpJoinSessionRequest struct {
	CollaboratorID    string `json:"collaboratorID"`
	CollaboratorName  string `json:"collaboratorName"`
	CollaboratorEmail string `json:"collaboratorEmail"`
}

type CollabOpLeaveSessionRequest struct {
	CollaboratorID string `json:"collaboratorID"`
}

type CollabOpPingSessionRequest struct{}

type CollabOpFetchTripRequest struct{}

type CollabOpUpdateTripRequest struct {
	CollabOpUpdateTripPatch
}

type CollabOpUpdateTripPatch struct {
	Op    string `json:"op"`
	Path  string `json:"path"`
	Value string `json:"value"` // JSON string
}
