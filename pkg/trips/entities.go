package trips

import (
	"fmt"
	"sort"
	"time"

	"github.com/google/uuid"
	"github.com/travelreys/travelreys/pkg/auth"
	"github.com/travelreys/travelreys/pkg/common"
	"github.com/travelreys/travelreys/pkg/images"
	"github.com/travelreys/travelreys/pkg/maps"
	"github.com/travelreys/travelreys/pkg/media"
	"github.com/travelreys/travelreys/pkg/ogp"
	"github.com/travelreys/travelreys/pkg/storage"
)

const (
	JSONPathLabelUiColor = "labels/ui|color"
	JSONPathLabelUiIcon  = "labels/ui|icon"

	LabelCreatedBy       = "createdBy"
	LabelDelimeter       = "|"
	LabelFractionalIndex = "fIndex"
	LabelSharingAccess   = "sharing|access"
	LabelUiColor         = "ui|color"
	LabelUiIcon          = "ui|icon"
)

const (
	SharingAccessViewer  = "view"
	SharingAccessComment = "comment"
)

type Trip struct {
	ID   string `json:"id" bson:"id"`
	Name string `json:"name" bson:"name"`

	CoverImage CoverImage `json:"coverImage" bson:"coverImage"`
	StartDate  time.Time  `json:"startDate" bson:"startDate"`
	EndDate    time.Time  `json:"endDate" bson:"endDate"`

	// Members
	Creator   Member            `json:"creator" bson:"creator"`
	Members   map[string]Member `json:"members" bson:"members"`
	MembersID map[string]string `json:"membersId" bson:"membersId"`

	// Logistics
	Notes    string                 `json:"notes" bson:"notes"`
	Transits map[string]BaseTransit `json:"transits" bson:"transits"`
	Lodgings map[string]Lodging     `json:"lodgings" bson:"lodgings"`
	Budget   Budget                 `json:"budget" bson:"budget"`
	Links    map[string]Link        `json:"links" bson:"links"`

	Itineraries map[string]Itinerary `json:"itineraries" bson:"itineraries"`

	// Media, Attachements
	MediaItems map[string]media.MediaItemMap `json:"mediaItems" bson:"mediaItems"`
	Files      map[string]storage.Object     `json:"files" bson:"files"`

	// To be deprecated:
	Media map[string]storage.Object `json:"media" bson:"media"`

	UpdatedAt time.Time `json:"updatedAt" bson:"updatedAt"`
	CreatedAt time.Time `json:"createdAt" bson:"createdAt"`

	Labels common.Labels `json:"labels" bson:"labels"`
	Tags   common.Tags   `json:"tags" bson:"tags"`
}

type TripOGP struct {
	ID   string `json:"id"`
	Name string `json:"name"`

	CoverImage CoverImage `json:"coverImage" bson:"coverImage"`
	Creator    auth.User  `json:"creator"`
}

type TripsList []Trip

func NewTrip(creator Member, name string) Trip {
	creator.Role = MemberRoleCreator

	return Trip{
		ID:          uuid.New().String(),
		Name:        name,
		CoverImage:  CoverImage{},
		StartDate:   time.Time{},
		EndDate:     time.Time{},
		Creator:     creator,
		Members:     map[string]Member{},
		MembersID:   map[string]string{},
		Transits:    map[string]BaseTransit{},
		Lodgings:    map[string]Lodging{},
		Itineraries: map[string]Itinerary{},
		Budget:      NewBudget(),
		Links:       LinkMap{},
		MediaItems:  map[string]media.MediaItemMap{},
		Media:       map[string]storage.Object{},
		Files:       map[string]storage.Object{},
		UpdatedAt:   time.Now(),
		CreatedAt:   time.Now(),
		Labels:      common.Labels{},
		Tags:        common.Tags{},
	}
}

func NewTripWithDates(creator Member, name string, start, end time.Time) Trip {
	plan := NewTrip(creator, name)
	plan.StartDate = start
	plan.EndDate = end
	return plan
}

func (trip Trip) PublicInfo() Trip {
	newTrip := NewTrip(trip.Creator, trip.Name)
	newTrip.CoverImage = trip.CoverImage
	newTrip.Lodgings = trip.Lodgings
	newTrip.Itineraries = trip.Itineraries
	return newTrip
}

func (trip Trip) OGP(creator auth.User) TripOGP {
	return TripOGP{
		ID:         trip.ID,
		Name:       trip.Name,
		CoverImage: trip.CoverImage,
		Creator:    creator,
	}
}

func (trip Trip) IsSharingEnabled() bool {
	_, ok := trip.Labels[LabelSharingAccess]
	return ok
}

const (
	CoverImageSourceWeb  = "web"
	CoverImageSourceTrip = "trip"
)

type CoverImage struct {
	Source   string               `json:"source" bson:"source"`
	WebImage images.ImageMetadata `json:"webImage" bson:"webImage"`

	// TripImage is the ID of mediaItem in the mediaItems["coverImage"]
	TripImage string `json:"tripImage" bson:"tripImage"`
}

// Members

const (
	MemberRoleCreator      = "creator"
	MemberRoleCollaborator = "collaborator"
	MemberRoleParticipant  = "participant"
)

