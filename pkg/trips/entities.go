package trips

import (
	"fmt"
	"sort"
	"time"

	"github.com/google/uuid"
	"github.com/tiinyplanet/tiinyplanet/pkg/common"
	"github.com/tiinyplanet/tiinyplanet/pkg/flights"
	"github.com/tiinyplanet/tiinyplanet/pkg/images"
	"github.com/tiinyplanet/tiinyplanet/pkg/maps"
)

const (
	JSONPathLabelUiColor = "labels/ui|color"
	JSONPathLabelUiIcon  = "labels/ui|icon"

	LabelDelimeter              = "|"
	LabelFractionalIndex        = "fIndex"
	LabelLocked                 = "locked"
	LabelItineraryDates         = "itinerary|dates"
	LabelItineraryDatesJSONPath = "labels/itinerary|dates"
	LabelTransportModePref      = "transportationPreference"
	LabelUiColor                = "ui|color"
	LabelUiIcon                 = "ui|icon"
)

type Trip struct {
	ID   string `json:"id" bson:"id"`
	Name string `json:"name" bson:"name"`

	CoverImage images.ImageMetadata `json:"coverImage" bson:"coverImage"`
	StartDate  time.Time            `json:"startDate" bson:"startDate"`
	EndDate    time.Time            `json:"endDate" bson:"endDate"`

	// Members
	Creator Member            `json:"creator" bson:"creator"`
	Members map[string]Member `json:"members" bson:"members"`

	// Logistics
	Notes    string                 `json:"notes" bson:"notes"`
	Flights  map[string]Flight      `json:"flights" bson:"flights"`
	Transits map[string]BaseTransit `json:"transits" bson:"transits"`
	Lodgings map[string]Lodging     `json:"lodgings" bson:"lodgings"`

	// Activities
	Activities map[string]ActivityList `json:"activities" bson:"activities"`
	Itinerary  []ItineraryList         `json:"itinerary" bson:"itinerary"`

	// Budget
	Budget Budget `json:"budget" bson:"budget"`

	UpdatedAt  time.Time `json:"updatedAt" bson:"updatedAt"`
	CreatedAt  time.Time `json:"createdAt" bson:"createdAt"`
	IsArchived bool      `json:"isArchived" bson:"isArchived"`

	Labels common.Labels `json:"labels" bson:"labels"`
	Tags   common.Tags   `json:"tags" bson:"tags"`
}

type TripsList []Trip

func NewTrip(creator Member, name string) Trip {
	creator.Role = MemberRoleCreator

	return Trip{
		ID:         uuid.New().String(),
		Name:       name,
		CoverImage: images.ImageMetadata{},
		StartDate:  time.Time{},
		EndDate:    time.Time{},

		Creator: creator,
		Members: map[string]Member{},

		Flights:  map[string]Flight{},
		Transits: map[string]BaseTransit{},
		Lodgings: map[string]Lodging{},

		Activities: map[string]ActivityList{},
		Itinerary:  []ItineraryList{},
		Budget:     NewBudget(),

		UpdatedAt: time.Now(),
		CreatedAt: time.Now(),

		IsArchived: false,
		Labels:     common.Labels{},
		Tags:       common.Tags{},
	}
}

func NewTripWithDates(creator Member, name string, start, end time.Time) Trip {
	plan := NewTrip(creator, name)
	plan.StartDate = start
	plan.EndDate = end
	return plan
}

// Members

const (
	MemberRoleCreator      = "creator"
	MemberRoleCollaborator = "collaborator"
	MemberRoleParticipant  = "participant"
)

type Member struct {
	ID   string `json:"id" bson:"id"`
	Role string `json:"role" bson:"role"`

	Labels map[string]string `json:"labels" bson:"labels"`
}

type MembersList []Member

func NewMember(id, role string) Member {
	return Member{
		ID:     id,
		Role:   role,
		Labels: map[string]string{},
	}
}

const (
	TransitTypeFlight = "flight"
)

type BaseTransit struct {
	ID   string `json:"id" bson:"id"`
	Type string `json:"type"`

	ConfirmationID string       `json:"confirmationID" bson:"confirmationID"`
	Notes          string       `json:"notes" bson:"notes"`
	Price          common.Price `json:"price" bson:"price"`

	Tags        common.Tags         `json:"tags" bson:"tags"`
	Labels      common.Labels       `json:"labels" bson:"labels"`
	Attachments []common.FileObject `json:"attachments" bson:"attachments"`
}

type Flight struct {
	BaseTransit
	ItineraryType string         `json:"itineraryType" bson:"itineraryType"`
	Depart        flights.Flight `json:"depart" bson:"depart"`
	Return        flights.Flight `json:"return" bson:"return"`
}

type Lodging struct {
	ID string `json:"id" bson:"id"`

	NumGuests      int32        `json:"numGuests" bson:"numGuests"`
	CheckinTime    time.Time    `json:"checkinTime" bson:"checkinTime"`
	CheckoutTime   time.Time    `json:"checkoutTime" bson:"checkoutTime"`
	Price          common.Price `json:"price" bson:"price"`
	ConfirmationID string       `json:"confirmationID" bson:"confirmationID"`
	Notes          string       `json:"notes" bson:"notes"`
	Place          maps.Place   `json:"place" bson:"place"`

	Tags        common.Tags         `json:"tags" bson:"tags"`
	Labels      common.Labels       `json:"labels" bson:"labels"`
	Attachments []common.FileObject `json:"attachments" bson:"attachments"`
}

