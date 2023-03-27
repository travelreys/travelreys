package maps

import (
	"context"

	"github.com/go-kit/kit/endpoint"
	"github.com/travelreys/travelreys/pkg/common"
)

type PlacesAutocompleteRequest struct {
	Query        string `json:"query"`
	Types        string `json:"types"`
	Lang         string `json:"lang"`
	Sessiontoken string `json:"sessiontoken"`
}
type PlacesAutocompleteResponse struct {
	Predictions AutocompletePredictionList `json:"predictions"`
	Err         error                      `json:"error,omitempty"`
}

func (r PlacesAutocompleteResponse) Error() error {
	return r.Err
}

func NewPlacesAutocompleteEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, epReq interface{}) (interface{}, error) {
		req, ok := epReq.(PlacesAutocompleteRequest)
		if !ok {
			return PlacesAutocompleteResponse{
				Err: common.ErrorEndpointReqMismatch,
			}, nil
		}
		preds, err := svc.PlacesAutocomplete(ctx, req.Query, req.Types, req.Sessiontoken, req.Lang)
		return PlacesAutocompleteResponse{Predictions: preds, Err: err}, nil
	}
}

type PlaceDetailsRequest struct {
	PlaceID      string   `json:"placeID"`
	Fields       []string `json:"fields"`
	Lang         string   `json:"lang"`
	Sessiontoken string   `json:"sessiontoken"`
}

type PlaceDetailsResponse struct {
	Place Place `json:"place"`
	Err   error `json:"error,omitempty"`
}

func (r PlaceDetailsResponse) Error() error {
	return r.Err
}

func NewPlaceDetailsEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, epReq interface{}) (interface{}, error) {
		req, ok := epReq.(PlaceDetailsRequest)
		if !ok {
			return PlaceDetailsResponse{
				Err: common.ErrorEndpointReqMismatch,
			}, nil
		}
		place, err := svc.PlaceDetails(ctx, req.PlaceID, req.Fields, req.Sessiontoken, req.Lang)
		return PlaceDetailsResponse{Place: place, Err: err}, nil
	}
}

type PlaceAtmosphereRequest struct {
	PlaceID      string   `json:"placeID"`
	Fields       []string `json:"fields"`
	Lang         string   `json:"lang"`
	Sessiontoken string   `json:"sessiontoken"`
}

type PlaceAtmosphereResponse struct {
	Place PlaceAtmosphere `json:"place"`
	Err   error           `json:"error,omitempty"`
}

func (r PlaceAtmosphereResponse) Error() error {
	return r.Err
}

func NewPlaceAtmosphereEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, epReq interface{}) (interface{}, error) {
		req, ok := epReq.(PlaceAtmosphereRequest)
		if !ok {
			return PlaceAtmosphereResponse{
				Err: common.ErrorEndpointReqMismatch,
			}, nil
		}
		place, err := svc.PlaceAtmosphere(ctx, req.PlaceID, req.Fields, req.Sessiontoken, req.Lang)
		return PlaceAtmosphereResponse{Place: place, Err: err}, nil
	}
}
