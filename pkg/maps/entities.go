package maps

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/travelreys/travelreys/pkg/common"
	"googlemaps.github.io/maps"
)

const (
	LabelPlaceID        = "gID"
	LabelPhotoReference = "p"
	LabelCountry        = "country"
	LabelCity           = "city"
	LabelState          = "state"

	DirectionModeDriving = "driving"
	DirectionModeWalking = "walking"
	DirectionModeTransit = "transit"
	DefaultDirectionMode = DirectionModeDriving

	DefaultWalkingMaxDuration = 20 * time.Minute
)

var (
	ErrRouteListEmpty     = errors.New("maps.ErrRouteListEmpty")
	DirectionModesAllList = []string{
		DirectionModeDriving,
		DirectionModeTransit,
		DirectionModeWalking,
	}
)

type LatLng struct {
	Lat float64 `json:"lat" bson:"lat" msgpack:"lat"`
	Lng float64 `json:"lng" bson:"lng" msgpack:"lng"`
}

type Place struct {
	ID          string        `json:"id" bson:"id" msgpack:"id"`
	Name        string        `json:"name" bson:"name" msgpack:"name"`
	Address     string        `json:"address" bson:"address" msgpack:"address"`
	LatLng      LatLng        `json:"latlng" bson:"latlng" msgpack:"latlng"`
	PhoneNumber string        `json:"phoneNumber" bson:"phoneNumber" msgpack:"phoneNumber"`
	Types       []string      `json:"types" bson:"types" msgpack:"types"`
	Website     string        `json:"website" bson:"website" msgpack:"website"`
	Labels      common.Labels `json:"labels" bson:"labels" msgpack:"labels"`
}

func PlaceFromPlaceDetailsResult(result maps.PlaceDetailsResult) Place {
	place := Place{
		ID:          uuid.NewString(),
		Name:        result.Name,
		Address:     result.FormattedAddress,
		LatLng:      LatLng(result.Geometry.Location),
		Website:     result.Website,
		Types:       result.Types,
		PhoneNumber: result.InternationalPhoneNumber,
		Labels: common.Labels{
			LabelPlaceID: result.PlaceID,
		},
	}
	if len(result.Photos) > 0 {
		place.Labels[LabelPhotoReference] = result.Photos[0].PhotoReference
	}
	for _, adr := range result.AddressComponents {
		for _, typ := range adr.Types {
			if typ == "country" {
				place.Labels[LabelCountry] = adr.LongName
			}
			if typ == "locality" {
				place.Labels[LabelCity] = adr.LongName
			}
			if typ == "administrative_area_level_1" {
				place.Labels[LabelState] = adr.LongName
			}
		}
	}

	return place
}

func (p Place) PlaceID() string {
	return p.Labels[LabelPlaceID]
}

type PlaceAtmosphere struct {
	maps.PlaceDetailsResult
}

type AutocompletePrediction struct {
	maps.AutocompletePrediction
}

type AutocompletePredictionList []AutocompletePrediction

// Polyline represents a list of lat,lng points encoded as a byte array.
// See: https://developers.google.com/maps/documentation/utilities/polylinealgorithm
type Polyline struct {
	Points string `json:"points" bson:"points" msgpack:"points"`
}

type Route struct {
	Polyline      Polyline      `json:"polyline" bson:"polyline" msgpack:"polyline"`
	Distance      int           `json:"distance" bson:"distance" msgpack:"distance"`
	Duration      time.Duration `json:"duration" bson:"duration" msgpack:"duration"`
	StartLocation LatLng        `json:"start" bson:"start" msgpack:"start"`
	EndLocation   LatLng        `json:"end" bson:"end" msgpack:"end"`
	TravelMode    string        `json:"travelMode" bson:"travelMode" msgpack:"travelMode"`
	Labels        common.Labels `json:"labels" bson:"labels" msgpack:"labels"`
}

func RouteFromRouteResult(result maps.Route, mode string) Route {
	if len(result.Legs) <= 0 {
		return Route{
			TravelMode: mode,
			Labels:     common.Labels{},
		}
	}
	return Route{
		Polyline: Polyline{
			Points: result.OverviewPolyline.Points,
		},
		Distance:      result.Legs[0].Distance.Meters,
		Duration:      result.Legs[0].Duration, // seconds
		StartLocation: LatLng(result.Legs[0].StartLocation),
		EndLocation:   LatLng(result.Legs[0].EndLocation),
		TravelMode:    mode,
		Labels:        common.Labels{},
	}
}

type RouteList []Route
type RouteListMap map[string]RouteList

// GetMostCommonSenseRoute returns the "most common sense" route
// amongs multiple route options, that is, when walking is avaliable and
// less than 20 mins. Else, it returns the shorter of driving or transit.
func (l RouteList) GetMostCommonSenseRoute() (Route, error) {
	if len(l) <= 0 {
		return Route{}, ErrRouteListEmpty
	}

	walkingIdx := -1
	for idx, r := range l {
		if r.TravelMode == DirectionModeWalking {
			walkingIdx = idx
			break
		}
	}

	// Driving or Transit
	shortest := l[0].Duration
	shortestIdx := 0
	for idx, r := range l {
		if idx == walkingIdx {
			continue
		}
		if r.Duration < shortest {
			shortest = r.Duration
			shortestIdx = idx
		}
	}

	// Walking not available, return shorter of transit or driving
	if walkingIdx <= 0 {
		return l[shortestIdx], nil
	}

	// If walking takes less than 20 mins, return walking
	walkingRoute := l[walkingIdx]
	if walkingRoute.Duration < DefaultWalkingMaxDuration {
		return walkingRoute, nil
	}
	return l[shortestIdx], nil

}

func GetShortestDurationGMapsRoute(routes []maps.Route) (maps.Route, error) {
	if len(routes) <= 0 {
		return maps.Route{}, ErrRouteListEmpty
	}
	shortest := routes[0].Legs[0].Duration
	shortestIdx := 0
	for idx, r := range routes {
		if r.Legs[0].Duration < shortest {
			shortest = r.Legs[0].Duration
			shortestIdx = idx
		}
	}
	return routes[shortestIdx], nil
}
