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
)

var (
	ErrRouteListEmpty     = errors.New("maps.entities.RouteListEmpty")
	DirectionModesAllList = []string{
		DirectionModeDriving,
		DirectionModeTransit,
		DirectionModeWalking,
	}
)

type LatLng struct {
	Lat float64 `json:"lat"`
	Lng float64 `json:"lng"`
}

type Place struct {
	ID          string        `json:"id"`
	Name        string        `json:"name"`
	Address     string        `json:"address"`
	LatLng      LatLng        `json:"latlng"`
	PhoneNumber string        `json:"phoneNumber"`
	Types       []string      `json:"types"`
	Website     string        `json:"website"`
	Labels      common.Labels `json:"labels"`
}

func (p Place) isEmpty() bool {
	return p.ID == ""
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
	Points string `json:"points"`
}

type Route struct {
	Polyline      Polyline      `json:"polyline"`
	Distance      int           `json:"distance"`
	Duration      time.Duration `json:"duration"`
	StartLocation LatLng        `json:"start"`
	EndLocation   LatLng        `json:"end"`
	TravelMode    string        `json:"travelMode" bson:"travelMode"`
	Labels        common.Labels `json:"labels"`
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

	walkingRoute := l[walkingIdx]
	if walkingRoute.Duration < 20*time.Minute {
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
