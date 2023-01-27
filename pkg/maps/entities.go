package maps

import "googlemaps.github.io/maps"

type Place struct {
	maps.PlaceDetailsResult
}

type AutocompletePrediction struct {
	maps.AutocompletePrediction
}

type AutocompletePredictionList []AutocompletePrediction