type Member struct {
	ID     string            `json:"id" bson:"id"`
	Role   string            `json:"role" bson:"role"`
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

func (t Trip) GetAllMembersID() []string {
	membersIDs := []string{t.Creator.ID}
	for id := range t.Members {
		membersIDs = append(membersIDs, id)
	}
	return membersIDs
}

const (
	TransitTypeFlight = "flight"
	TransitTypeTrain  = "train"
	TransitTypeBus    = "bus"
	TransitTypeOthers = "others"
)

type BaseTransit struct {
	ID              string           `json:"id" bson:"id"`
	Type            string           `json:"type"`
	DepartTime      time.Time        `json:"departTime" bson:"departTime"`
	DepartLocation  maps.Place       `json:"departLocation" bson:"departLocation"`
	ArrivalTime     time.Time        `json:"arrivalTime" bson:"arrivalTime"`
	ArrivalLocation maps.Place       `json:"arrivalLocation" bson:"arrivalLocation"`
	ConfirmationID  string           `json:"confirmationID" bson:"confirmationID"`
	Notes           string           `json:"notes" bson:"notes"`
	PriceItem       common.PriceItem `json:"price" bson:"price"`

	Tags   common.Tags   `json:"tags" bson:"tags"`
	Labels common.Labels `json:"labels" bson:"labels"`
}

type Lodging struct {
	ID             string           `json:"id" bson:"id"`
	NumGuests      int32            `json:"numGuests" bson:"numGuests"`
	CheckinTime    time.Time        `json:"checkinTime" bson:"checkinTime"`
	CheckoutTime   time.Time        `json:"checkoutTime" bson:"checkoutTime"`
	PriceItem      common.PriceItem `json:"price" bson:"price"`
	ConfirmationID string           `json:"confirmationID" bson:"confirmationID"`
	Notes          string           `json:"notes" bson:"notes"`
	Place          maps.Place       `json:"place" bson:"place"`
	Tags           common.Tags      `json:"tags" bson:"tags"`
	Labels         common.Labels    `json:"labels" bson:"labels"`
}

type ActivityComment struct {
	ID        string    `json:"id" bson:"id"`
	Comment   string    `json:"comment" bson:"comment"`
	Member    Member    `json:"member" bson:"member"`
	CreatedAt time.Time `json:"createdAt" bson:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt" bson:"updatedAt"`
}

type Activity struct {
	ID        string            `json:"id" bson:"id"`
	Title     string            `json:"title" bson:"title"`
	Place     maps.Place        `json:"place" bson:"place"`
	Notes     string            `json:"notes" bson:"notes"`
	PriceItem common.PriceItem  `json:"price" bson:"price"`
	StartTime time.Time         `json:"startTime" bson:"startTime"`
	EndTime   time.Time         `json:"endTime" bson:"endTime"`
	Comments  []ActivityComment `json:"comments" bson:"comments"`
	Labels    common.Labels     `json:"labels" bson:"labels"`
}

func (a Activity) HasPlace() bool {
	return a.Place.Name != ""
}

type ActivityList []Activity

func (l ActivityList) Len() int {
	return len(l)
}
func (l ActivityList) Swap(i, j int) {
	l[i], l[j] = l[j], l[i]
}
func (l ActivityList) Less(i, j int) bool {
	return l[i].Labels[LabelFractionalIndex] < l[j].Labels[LabelFractionalIndex]
}

type Itinerary struct {
	ID          string                    `json:"id" bson:"id"`
	Date        time.Time                 `json:"date" bson:"date"`
	Description string                    `json:"desc" bson:"desc"`
	Activities  map[string]Activity       `json:"activities" bson:"activities"`
	Routes      map[string]maps.RouteList `json:"routes" bson:"routes"`
	Labels      common.Labels             `json:"labels" bson:"labels"`
}

func NewItinerary(date time.Time) Itinerary {
	return Itinerary{
		ID:          uuid.New().String(),
		Date:        date,
		Description: "",
		Activities:  map[string]Activity{},
		Routes:      map[string]maps.RouteList{},
		Labels:      common.Labels{},
	}
}

func GetSortedItineraryKeys(trip *Trip) []string {
	list := []string{}
	for key := range trip.Itineraries {
		list = append(list, key)
	}
	sort.Sort(sort.StringSlice(list))
	return list
}

// SortActivities returns Activities sorted by their fractional index
func (l Itinerary) SortActivities() []Activity {
	sorted := ActivityList{}
	for _, act := range l.Activities {
		sorted = append(sorted, act)
	}
	sort.Sort(sorted)
	return sorted
}

func GetFracIndexes(acts []Activity) []string {
	result := []string{}
	for _, a := range acts {
		result = append(result, a.Labels[LabelFractionalIndex])
	}
	return result
}

func (l Itinerary) routePairingKey(a1 Activity, a2 Activity) string {
	return fmt.Sprintf("%s%s%s", a1.ID, LabelDelimeter, a2.ID)
}

func (l Itinerary) MakeRoutePairings() map[string]bool {
	pairings := map[string]bool{}
	sorted := l.SortActivities()
	for i := 1; i < len(sorted); i++ {
		// We need the origin and destination to have a place
		if sorted[i-1].Place.ID == "" || sorted[i].Place.ID == "" {
			continue
		}
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
	Title     string           `json:"title" bson:"title"`
	Desc      string           `json:"desc" bson:"desc"`
	PriceItem common.PriceItem `json:"price" bson:"price"`

	Labels common.Labels `json:"labels" bson:"labels"`
	Tag    common.Tags   `json:"tags" bson:"tags"`
}

type BudgetItemsList []BudgetItem

type Link struct {
	ID     string        `json:"id" bson:"id"`
	Notes  string        `json:"notes" bson:"notes"`
	OGP    ogp.Opengraph `json:"ogp"`
	Labels common.Labels `json:"labels" bson:"labels"`
	Tags   common.Tags   `json:"tags" bson:"tags"`
}

type LinkMap map[string]Link

const (
	AttachmentTypeFile  = "file"
	AttachmentTypeMedia = "media"
)
