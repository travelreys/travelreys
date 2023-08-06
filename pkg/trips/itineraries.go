package trips

import (
	"fmt"
	"sort"
	"time"

	"github.com/google/uuid"
	"github.com/travelreys/travelreys/pkg/common"
	"github.com/travelreys/travelreys/pkg/finance"
	"github.com/travelreys/travelreys/pkg/maps"
)

type Activity struct {
	ID        string            `json:"id" bson:"id"`
	Title     string            `json:"title" bson:"title"`
	Place     maps.Place        `json:"place" bson:"place"`
	Notes     string            `json:"notes" bson:"notes"`
	PriceItem finance.PriceItem `json:"price" bson:"price"`
	StartTime time.Time         `json:"startTime" bson:"startTime"`
	EndTime   time.Time         `json:"endTime" bson:"endTime"`
	Labels    common.Labels     `json:"labels" bson:"labels"`
}

func (a Activity) HasPlace() bool {
	return a.Place.Name != ""
}

type ActivityMap map[string]*Activity
type ActivityList []*Activity

func (l ActivityList) GetFracIndexes() []string {
	result := []string{}
	for _, a := range l {
		result = append(result, a.Labels[LabelFractionalIndex])
	}
	return result
}

func (l ActivityList) Len() int {
	return len(l)
}
func (l ActivityList) Swap(i, j int) {
	l[i], l[j] = l[j], l[i]
}
func (l ActivityList) Less(i, j int) bool {
	return l[i].Labels[LabelFractionalIndex] < l[j].Labels[LabelFractionalIndex]
}

type RouteMap map[string]maps.RouteList

type Itinerary struct {
	ID          string        `json:"id" bson:"id"`
	Date        time.Time     `json:"date" bson:"date"`
	Description string        `json:"desc" bson:"desc"`
	Activities  ActivityMap   `json:"activities" bson:"activities"`
	Routes      RouteMap      `json:"routes" bson:"routes"`
	Labels      common.Labels `json:"labels" bson:"labels"`
}

func NewItinerary(date time.Time) Itinerary {
	return Itinerary{
		ID:          uuid.New().String(),
		Date:        date,
		Description: "",
		Activities:  ActivityMap{},
		Routes:      RouteMap{},
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

func (itin Itinerary) GetDate() time.Time {
	return time.Date(
		itin.Date.Year(),
		itin.Date.Month(),
		itin.Date.Day(),
		0, 0, 0, 0,
		itin.Date.Location(),
	)
}

// SortActivities returns Activities sorted by their fractional index
func (itin Itinerary) SortActivities() ActivityList {
	sorted := ActivityList{}
	for _, act := range itin.Activities {
		sorted = append(sorted, act)
	}
	sort.Sort(sorted)
	return sorted
}

func (itin Itinerary) routePairingKey(a1 *Activity, a2 *Activity) string {
	return fmt.Sprintf("%s%s%s", a1.ID, LabelDelimeter, a2.ID)
}

func (itin Itinerary) RoutePairings(lodgings LodgingsMap) map[string]bool {
	pairings := map[string]bool{}
	sorted := itin.SortActivities()

	if len(sorted) <= 0 {
		return pairings
	}

	// Find routes between lodging and first activity
	if sorted[0].Place.ID != "" {
		for _, l := range lodgings {
			act := Activity{ID: l.ID, Place: l.Place}
			pairings[itin.routePairingKey(&act, sorted[0])] = true
		}
	}

	// Find route between activities
	for i := 1; i < len(sorted); i++ {
		// We need the origin and destination to have a place
		if sorted[i-1].Place.ID == "" || sorted[i].Place.ID == "" {
			continue
		}
		pairings[itin.routePairingKey(sorted[i-1], sorted[i])] = true
	}
	return pairings
}
