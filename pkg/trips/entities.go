package trips

import (
	"time"

	"github.com/awhdesmond/tiinyplanet/pkg/common"
)

// Trip Plan

type TripPlan struct {
	Name      string    `json:"name"`
	ImageURL  string    `json:"imageURL"`
	StartDate time.Time `json:"startDate"`
	EndDate   time.Time `json:"endDate"`

	CreatorID string                     `json:"creatorID"` // ID of user that created the plan
	Contents  map[string]TripContentList `json:"contents"`  // Map of trip content's list

	Flights  map[string]Flight      `json:"flights"`
	Transits map[string]BaseTransit `json:"transits"`
	Lodgings map[string]Lodging     `json:"lodgings"`

	UpdatedAt time.Time `json:"updatedAt"`
	CreatedAt time.Time `json:"createdAt"`
}

// Transit

type TransitCenter struct {
	ID                string             `json:"id"`
	TransitCenterType string             `json:"transitCenterType"`
	Name              string             `json:"name"`
	Positioning       common.Positioning `json:"positioning"`
	Labels            map[string]string  `json:"labels"`
}

type BaseTransit struct {
	ID string `json:"id"`

	Cost           float64           `json:"cost"`
	ConfirmationID string            `json:"confirmationID"`
	Notes          string            `json:"notes"`
	Tags           map[string]string `json:"tags"`
	Labels         map[string]string `json:"labels"`

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

	Cost           float64           `json:"cost"`
	ConfirmationID string            `json:"confirmationID"`
	Notes          string            `json:"notes"`
	Tags           map[string]string `json:"tags"`
	Labels         map[string]string `json:"labels"`

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