type Activity struct {
	ID    string `json:"id" bson:"id"`
	Title string `json:"title" bson:"title"`

	Place maps.Place `json:"place" bson:"place"`
	Notes string     `json:"notes" bson:"notes"`

	Comments []ActivityComment `json:"comments" bson:"comments"`
	Labels   common.Labels     `json:"labels" bson:"labels"`
}

func (a Activity) HasPlace() bool {
	return a.Place.Name != ""
}

type ActivityList struct {
	ID         string              `json:"id" bson:"id"`
	Name       string              `json:"name" bson:"name"`
	Activities map[string]Activity `json:"activities" bson:"activities"`
	Labels     common.Labels       `json:"labels" bson:"labels"`
}

func NewActivityList(name string) ActivityList {
	return ActivityList{
		ID:         uuid.New().String(),
		Name:       name,
		Activities: map[string]Activity{},
		Labels:     common.Labels{},
	}
}

type ActivityComment struct {
	ID        string    `json:"id" bson:"id"`
	Comment   string    `json:"comment" bson:"comment"`
	Member    Member    `json:"member" bson:"member"`
	CreatedAt time.Time `json:"createdAt" bson:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt" bson:"updatedAt"`
}

type ItineraryActivity struct {
	ID             string        `json:"id" bson:"id"`
	ActivityListID string        `json:"activityListId" bson:"activityListId"`
	ActivityID     string        `json:"activityId" bson:"activityId"`
	Price          common.Price  `json:"price" bson:"price"`
	StartTime      time.Time     `json:"startTime" bson:"startTime"`
	EndTime        time.Time     `json:"endTime" bson:"endTime"`
	Labels         common.Labels `json:"labels" bson:"labels"`
}

type ItineraryActivityList []ItineraryActivity

func (l ItineraryActivityList) Len() int {
	return len(l)
}
func (l ItineraryActivityList) Swap(i, j int) {
	l[i], l[j] = l[j], l[i]
}
func (l ItineraryActivityList) Less(i, j int) bool {
	return l[i].Labels[LabelFractionalIndex] < l[j].Labels[LabelFractionalIndex]
}

type ItineraryList struct {
	ID          string                       `json:"id" bson:"id"`
	Date        time.Time                    `json:"date" bson:"date"`
	Description string                       `json:"desc" bson:"desc"`
	Activities  map[string]ItineraryActivity `json:"activities" bson:"activities"`
	Routes      map[string]maps.RouteList    `json:"routes" bson:"routes"`
	Labels      common.Labels                `json:"labels" bson:"labels"`
}

func NewItineraryList(date time.Time) ItineraryList {
	return ItineraryList{
		ID:          uuid.New().String(),
		Date:        date,
		Description: "",
		Activities:  map[string]ItineraryActivity{},
		Routes:      map[string]maps.RouteList{},
		Labels:      common.Labels{},
	}
}

// SortActivities returns a list of ItineraryActivities sorted
// by their fractional index
func (l ItineraryList) SortActivities() []ItineraryActivity {
	sorted := ItineraryActivityList{}
	for _, act := range l.Activities {
		sorted = append(sorted, act)
	}
	sort.Sort(sorted)
	return sorted
}

func GetFracIndexes(acts []ItineraryActivity) []string {
	result := []string{}
	for _, a := range acts {
		result = append(result, a.Labels[LabelFractionalIndex])
	}
	return result
}

func (l ItineraryList) routePairingKey(a1 ItineraryActivity, a2 ItineraryActivity) string {
	return fmt.Sprintf("%s%s%s", a1.ID, LabelDelimeter, a2.ID)
}

func (l ItineraryList) MakeRoutePairings() map[string]bool {
	pairings := map[string]bool{}
	sorted := l.SortActivities()
	for i := 1; i < len(sorted); i++ {
		pairings[l.routePairingKey(sorted[i-1], sorted[i])] = true
	}
	return pairings
}

type Budget struct {
	ID     string          `json:"id" bson:"id"`
	Amount common.Price    `json:"amount" bson:"amount"`
	Items  BudgetItemsList `json:"items" bson:"items"`

	Labels common.Labels `json:"labels" bson:"labels"`
	Tags   common.Tags   `json:"tags" bson:"tags"`
}

func NewBudget() Budget {
	return Budget{
		Amount: common.Price{},
		Items:  BudgetItemsList{},
		Labels: common.Labels{},
		Tags:   common.Tags{},
	}
}

type BudgetItem struct {
	Title string       `json:"title" bson:"title"`
	Desc  string       `json:"desc" bson:"desc"`
	Price common.Price `json:"price" bson:"price"`

	Labels common.Labels `json:"labels" bson:"labels"`
	Tag    common.Tags   `json:"tags" bson:"tags"`
}

type BudgetItemsList []BudgetItem
