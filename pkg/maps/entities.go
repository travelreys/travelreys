package maps

import (
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
		Duration:      result.Legs[0].Duration,
		StartLocation: LatLng(result.Legs[0].StartLocation),
		EndLocation:   LatLng(result.Legs[0].EndLocation),
		TravelMode:    mode,
		Labels:        common.Labels{},
	}
}

type RouteList []Route
